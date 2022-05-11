package models

import (
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
	Cmd            string      `yaml:",omitempty"`
	DNS            string      `yaml:"dns,omitempty"`
	AllowLatest    bool        `yaml:"allow_latest,omitempty"`
	Overwrite      bool        `yaml:",omitempty"`
	Registry       string      `yaml:",omitempty"`
	RegistryUserId string      `yaml:"registry_user,omitempty"`
	OutDir         string      `yaml:",omitempty"`
	ToolsImage     string      `yaml:"tools_image,omitempty"`
	ValuesImage    string      `yaml:"values_image,omitempty"`
	Image          string      `yaml:",omitempty"`
	Starter        string      `yaml:",omitempty"`
	Params         interface{} `yaml:",omitempty"`
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

	// }
	if t.Type == "helm" {
		fmt.Println("Create chart")
		vol := path.CurrentDir() + ":/workspace"
		starterVol := path.HelmDefaultsStarters() + ":/root/.local/share/helm/starters"
		chartLoc := "/workspace/gp/helm/" + comp.Slug
		helmCmd := cmd.GenetateHelmChart(vol, t.ToolsImage, chartLoc, t.Starter, starterVol)
		fmt.Println(helmCmd)
		_, err := shell.Execute(dockerBin, helmCmd)
		if err != nil {
			fmt.Println(err)
		}

		// volWorkspace, targetVol, configVol, dataVol string, dockerImage string
		targetVol := path.Targets() + t.Name + "/target.yaml:/targets/target.yaml"
		configVol := path.AppConfig() + "config.yaml:/gp/config.yaml"
		dataVol := path.AppRoot() + "/data:/data"
		volumes := []string{
			vol, targetVol, configVol, dataVol,
		}
		helmValuesCmd := cmd.HelmValues(volumes, t.ValuesImage)
		fmt.Println(helmValuesCmd)
		_, err = shell.Execute(dockerBin, helmValuesCmd)
		if err != nil {
			fmt.Println(err)
		}
	}
	if t.Type == "oam" {
		fmt.Println("generate manifests")
	}
	if t.Type == "go" {
		fmt.Println("build go")
	}
	return nil
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
