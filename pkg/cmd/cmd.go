package cmd

import (
	"fmt"
	"time"
)

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

func ComposeTarget(composeFile, service string) []string {
	actionCMD := []string{
		"docker",
		"compose",
		"-f",
		composeFile,
		"run",
		"-T",
		service,
	}
	return actionCMD
}
