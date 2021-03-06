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
	"io/ioutil"
	"strings"

	"github.com/ipedrazas/gp/pkg/files"
	"github.com/ipedrazas/gp/pkg/models"
	"github.com/ipedrazas/gp/pkg/path"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	name        string
	global      bool
	params      string
	data        string
	fromDefault string
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
		dt := &models.Target{}
		fileName := path.DefaultTargets() + fromDefault + "/target.yaml"
		err := files.Load(fileName, dt)
		if err != nil {
			cobra.CheckErr(err)
		}
		dt.Name = name
		fmt.Println(dt)
		c := &models.Component{}
		err = c.Hydrate(viper.GetViper())
		if err != nil {
			cobra.CheckErr(err)
		}
		if dt.InDefaults(fromDefault) {
			if dt.IsAvailable() {
				dt.Save()
				writeCompose(dt, c)
			} else {
				cobra.CheckErr(errors.New("target already exists in " + path.Targets() + dt.Name))
			}
		} else {
			cobra.CheckErr(errors.New("target not found in Defaults directory"))
		}
	},
}

func init() {
	targetCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&name, "name", "n", "", "name of the target")
	addCmd.Flags().StringVarP(&fromDefault, "default", "f", "", "Target name defined in the Default repo")
	addCmd.Flags().BoolVarP(&global, "global", "g", false, "create a global target, for all projects, or just add it to the current project")
	addCmd.Flags().StringVarP(&params, "params", "p", "", "comma separated key value pairs, -p \"KEY1=Value,KEY2=Value\"")
	addCmd.Flags().StringVarP(&data, "data", "d", "", "data directory. If the target needs to feed data from a directory, this is the path")
	addCmd.Flags().StringVarP(&params, "params", "p", "", "comma separated key value pairs, -p \"KEY1=Value,KEY2=Value\"")

}

func writeCompose(t *models.Target, c *models.Component) error {

	tpl := path.DefaultTargets() + fromDefault + "/" + t.Compose
	dst := path.Targets() + t.Name + "/" + t.Compose
	content, err := ioutil.ReadFile(tpl)
	if err != nil {
		fmt.Println(err, tpl)
		return err
	}

	text := string(content)

	text = strings.Replace(text, "__TARGET_NAME__", t.Name, -1)
	text = strings.Replace(text, "__WORKSPACE__", path.CurrentDir(), -1)
	text = strings.Replace(text, "__GP_ROOT_CONFIG__/", path.AppConfig(), -1)
	text = strings.Replace(text, "__DATA__", data, -1)
	text = strings.Replace(text, "__NAME__", c.Slug, -1)

	err = ioutil.WriteFile(dst, []byte(text), 0644)
	if err != nil {
		fmt.Println(err, tpl, dst)
		return err
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
