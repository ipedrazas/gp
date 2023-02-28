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
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ipedrazas/gp/pkg/files"
	"github.com/ipedrazas/gp/pkg/models"
	"github.com/ipedrazas/gp/pkg/path"
	"github.com/ipedrazas/gp/pkg/remote"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	name        string
	global      bool
	params      string
	data        string
	fromDefault string
	oci         string
	url         string
	targetPath  string
	// build       bool
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new Target defined in your catalog",
	Long: `Adds a new Target. Targets are defined in $USER/.config/gp/defaults/gp/targets
	
	gp target add -n [New Target Name] -f [Target Name defined in Defaults] 

	If you need to add more KV, you can pass them as 
	gp target add -n [New Target Name] -f [Target Name defined in Defaults] -p [KEY=Value1,KEY2=value2]
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if oci != "" {
			remote.FetchOCI(oci, "gp/targets/"+name+"/")
			return
		}
		fromFilesystem()
		if build {
			buildCmd.Run(cmd, []string{})
		}
	},
}

func init() {
	targetCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&name, "name", "n", "", "name of the target")
	addCmd.Flags().StringVarP(&fromDefault, "default", "f", "", "Target name defined in the Default repo")
	addCmd.Flags().BoolVarP(&global, "global", "g", false, "create a global target, for all projects, or just add it to the current project")
	addCmd.Flags().StringVarP(&params, "params", "p", "", "comma separated key value pairs, -p \"KEY1=Value,KEY2=Value\"")
	addCmd.Flags().StringVar(&data, "data", "", "data directory. If the target needs to feed data from a directory, this is the path")
	addCmd.Flags().StringVarP(&domain, "domain", "d", "", "partial dns for public access to resources: .mycontext.mydomain.com ")
	addCmd.Flags().StringVarP(&registryUser, "user", "u", "", "User to access the registry")
	addCmd.Flags().StringVarP(&registry, "registry", "r", "", "Flag to specify the Docker registry to use.")
	addCmd.Flags().StringVarP(&oci, "oci", "o", "", "OCI image that contains the Target.")
	addCmd.Flags().StringVar(&url, "url", "", "URL to fetch the Target from.")
	addCmd.Flags().StringVar(&targetPath, "path", "", "path where target resources are located.")

}

func TargetFromFile(filepath string) *models.Target {
	dt := &models.Target{}

	err := files.Load(filepath, dt)
	if err != nil {
		cobra.CheckErr(err)
	}
	dt.Name = name
	dt.Domain = domain
	dt.Registry = registry
	dt.RegistryUserId = registryUser

	v := viper.GetViper()

	if domain == "" && v.IsSet("domain") {
		dt.Domain = v.GetString("domain")
	}
	if registryUser == "" && v.IsSet("registry_user") {
		dt.RegistryUserId = v.GetString("registry_user")
	}
	if registry == "" && v.IsSet("registry") {
		dt.Registry = v.GetString("registry")
	}

	return dt
}

func fromFilesystem() {

	if fromDefault == "" {
		fromDefault = name
	}
	fileName := path.DefaultTargets() + fromDefault + "/target.yaml"

	dt := TargetFromFile(fileName)

	if dt.InDefaults(fromDefault) {
		if dt.IsAvailable() {

			err := dt.Save()
			if err != nil {
				fmt.Println(err)
			}
			files, err := os.ReadDir(path.DefaultTargets() + fromDefault)
			if err != nil {
				cobra.CheckErr(err)
			}
			c := &models.Component{}
			v := viper.GetViper()
			err2 := c.Hydrate(v, false)
			if err2 != nil {
				cobra.CheckErr(err2)
			}
			for _, file := range files {
				if !file.IsDir() && file.Name() != "target.yaml" {
					err = parseTemplate(dt, c, file.Name())
					if err != nil {
						fmt.Println(err)
					}
				}
			}

		} else {
			cobra.CheckErr(errors.New("target already exists in " + path.Targets() + dt.Name))
		}
	} else {
		cobra.CheckErr(errors.New("target not found in Defaults directory"))
	}
}

func parseTemplate(t *models.Target, c *models.Component, filename string) error {

	tpl := path.DefaultTargets() + fromDefault + "/" + filename
	// dst := path.Targets() + t.Name + "/" + filename
	dst := "gp/targets/" + t.Name + "/" + filename
	content, err := os.ReadFile(tpl)
	if err != nil {
		// fmt.Println(tpl, dst)
		return err
	}

	text := string(content)

	text = strings.Replace(text, "__TARGET_NAME__", t.Name, -1)
	text = strings.Replace(text, "__WORKSPACE__", path.CurrentDir(), -1)
	text = strings.Replace(text, "__GP_ROOT_CONFIG__/", path.AppConfig(), -1)
	text = strings.Replace(text, "__DATA__", data, -1)
	text = strings.Replace(text, "__NAME__", c.Slug, -1)

	content = []byte(text)
	err2 := os.WriteFile(dst, content, 0666)
	if err2 != nil {
		return err2
	}
	return nil
}

type ComposeInput struct {
	TargetName    string `cue:"targetName" json:"target_name,omitempty"`
	Workspace     string `cue:"workspaceDir" json:"workspace,omitempty"`
	ConfigDir     string `cue:"configDir" json:"config_dir,omitempty"`
	DataDir       string `cue:"dataDir" json:"data_dir,omitempty"`
	CompName      string `cue:"compName" json:"comp_name,omitempty"`
	ComposeSource string
	YamlOut       string
}

// func writeComposeFromCUE(input *ComposeInput) error {
// 	ctx := cuecontext.New()

// 	// read and compile value
// 	d, _ := os.ReadFile(input.ComposeSource)
// 	schema := ctx.CompileBytes(d)
// 	cueData := ctx.Encode(input)

// 	return nil
// }

// slug
// workspace
// gp_config
// data
