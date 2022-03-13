package systemCommand

import (
	"os/exec"
)

func RunCmd(cmd string) (string, error) {
	CMD := exec.Command("sh", "-c", cmd)
	result, err := CMD.CombinedOutput()
	return string(result), err
}

func RunFile(path string) (string, error) {
	CMD := exec.Command("sh", path)
	result, err := CMD.CombinedOutput()
	return string(result), err
}
