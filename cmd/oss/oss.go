package oss

import (
	"fmt"
	"log"
	"os/exec"
)

func LaunchScript(cmdName string, pathToScript string, arg string) string {
	cmd, err := exec.Command(cmdName, pathToScript, arg).Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	cmdOutput := string(cmd)

	log.Printf("%v\n", cmdOutput)
	return cmdOutput

}
