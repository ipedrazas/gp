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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func getHome() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return homedir
}

// AppRoot returns "$PWD/gp/"
func AppRoot() string {
	return CurrentDir() + "/gp/"
}

func OAMRoot() string {
	return AppRoot() + "k8s/oam/"
}

// AppConfig returns "$HOME/.config/gp/"
func AppConfig() string {
	return getHome() + "/.config/gp/"
}

// AppConfig returns "$HOME/.config/gp/secrets.txt"
func AppSecrets() string {
	return AppConfig() + "secrets.txt"
}

// Defaults returns "$HOME/.config/gp/defaults/"
func Defaults() string {
	return AppConfig() + "defaults/"
}

func DefaultsK8S() string {
	return Defaults() + "k8s/"
}

// AppFile returns "$PWD/.app.yaml"
func AppFile() string {
	return CurrentDir() + "/.app.yaml"
}

// CurrentDir returns the current directory, no
// final slash /mydir
func CurrentDir() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	return path
}

func HelmChart() string {
	return AppRoot() + "k8s/helm/"
}

func Kubeconfig() string {
	return getHome() + "/.kube/config"
}

// Targets returns  "$HOME/.config/gp/targets/"
func Targets() string {
	if Exists(AppRoot() + "targets/") {
		return AppRoot() + "targets/"
	}
	return AppConfig() + "targets/"
}

func Dockerfiles() string {
	return Defaults() + "dockerfiles/"
}

func Exists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {

		return true
	}
	return false
}
func MakeDirectoryIfNotExists(path string) error {

	if !Exists(path) {
		return os.MkdirAll(path, os.ModeDir|0755)
	}
	return nil
}

func GetBinPath(tool string) string {
	path, err := exec.LookPath(tool)
	if err != nil {
		fmt.Printf("didn't find '%s' executable\n", tool)
		return ""
	}
	return path
}

func GetDirNames(dir string) []string {
	dirs := []string{}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		// If dir doesn't exist, return empty array
		return dirs
	}

	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			dirs = append(dirs, fileInfo.Name())
		}
	}
	return dirs

}
