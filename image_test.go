package undocker

import (
	"io/ioutil"
	"os"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestImage_Extract(t *testing.T) {
	tests := []struct {
		name             string
		rURL             string
		repo             string
		tag              string
		overwriteSymlink bool
		wantDirs         []string
		wantErr          bool
	}{
		{
			name: "Extract busybox from local registry",
			// tokibi/busybox-bundle-registry
			rURL:             "http://localhost:5000",
			repo:             "busybox",
			tag:              "latest",
			overwriteSymlink: false,
			wantDirs:         []string{"bin", "dev", "etc", "home", "root", "tmp", "usr", "var"},
			wantErr:          false,
		},
		{
			name:             "Extract lolipopmc/php:7.4 from docker registry",
			rURL:             "https://registry.hub.docker.com",
			repo:             "lolipopmc/php",
			tag:              "7.4",
			overwriteSymlink: false,
			wantDirs:         []string{"bin", "boot", "dev", "etc", "home", "lib", "lib64", "media", "mnt", "opt", "proc", "root", "run", "sbin", "srv", "sys", "tmp", "usr", "var"},
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rtmpdir, err := ioutil.TempDir("", "registry")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(rtmpdir)
			etmpdir, err := ioutil.TempDir("", "undocker")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(etmpdir)

			registry, err := NewRegistry(tt.rURL, "", "", rtmpdir)
			if err != nil {
				t.Fatal(err)
			}
			i := Image{
				Source:     registry,
				Repository: tt.repo,
				Tag:        tt.tag,
			}
			if err := i.Extract(etmpdir, tt.overwriteSymlink); (err != nil) != tt.wantErr {
				t.Errorf("Image.Extract() error = %v, wantErr %v", err, tt.wantErr)
			}

			files, err := ioutil.ReadDir(etmpdir)
			if err != nil {
				t.Fatal(err)
			}
			dirs := []string{}
			for _, f := range files {
				if !f.IsDir() {
					continue
				}
				dirs = append(dirs, f.Name())
			}
			sort.Slice(tt.wantDirs, func(i, j int) bool { return tt.wantDirs[i] < tt.wantDirs[j] })
			sort.Slice(dirs, func(i, j int) bool { return dirs[i] < dirs[j] })
			if diff := cmp.Diff(dirs, tt.wantDirs, nil); diff != "" {
				t.Errorf("got %v\nwant %v", dirs, tt.wantDirs)
			}
		})
	}
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}
