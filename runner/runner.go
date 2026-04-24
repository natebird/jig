package runner

import (
	"fmt"
	"os"
	"os/exec"
)

// Run executes a command in the given directory, streaming stdout/stderr live.
func Run(dir string, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command provided")
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
