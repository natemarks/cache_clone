package utility

import (
	"reflect"
	"testing"
)

func Test_urlHostAndPath(t *testing.T) {
	type args struct {
		urlString string
	}
	tests := []struct {
		name     string
		args     args
		wantHost string
		wantPath string
		wantErr  bool
	}{
		{name: "valid", args: args{urlString: "https://my.git.host/scm/group/project.git"},
			wantHost: "my.git.host",
			wantPath: "/scm/group/project.git",
			wantErr:  false},
		{name: "empty_host_from_missing_protcol", args: args{urlString: "my.git.host/scm/group/project.git"},
			wantHost: "",
			wantPath: "",
			wantErr:  true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, gotPath, err := UrlHostAndPath(tt.args.urlString)
			if (err != nil) != tt.wantErr {
				t.Errorf("UrlHostAndPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHost != tt.wantHost {
				t.Errorf("UrlHostAndPath() gotHost = %v, want %v", gotHost, tt.wantHost)
			}
			if gotPath != tt.wantPath {
				t.Errorf("UrlHostAndPath() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}

func TestMakeParentDir(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name          string
		args          args
		wantParentDir string
		wantErr       bool
	}{
		{name: "valid", args: args{p: "/tmp/gg/hh"},
			wantParentDir: "/tmp/gg",
			wantErr:       false},
		{name: "perm_denied", args: args{p: "/etc/gg/hh"},
			wantParentDir: "/etc/gg",
			wantErr:       true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParentDir, err := MakeParentDir(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeParentDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotParentDir != tt.wantParentDir {
				t.Errorf("MakeParentDir() gotParentDir = %v, want %v", gotParentDir, tt.wantParentDir)
			}
		})
	}
}

func Test_run(t *testing.T) {
	type args struct {
		c []string
	}
	tests := []struct {
		name    string
		args    args
		want    Result
		wantErr bool
	}{
		// Run ls successfully to find a file that exists
		{"succeed:list_readme", args{
			c: []string{"ls", "-b", "README.md"},
		}, Result{
			ReturnCode: 0,
			StdOut:     "README.md\n",
			StdErr:     "",
		}, false},
		// Run an executable that doesn't exist
		{"fail: executable doesn't exist", args{
			c: []string{"i_dont_exist", "first_arg", "second_arg"},
		}, Result{
			ReturnCode: 0,
			StdOut:     "",
			StdErr:     "",
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Run(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() got = %v, want %v", got, tt.want)
			}
		})
	}
}
