package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ipedrazas/gp/pkg/cmd"
	"github.com/ipedrazas/gp/pkg/files"
	"github.com/ipedrazas/gp/pkg/path"
	"github.com/ipedrazas/gp/pkg/shell"
)

type Target struct {
	Name           string
	Platform       []string
	Domain         string `yaml:"dns,omitempty"`
	Overwrite      bool   `yaml:",omitempty"`
	Registry       string `yaml:",omitempty"`
	RegistryUserId string `yaml:"registry_user,omitempty"`

	// Image is a setting to override the docker image name (app/slug name)
	Image   string `yaml:",omitempty"`
	Compose string `yaml:"compose,omitempty"`

	// Actions are the ordered list of actions the target will execute
	// A Compose file can have several services, in case order is important
	// this parameter allows you to define the right order
	Actions     []string `yaml:"actions,omitempty"`
	DockerBuild bool     `yaml:"docker_build"`
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
	dir := path.AppRoot() + "targets/" + t.Name
	return !path.Exists(dir)

}

func (t *Target) InDefaults(fromDefault string) bool {
	dir := path.DefaultTargets() + fromDefault
	return path.Exists(dir)
}

func (t *Target) Run(comp *Component, gitSha string) error {
	fmt.Println("Processing target ", t.Name)

	if t.DockerBuild {
		// v := viper.GetViper()
		dockerBin := path.GetBinPath("docker")
		plat := strings.Join(t.Platform, ",")
		// for _, plat := range t.Platform {
		t.SetDockerImage(comp.Slug, comp.Version)

		var dockerBuildCMD []string
		if path.HasBakeFile() {
			dockerBuildCMD = cmd.BuildxWithBake()
		} else {
			dockerBuildCMD = cmd.Buildx(plat, t.Image, gitSha, comp.Version, true)
		}

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
	dockerBin := path.GetBinPath("docker")
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
