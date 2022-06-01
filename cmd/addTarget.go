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
	"fmt"
	"runtime"
	"strings"

	"github.com/ipedrazas/gp/pkg/models"
	"github.com/ipedrazas/gp/pkg/path"
	"github.com/spf13/cobra"
)

var (
	t        models.Target
	platform string
	actions  string
	global   bool
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
		if !global {
			path.MakeDirectoryIfNotExists(path.AppRoot() + "targets/")
		}
		if t.IsAvailable() {
			t.Save()
		}
	},
}

func init() {
	targetCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&t.Name, "name", "n", "", "name of the target")
	addCmd.Flags().StringVarP(&platform, "platform", "p", "linux/"+runtime.GOARCH, "platform of the target: linux/amd64, linux/arm64, windows/amd64")
	addCmd.Flags().StringVarP(&t.Compose, "compose", "c", "target-compose.yaml", "docker compose file that executes the target")
	addCmd.Flags().StringVarP(&t.Registry, "registry", "r", "docker.io", "docker registry to push images")
	addCmd.Flags().StringVarP(&t.RegistryUserId, "user", "u", "", "docker registry username")
	addCmd.Flags().StringVarP(&t.DNS, "dns", "d", "", "DNS to use in artifacts")
	addCmd.Flags().StringVarP(&actions, "actions", "a", "", "comma separated docker compose services. Use this option to override the predefined order of the compose file")

	addCmd.Flags().BoolVarP(&t.DockerBuild, "build", "b", true, "execute docker build for the defined target platforms")
	addCmd.Flags().BoolVarP(&global, "global", "g", false, "create a global target, for all projects, or just add it to the current project")
	t.Platform = append(t.Platform, platform)

	if actions != "" {
		t.Actions = strings.Split(actions, ",")
	}
}
