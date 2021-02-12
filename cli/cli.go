package cli

import (
	"fmt"
	"github.com/blang/vfs"
	"github.com/mdev5000/runnr"
	"os"
	"path/filepath"
)

func Run(fs vfs.Filesystem, workingDir string) error {
	globalWasRun, err := checkForGlobal(fs, workingDir)
	if globalWasRun || err != nil {
		return err
	}

	config, err := runnr.ReadConfig(fs)
	if err != nil {
		return err
	}

	outPath := config.Outpath
	exePath := config.Executable
	fullOutPath := filepath.Join(workingDir, outPath)

	recompile, _ := shouldRecompile(os.Args)
	if recompile || !fileExists(fullOutPath) {
		if err := rebuildApp(fullOutPath, outPath, exePath); err != nil {
			return err
		}
	}

	return runnr.RunExe(outPath, os.Args)
}

func shouldRecompile(args []string) (bool, []string) {
	if len(args) < 2 {
		return false, args
	}
	lastArgIndex := len(args) - 1
	if args[lastArgIndex] == "-r" || args[lastArgIndex] == "--recompile" {
		args = args[0 : len(args)-1]
		return true, args
	}
	return false, args
}

func fileExists(path string) bool {
	f, err := os.Lstat(path)
	if err != nil || f.IsDir() {
		return false
	}
	return true
}

func rebuildApp(outPathFull, outPath, exePath string) error {
	os.Remove(outPathFull)
	fmt.Println("Recompiling....")
	if err := runnr.RunCommand("go", "build", "-o", outPath, exePath); err != nil {
		return err
	}
	if !fileExists(outPathFull) {
		return fmt.Errorf("failed to build")
	}
	fmt.Println("Recompiling done.")
	return nil
}
