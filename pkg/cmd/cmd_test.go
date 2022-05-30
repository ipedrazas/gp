package cmd

import (
	"reflect"
	"testing"
	"time"
)

func TestBuildx(t *testing.T) {

	bdate := time.Now().Format("2006-01-02 15:04:05")
	res := []string{
		"docker", "buildx", "build", "--platform", "linux/amd64",
		"--tag", "unittest", "--build-arg", "GIT_SHA=111",
		"--build-arg", "VERSION=1.1.1", "--build-arg",
		"BUILD_DATE=" + bdate, "--push", ".",
	}
	type args struct {
		platform string
		tag      string
		shaGit   string
		version  string
		push     bool
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "t01", args: args{
			platform: "linux/amd64",
			tag:      "unittest",
			shaGit:   "111",
			version:  "1.1.1",
			push:     true,
		}, want: res},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Buildx(tt.args.platform, tt.args.tag, tt.args.shaGit, tt.args.version, tt.args.push); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Buildx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComposeTarget(t *testing.T) {
	type args struct {
		composeFile string
		service     string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "t01", args: args{
			composeFile: "docker-compose.yaml",
			service:     "DoIt",
		}, want: []string{"docker", "compose", "-f", "docker-compose.yaml", "run", "-T", "DoIt"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComposeTarget(tt.args.composeFile, tt.args.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComposeTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}
