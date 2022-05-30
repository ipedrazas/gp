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

	"github.com/ipedrazas/gp/pkg/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:     "build",
	Short:   "Build the application resources",
	Long:    `Build an application generating Dockerfile, docker-compose, docker image(s) and all the Kubernetes manifests in YAML.`,
	Example: `  gp build `,
	Run: func(cmd *cobra.Command, args []string) {
		c := &models.Component{}
		err := c.Hydrate(viper.GetViper())
		if err != nil {
			cobra.CheckErr(err)
		}
		err = c.GenerateDockerfile()
		if err != nil {
			cobra.CheckErr(err)
		}
		gitSha, err := GetGitSha()
		if err != nil {
			gitSha = "no-git-repo"

		}
		fmt.Println("building for targets:")
		for _, target := range c.Targets {
			fmt.Println("\t- ", target.Name)
		}

		for _, target := range c.Targets {
			err = target.Run(c, gitSha)
			if err != nil {
				cobra.CheckErr(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func GetGitSha() (string, error) {
	revision := "HEAD"
	path := "."
	r, err := git.PlainOpen(path)
	if err != nil {
		return "", err
	}

	h, err := r.ResolveRevision(plumbing.Revision(revision))
	if err != nil {
		return "", err
	}

	return h.String(), nil
}
