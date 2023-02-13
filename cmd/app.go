/*
Copyright © 2023 Ivan Pedrazas <ipedrazas@gmail.com>

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
	"strings"

	"github.com/ipedrazas/gp/pkg/models"
	"github.com/ipedrazas/gp/pkg/path"
	"github.com/ipedrazas/gp/pkg/shell"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	app      string
	username string
	appName  string
	apps     []models.App
)

// appCmd represents the app command
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("app called")
		initApps()
		for _, a := range apps {
			fmt.Println(a.DockerRun())
			fmt.Println(a.Name + " " + appName)
			if a.Name == app {
				dockerBin := path.GetBinPath("docker")
				_, err := shell.Execute(dockerBin, a.CMD())
				if err != nil {
					cobra.CheckErr(err)
				}
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(appCmd)
	appCmd.Flags().StringVarP(&app, "app", "a", "go", "create a new app from type.")
	appCmd.Flags().StringVarP(&username, "full-name", "w", "default", "user name in the format: \"My Name myname@myemail.com\".")
	appCmd.Flags().StringVarP(&appName, "appName", "n", "go", "new app name")
}

func initApps() {

	v := viper.GetViper()
	if v.IsSet("name") {
		username = v.GetString("name")
		if v.IsSet("email") {
			username += " " + v.GetString("email")
		}
	}

	envs := []string{"USER=" + username, "NAME=" + appName, "SLUG=" + Slugify(appName)}
	app := &models.App{
		Name:    "cobra",
		Image:   "harbor.alacasa.uk/library/gotools:cobra",
		EnvVars: envs,
		Volume:  "/workspace",
	}
	apps = append(apps, *app)
	app2 := &models.App{
		Name:    "gin",
		Image:   "harbor.alacasa.uk/library/gotools:gin",
		EnvVars: envs,
		Volume:  "/workspace",
	}
	apps = append(apps, *app2)
}

func Slugify(s string) string {

	ss := strings.Split(s, "/")
	s = ss[len(ss)-1]
	s = strings.Replace(strings.ToLower(s), " ", "-", -1)
	return strings.ToLower(s)

}
