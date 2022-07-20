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
	"os"

	"github.com/ipedrazas/gp/pkg/path"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
)

var (
	registryUser string
	registry     string
	dns          string
	directory    string
	defaultsURL  string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialises the Defaults resources",
	Long: `This command checks out the git repo where the Defaults are defined.
	Defaults are a set of resources used by the tool to generate artifacts.
	These can be Dockerfiles, helm charts, kubernetes manifests or any other template needed.`,
	Example: `  gp init --defaults https://github.com/ipedrazas/defaults.git] `,
	Run: func(cmd *cobra.Command, args []string) {

		os.RemoveAll(path.Defaults())

		err := path.MakeDirectoryIfNotExists(path.AppConfig())
		if err != nil {
			fmt.Println(err)
		}

		_, err = git.PlainClone(path.Defaults(), false,
			&git.CloneOptions{
				URL:               defaultsURL,
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			})

		if err != nil {
			fmt.Println(err)
		}
		v := viper.GetViper()
		v.SetDefault("defaultsUrl", defaultsURL)

		err = v.WriteConfig()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initCmd.Flags().StringVar(&defaultsURL, "defaults", "https://github.com/ipedrazas/defaults.git", "Git repo for the Defaults.")
	initCmd.Flags().StringVarP(&directory, "dir", "d", path.AppConfig(), "directory to store defaults")
	initCmd.Flags().StringVarP(&registryUser, "user", "u", "", "User to access the registry")
	initCmd.Flags().StringVarP(&registry, "registry", "r", "docker.io", "Flag to specify the Docker registry to use.")
	initCmd.Flags().StringVar(&dns, "dns", ".home.local", "Default subdomain for dns and ingress.")

}
