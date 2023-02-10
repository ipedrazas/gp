package models

import (
	"fmt"
	"os"
	"strings"

	"github.com/ipedrazas/gp/pkg/files"
	"github.com/ipedrazas/gp/pkg/path"
	"github.com/spf13/viper"
)

type Component struct {
	Name      string   `yaml:",omitempty"`
	Slug      string   `yaml:",omitempty"`
	Port      int      `yaml:",omitempty"`
	Version   string   `yaml:",omitempty"`
	Workspace string   `yaml:",omitempty"`
	Src       string   `yaml:",omitempty"`
	Lang      string   `yaml:",omitempty"`
	Cmd       string   `yaml:",omitempty"`
	Overwrite bool     `yaml:",omitempty"`
	Config    Conf     `yaml:",omitempty"`
	Targets   []Target `yaml:",omitempty"`
	Secrets   Secrets  `yaml:",omitempty"`
	Path      string   `yaml:",omitempty"`
}

func (c *Component) GenerateDockerfile() error {

	if !path.Exists(path.CurrentDir() + "/Dockerfile") {
		fmt.Println("Generating Dockerfile", path.CurrentDir()+"/Dockerfile")
		if strings.ToLower(c.Lang) == "go" {
			tpl := path.Dockerfiles() + "go/go.Dockerfile"
			err := c.writeDockerfile(tpl)
			if err != nil {
				fmt.Println("Error Generating Dockerfile", path.CurrentDir()+"/Dockerfile")
				return err
			}
		}
		if strings.ToLower(c.Lang) == "python" {
			tpl := path.Dockerfiles() + "python/python.Dockerfile"
			err := c.writeDockerfile(tpl)
			if err != nil {
				return err
			}
		}
		if strings.ToLower(c.Lang) == "static" {
			tpl := path.Dockerfiles() + "static/nginx.Dockerfile"
			err := c.writeDockerfile(tpl)
			if err != nil {
				return err
			}
		}
		if strings.ToLower(c.Lang) == "node" {
			tpl := path.Dockerfiles() + "node/node.Dockerfile"
			err := c.writeDockerfile(tpl)
			if err != nil {
				return err
			}
		}

		err := files.Copy(path.Dockerfiles()+"dockerignore", path.CurrentDir()+"/.dockerignore")
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Component) writeDockerfile(tpl string) error {

	dst := path.CurrentDir() + "/Dockerfile"
	content, err := os.ReadFile(tpl)
	if err != nil {
		fmt.Println(err, tpl)
		return err
	}
	text := c.RenderTemplate(content)
	err = os.WriteFile(dst, []byte(text), 0644)
	if err != nil {
		fmt.Println(err, tpl, dst)
		return err
	}
	return nil
}

func (c *Component) RenderTemplate(content []byte) string {
	text := string(content)
	text = strings.Replace(text, "__NAME__", c.Slug, -1)
	text = strings.Replace(text, "__GIT_REPO__", c.Src, -1)
	text = strings.Replace(text, "__ENTRYPOINT__", parseCMD(c.Cmd), -1)
	return text
}

func parseCMD(cmd string) string {
	res := "["
	cmds := strings.Split(cmd, " ")
	for k, entry := range cmds {
		if len(cmds) == 1 {
			if entry[0:1] != "/" {
				entry = "/" + entry
			}
		}

		res += "\"" + entry + "\", "
		lastItem := len(cmds) - 1
		if k == lastItem {
			// remove the trailing comma plus space
			res = res[:len(res)-2]
		}
	}
	res += "]"
	return res
}

func (c *Component) Hydrate(v *viper.Viper) error {

	files.Load(path.AppFile(), c)
	err := path.MakeDirectoryIfNotExists(path.AppRoot())
	if err != nil {
		fmt.Println("failed loading the app file")
		return err
	}

	targetDirs := path.GetDirNames(path.Targets())
	c.Targets = []Target{}
	for _, target := range targetDirs {
		t := &Target{}
		fileName := path.Targets() + target + "/target.yaml"
		err = files.Load(fileName, t)
		if err != nil {
			fmt.Printf("Warn: target %s cannot be read\n\n", fileName)
			continue
		}
		if t.Registry == "" {
			t.Registry = v.GetString("docker.registry")
		}
		if t.RegistryUserId == "" {
			t.RegistryUserId = v.GetString("docker.user")
		}

		c.Targets = append(c.Targets, *t)
	}
	if c.Slug == "" {
		c.Slug = c.Name
	}

	if c.Cmd == "" {
		c.Cmd = c.Slug
	}

	return nil
}
