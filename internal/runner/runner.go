package runner

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/febriliankr/lazymake/internal/parser"
)

func Run(target parser.Target, extraArgs []string, dryRun bool) error {
	args := []string{}

	// Only add -f if the Makefile isn't in the current directory with default name
	cwd, _ := os.Getwd()
	defaultPath := filepath.Join(cwd, "Makefile")
	absFile, _ := filepath.Abs(target.File)
	if absFile != defaultPath {
		args = append(args, "-f", target.File)
	}

	args = append(args, target.Name)
	args = append(args, extraArgs...)

	if dryRun {
		cmdStr := "make"
		for _, a := range args {
			cmdStr += " " + a
		}
		os.Stdout.WriteString(cmdStr + "\n")
		return nil
	}

	cmd := exec.Command("make", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
