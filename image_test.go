package undocker

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestImage_Extract(t *testing.T) {
	// tokibi/busybox-bundle-registry
	rtmpdir, err := ioutil.TempDir("", "registry")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(rtmpdir)
	registry, err := NewRegistry("http://localhost:5000", "", "", rtmpdir)
	if err != nil {
		t.Error(err)
	}
	etmpdir, err := ioutil.TempDir("", "undocker")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(etmpdir)

	type fields struct {
		Source     Source
		Repository string
		Tag        string
	}
	type args struct {
		dir              string
		overwriteSymlink bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Source:     registry,
				Repository: "busybox",
				Tag:        "latest",
			},
			args: args{
				dir:              etmpdir,
				overwriteSymlink: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Image{
				Source:     tt.fields.Source,
				Repository: tt.fields.Repository,
				Tag:        tt.fields.Tag,
			}
			if err := i.Extract(tt.args.dir, tt.args.overwriteSymlink); (err != nil) != tt.wantErr {
				t.Errorf("Image.Extract() error = %v, wantErr %v", err, tt.wantErr)
			}
			if info, err := os.Stat(filepath.Join(etmpdir, "bin")); err != nil || !info.IsDir() {
				t.Error("Extract failed")
			}
		})
	}
}
