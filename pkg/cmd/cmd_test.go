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

func TestGenetateHelmChart(t *testing.T) {
	res := []string{
		"docker",
		"run",
		"--rm",
		"-v",
		"volWorkspace:/workspace",
		"-v",
		"/Users/ivan/.kube:/root/.kube",
		"-v",
		"starter:/starter",
		"dockerImage",
		"sh",
		"-c",
		"helm create name --starter starter",
	}
	type args struct {
		volWorkspace  string
		dockerImage   string
		name          string
		starterName   string
		starterVolume string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "t01", args: args{
			volWorkspace:  "volWorkspace:/workspace",
			dockerImage:   "dockerImage",
			name:          "name",
			starterName:   "starter",
			starterVolume: "starter:/starter",
		}, want: res},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenetateHelmChart(tt.args.volWorkspace, tt.args.dockerImage, tt.args.name, tt.args.starterName, tt.args.starterVolume); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenetateHelmChart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHelmValues(t *testing.T) {
	res := []string{
		"docker",
		"run",
		"--rm",
		"-v",
		"/data:/data",
		"image",
	}

	type args struct {
		volumes     []string
		dockerImage string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "t01",
			args: args{
				volumes: []string{
					"/data:/data",
				},
				dockerImage: "image",
			},
			want: res,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HelmValues(tt.args.volumes, tt.args.dockerImage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HelmValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
