package util

import "os/exec"

func CommandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func RunCommand(command string) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Run()
}
