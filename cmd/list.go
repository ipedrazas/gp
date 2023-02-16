/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/ipedrazas/gp/pkg/files"
	"github.com/ipedrazas/gp/pkg/models"
	"github.com/ipedrazas/gp/pkg/path"
	"github.com/spf13/cobra"
)

var all bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nd2 Target List:")
		fmt.Println("\n  Targets defined in the project: " + path.Targets())
		localTargets := getTargets(path.Targets())
		printTargets(localTargets)

		if all {
			fmt.Println("\nTargets Availables in the system: " + path.DefaultTargets())
			defTargets := getTargets(path.DefaultTargets())
			printTargets(defTargets)
		}

	},
}

func printTargets(targets []*models.Target) {
	for _, t := range targets {
		fmt.Println("    " + t.Name)
		for _, act := range t.Actions {
			fmt.Println("      - " + act)
		}
	}
}

func getTargets(targetPath string) []*models.Target {

	targets := []*models.Target{}

	entries, err := os.ReadDir(path.DefaultTargets())
	if err != nil {
		cobra.CheckErr(err)
	}
	for _, file := range entries {
		if file.IsDir() {
			fileName := path.DefaultTargets() + file.Name() + "/target.yaml"
			dt := &models.Target{}
			err := files.Load(fileName, dt)
			if err != nil {
				cobra.CheckErr(err)
			}
			if dt.Name == "__TARGET_NAME__" {
				dt.Name = file.Name()
			}
			targets = append(targets, dt)
		}
	}
	return targets
}

func init() {
	targetCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&all, "all", false, "return all targets found in the local project and the system")
}
