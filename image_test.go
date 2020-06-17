package undocker

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestImage_Extract(t *testing.T) {
	tests := []struct {
		name             string
		rURL             string
		repo             string
		tag              string
		overwriteSymlink bool
		wantErr          bool
	}{
		{
			name: "Extract busybox from local registry",
			// tokibi/busybox-bundle-registry
			rURL:             "http://localhost:5000",
			repo:             "busybox",
			tag:              "latest",
			overwriteSymlink: false,
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
				t.Error(err)
			}
			i := Image{
				Source:     registry,
				Repository: tt.repo,
				Tag:        tt.tag,
			}
			if err := i.Extract(etmpdir, tt.overwriteSymlink); (err != nil) != tt.wantErr {
				t.Errorf("Image.Extract() error = %v, wantErr %v", err, tt.wantErr)
			}
			if info, err := os.Stat(filepath.Join(etmpdir, "bin")); err != nil || !info.IsDir() {
				t.Error("Extract failed")
			}
		})
	}
}
