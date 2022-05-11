package cmd

import (
	"fmt"
	"time"
)

func GenetateHelmChart(volWorkspace, dockerImage, name, starterName, starterVolume string) []string {

	// Helm creates the dir structure if it doesn't exists
	command := []string{
		"docker",
		"run",
		"--rm",
		"-v",
		volWorkspace,
		"-v",
		"/Users/ivan/.kube:/root/.kube",
		"-v",
		starterVolume,
		dockerImage,
		"sh",
		"-c",
		"helm create " + name + " --starter " + starterName,
	}
	return command
}

func Buildx(platform, tag, shaGit, version string, push bool) []string {

	command := []string{
		"docker",
		"buildx",
		"build",
		"--platform",
		platform,
		"--tag",
		tag,
		"--build-arg",
		fmt.Sprintf("GIT_SHA=%s", shaGit),
		"--build-arg",
		fmt.Sprintf("VERSION=%s", version),
	}
	command = append(command, "--build-arg")

	now := time.Now().Format("2006-01-02 15:04:05")
	bd := fmt.Sprintf("BUILD_DATE=%s", now)
	command = append(command, bd)

	if push {
		command = append(command, "--push")
	}

	command = append(command, ".")

	return command
}

func HelmValues(volumes []string, dockerImage string) []string {

	// Helm creates the dir structure if it doesn't exists
	command := []string{
		"docker",
		"run",
		"--rm",
	}
	for _, vol := range volumes {
		command = append(command, "-v")
		command = append(command, vol)
	}

	command = append(command, dockerImage)
	return command
}

// docker run -it
// -v $(pwd):/workspace
// -v /Users/ivan/.config/gp/targets/local/target.yaml:/targets/target.yaml
// -v /Users/ivan/.config/gp/config.yaml:/gp/config.yaml
// -v /tmp/catalog.json:/data/catalog.json ipedrazas/python
