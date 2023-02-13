package models

import (
	"fmt"

	"github.com/ipedrazas/gp/pkg/path"
)

type App struct {
	Name    string
	Image   string
	Cmd     []string
	EnvVars []string
	Volume  string
}

func (app *App) DockerRun() string {
	cmd := "docker run --rm "
	for _, env := range app.EnvVars {
		cmd += fmt.Sprintf("-e %s ", env)
	}
	currentDir := path.CurrentDir()
	cmd += fmt.Sprintf("-v %s:%s ", currentDir, app.Volume)
	cmd += app.Image
	return cmd
}

func (app *App) CMD() []string {
	cmd := []string{"docker", "run", "--rm"}
	for _, env := range app.EnvVars {
		cmd = append(cmd, "-e", env)
	}
	currentDir := path.CurrentDir()
	cmd = append(cmd, "-v", currentDir+":"+app.Volume)
	cmd = append(cmd, app.Image)
	return cmd
}
