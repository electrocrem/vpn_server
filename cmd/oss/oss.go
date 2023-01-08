package oss

import (
	"fmt"
	"log"
	"os/exec"
)

func LaunchScript(pathToScript string, arg string) string {
	cmd, err := exec.Command(pathToScript, arg).Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	cmdOutput := string(cmd)

	log.Printf("%v\n", cmdOutput)
	return cmdOutput

}
