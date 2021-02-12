package runnr

import (
	"os"
	"os/exec"
	"syscall"
)

func RunCommand(cmd string, args ...string) error {
	pre := exec.Command(cmd, args...)
	pre.Stdin = os.Stdin
	pre.Stdout = os.Stdout
	pre.Stderr = os.Stderr
	if err := pre.Start(); err != nil {
		return err
	}
	if err := pre.Wait(); err != nil {
		return err
	}
	return nil
}

func RunExe(outPath string, args []string) error {
	return syscall.Exec(outPath, args, os.Environ())
}
