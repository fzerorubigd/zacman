package zsh

import (
	"assert"
	"fmt"
	"os"
	"os/exec"
)

var (
	zshPath string = "/usr/bin/zsh"
)

func Exec(command string, env ...string) error {
	cmd := exec.Command(zshPath, "-c", command)
	cmd.Env = env

	f, err := cmd.CombinedOutput()
	fmt.Print(string(f))
	return err
}

// SetGitPath to change the default path
func SetZshPath(new string) error {
	cmd := exec.Command(new, "--version")
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}
	zshPath = new

	return nil
}

func init() {
	tmp := os.Getenv("ZACMAN_ZSH_PATH")
	if tmp != "" {
		assert.Nil(SetZshPath(tmp), tmp)
	}
}
