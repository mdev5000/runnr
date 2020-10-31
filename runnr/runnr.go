package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"gopkg.in/yaml.v2"
	"syscall"
)

type config struct {
	Outpath     string     `yml:"outpath"`
	Executable  string     `yml:"executable"`
}

func ensureOk(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	wd, err := os.Getwd()
	ensureOk(err)
	configFile := filepath.Join(wd, "runnr.yml")
	if !fileExists(configFile) {
		ensureOk(fmt.Errorf("Missing config file 'runnr.yml'"))
	}
	configBytes, err := ioutil.ReadFile("runnr.yml")
	ensureOk(err)
	config := &config{}
	err = yaml.Unmarshal(configBytes, config)
	ensureOk(err)
	validateConfig(config)

	outPath := config.Outpath
	exePath := config.Executable
	//args := os.Args

	if len(os.Args) >= 2 {
		if os.Args[1] == "-n" || os.Args[1] == "--new" {
			fmt.Println("initializing")
			return
		}
		if os.Args[1] == "-r" || os.Args[1] == "--recompile" {
			rebuildApp(outPath, exePath)
			//args := make([]string, len(os.Args) - 1)
			//args[0] = os.Args[1]
			//for i, arg := os.Args[1:len(os.Args) - 1] {
			//	args[i + 1] = arg
			//}
		}
	}

	fullOutPath := filepath.Join(wd, outPath)
	if !fileExists(fullOutPath) {
		rebuildApp(outPath, exePath)
	}

	// otherwise run
	ensureOk(syscall.Exec(outPath, os.Args, os.Environ()))
}

func validateConfig(config *config) {
}

func fileExists(path string) bool {
	f, err := os.Lstat(path)
	if err != nil || f.IsDir() {
		return false
	}
	return true
}

func rebuildApp(outPath, exePath string) {
	fmt.Println("rebuilding....")
	ensureOk(runCommand("go", "build", "-o", outPath, exePath))
	fmt.Println("rebuilding done.")
}

func runCommand(cmd string, args ...string) error {
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
