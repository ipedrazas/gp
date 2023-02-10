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

	"github.com/ipedrazas/gp/pkg/path"
	"github.com/ipedrazas/gp/pkg/shell"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	app      string
	username string
	appName  string
	apps     []App
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
		for _, app := range apps {
			fmt.Println(app.toDockerRun())
			if app.Name == "cobra" {
				dockerBin := path.GetBinPath("docker")
				_, err := shell.Execute(dockerBin, app.toCMD())
				if err != nil {
					cobra.CheckErr(err)
				}
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

type App struct {
	Name    string
	Image   string
	CMD     []string
	ENVVARS []string
	Volume  string
}

func initApps() {

	// docker run --rm -e USER="Ivan Pedrazas ipedrazas@gmail.com" -e NAME="testabc" -v $(pwd):/workspace
	v := viper.GetViper()
	if v.IsSet("name") {
		username = v.GetString("name")
		if v.IsSet("email") {
			username += " " + v.GetString("email")
		}
	}

	envs := []string{"USER=" + username, "NAME=" + appName}
	app := &App{
		Name:    "cobra",
		Image:   "harbor.alacasa.uk/library/gotools:cobra",
		ENVVARS: envs,
		Volume:  "/workspace",
	}
	apps = append(apps, *app)
}

func (app *App) toDockerRun() string {
	cmd := "docker run --rm "
	for _, env := range app.ENVVARS {
		cmd += fmt.Sprintf("-e %s ", env)
	}
	currentDir := path.CurrentDir()
	cmd += fmt.Sprintf("-v %s:%s ", currentDir, app.Volume)
	cmd += app.Image
	return cmd
}

func (app *App) toCMD() []string {
	cmd := []string{"docker", "run", "--rm"}
	for _, env := range app.ENVVARS {
		cmd = append(cmd, "-e", env)
	}
	currentDir := path.CurrentDir()
	cmd = append(cmd, "-v", currentDir+":"+app.Volume)
	cmd = append(cmd, app.Image)
	return cmd
}
