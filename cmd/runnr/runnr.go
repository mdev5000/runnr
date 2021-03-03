package main

import (
	"github.com/blang/vfs"
	"github.com/mdev5000/runnr/cli"
	"os"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	fs := vfs.OS()
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}
	return cli.Run(fs, workingDir)
}
