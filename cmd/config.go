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

	"github.com/ipedrazas/gp/pkg/models"
	"github.com/ipedrazas/gp/pkg/path"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	registryUser string
	registry     string
	kubeconfig   string
	dns          string
	push         bool
	overwrite    bool
	verbose      bool
	directory    string
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:     "config",
	Short:   "Config gp with sane Deafults",
	Long:    `Config gp. Creates folder structure, config file and pulls the defaults repo definition`,
	Example: `  gp configure --user myuser --registry docker.io --defaults https://gitea.alacasa.uk/ivan/defaults.git] `,
	Run: func(cmd *cobra.Command, args []string) {

		v := viper.GetViper()
		v.SetDefault("docker.registry", registry)
		v.SetDefault("docker.user", registryUser)
		if dns != "" {
			v.SetDefault("dns", dns)
		}
		v.SetDefault("docker.push", push)
		v.SetDefault("docker.overwrite", overwrite)
		v.SetDefault("docker.verbose", verbose)
		v.SetDefault("docker.bin", path.GetBinPath("docker"))
		v.SetDefault("helm.registry", registry)

		err := v.WriteConfig()
		if err != nil {
			fmt.Println(err)
		}
		// we add the local target
		t := &models.Target{
			Name:     "local",
			Platform: []string{"linux/" + runtime.GOARCH},
		}
		if t.IsAvailable() {
			t.Save()
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringVarP(&directory, "dir", "d", path.AppConfig(), "directory to store defaults")
	configCmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "k", path.Kubeconfig(), "kubeconfig location")
	configCmd.Flags().StringVarP(&registryUser, "user", "u", "", "User to access the registry")
	configCmd.Flags().StringVarP(&registry, "registry", "r", "docker.io", "Flag to specify the Docker registry to use.")
	configCmd.Flags().StringVar(&dns, "dns", "", "Default subdomain for dns and ingress.")

	configCmd.Flags().BoolVar(&push, "push", true, "Push image after build")
	configCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite local Dockerfile with the defaults version.")
	configCmd.Flags().BoolVar(&verbose, "verbose", false, "Display build info.")

}
