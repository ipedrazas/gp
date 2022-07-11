/*
Copyright © 2022 Ivan Pedrazas <ipedrazas@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package path

import (
	"reflect"
	"testing"
)

const (
	home       = "/Users/ivan"
	currentDir = home + "/workspace/gp/pkg/path"
)

func Test_getHome(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "t01", want: home},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHome(); got != tt.want {
				t.Errorf("getHome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppRoot(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "t01", want: currentDir + "/gp/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppRoot(); got != tt.want {
				t.Errorf("AppRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOAMRootPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "t01", want: currentDir + "/gp/k8s/oam/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OAMRoot(); got != tt.want {
				t.Errorf("OAMRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppConfigPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "t01", want: home + "/.config/gp/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppConfig(); got != tt.want {
				t.Errorf("D2ConfigPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestD2K8sDefaultsPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "t01", want: home + "/.config/gp/defaults/k8s/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultsK8S(); got != tt.want {
				t.Errorf("D2K8sDefaultsPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppFile(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "t01", want: currentDir + "/.meta.yaml"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppFile(); got != tt.want {
				t.Errorf("AppFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCurrentDir(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "t01", want: currentDir},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CurrentDir(); got != tt.want {
				t.Errorf("getCurrentDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHelmChart(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "t01", want: currentDir + "/gp/k8s/helm/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HelmChart(); got != tt.want {
				t.Errorf("HelmChartPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDirNames(t *testing.T) {
	dirs := []string{"cmd", "files", "models", "path", "shell", "targets"}
	dirs2 := []string{".git", "cmd", "dist", "gp", "pkg"}
	type args struct {
		dir string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "t01", args: args{dir: home + "/workspace/gp/pkg/"}, want: dirs},
		{name: "t02", args: args{dir: home + "/workspace/gp/"}, want: dirs2},
		{name: "t03", args: args{dir: home + "/workspace/gp/pkg/path/"}, want: []string{}},
		// Check that If dir doesn't exist, return empty array
		{name: "t04", args: args{dir: home + "/workspace/gp/pkg/path/nil"}, want: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDirNames(tt.args.dir); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDirNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppSecrets(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "t01", want: home + "/.config/gp/secrets.txt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppSecrets(); got != tt.want {
				t.Errorf("AppSecrets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKubeconfig(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "t01", want: home + "/.kube/config"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Kubeconfig(); got != tt.want {
				t.Errorf("Kubeconfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "t01", args: args{path: "/tmp/doesnotexists"}, want: false},
		{name: "t01", args: args{path: "/tmp"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Exists(tt.args.path); got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTargets(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "t01", want: home + "/.config/gp/targets/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Targets(); got != tt.want {
				t.Errorf("Targets() = %v, want %v", got, tt.want)
			}
		})
	}
}
