package oss

import (
	"os/exec"
)

func LaunchScript(cmdName string, pathToScript string, arg string) error {
	cmd := exec.Command(cmdName, pathToScript, arg)
	err := cmd.Run()

	return err

}
