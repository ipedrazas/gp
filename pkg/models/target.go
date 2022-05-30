package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ipedrazas/gp/pkg/cmd"
	"github.com/ipedrazas/gp/pkg/files"
	"github.com/ipedrazas/gp/pkg/path"
	"github.com/ipedrazas/gp/pkg/shell"

	"github.com/spf13/viper"
)

type Target struct {
	Name           string
	Type           string
	Platform       string
	Cmd            string   `yaml:",omitempty"`
	DNS            string   `yaml:"dns,omitempty"`
	AllowLatest    bool     `yaml:"allow_latest,omitempty"`
	Overwrite      bool     `yaml:",omitempty"`
	Registry       string   `yaml:",omitempty"`
	RegistryUserId string   `yaml:"registry_user,omitempty"`
	OutDir         string   `yaml:"out_dir,omitempty"`
	Image          string   `yaml:",omitempty"`
	Compose        string   `yaml:"compose,omitempty"`
	Actions        []string `yaml:"actions,omitempty"`
	DockerBuild    bool     `yaml:"docker_build,omitempty"`
}

func (target *Target) SetDockerImage(appName string, tag string) {
	if target.Registry == "docker.io" || target.Registry == "" {
		target.Image = target.RegistryUserId + "/" + appName
	} else {
		target.Image = target.Registry + "/" + target.RegistryUserId + "/" + appName
	}
	if tag != "" {
		target.Image += ":" + tag
	}
}

func (t *Target) Save() {
	targetDir := path.Targets() + t.Name
	path.MakeDirectoryIfNotExists(targetDir)
	files.SaveAsYaml(targetDir+"/target.yaml", t)
}

func (t *Target) IsAvailable() bool {
	dir := path.Targets() + t.Name
	return !path.Exists(dir)

}

func (t *Target) Run(comp *Component, gitSha string) error {
	fmt.Println("Processing target ", t.Name)

	if t.DockerBuild {
		v := viper.GetViper()
		dockerBin := v.GetString("docker.bin")
		t.SetDockerImage(comp.Slug, comp.Version+"-"+getArch(t.Platform))

		dockerBuildCMD := cmd.Buildx(t.Platform, t.Image, gitSha, comp.Version, true)
		fmt.Println(dockerBuildCMD)
		_, err := shell.Execute(dockerBin, dockerBuildCMD)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return t.ExecuteTargetCompose()

}

func getArch(platform string) string {
	if strings.Contains(platform, "/") {
		res := strings.Split(platform, "/")
		if len(res) > 1 {
			return res[1]
		}
	}
	return ""
}

func (t *Target) GetCompose() *Compose {
	c := &Compose{}
	composePath := path.Targets() + t.Name + "/" + t.Compose
	files.Load(composePath, c)
	return c
}

func (t *Target) ExecuteTargetCompose() error {
	v := viper.GetViper()
	dockerBin := v.GetString("docker.bin")
	c := t.GetCompose()
	services := c.GetServiceNames()
	if len(t.Actions) == 0 {
		t.Actions = services
	}
	if validateActions(t.Actions, services) {
		for _, action := range t.Actions {
			compose := path.Targets() + t.Name + "/" + t.Compose
			if !path.Exists(compose) {
				return errors.New("compose file not found in target directory " + compose)
			}
			actionCMD := cmd.ComposeTarget(compose, action)
			fmt.Println(actionCMD)
			_, err := shell.Execute(dockerBin, actionCMD)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		return errors.New("target actions and compose services mismatch")
	}
	return nil
}

func validateActions(actions, services []string) bool {
	for _, action := range actions {
		if !contains(action, services) {
			return false
		}
	}
	return true
}

func contains(elem string, target []string) bool {
	for _, t := range target {
		if elem == t {
			return true
		}
	}
	return false
}
