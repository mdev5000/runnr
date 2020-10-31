package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

const ConfigFileName = "runnr.yml"

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
	configFile := filepath.Join(wd, ConfigFileName)

	if tryInitialize(wd, configFile) {
		return
	}

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
	fullOutPath := filepath.Join(wd, outPath)

	if shouldRecompile() {
		rebuildApp(fullOutPath, outPath, exePath)
	}

	if !fileExists(fullOutPath) {
		rebuildApp(fullOutPath, outPath, exePath)
	}

	ensureOk(syscall.Exec(outPath, os.Args, os.Environ()))
}

func tryInitialize(workingDir string, configPath string) bool {
	if len(os.Args) < 2 {
		return false
	}
	if !(os.Args[1] == "-n" || os.Args[1] == "--new") {
		return false
	}
	fmt.Println("Initializing...")
	if fileExists(configPath) {
		panic(fmt.Errorf("cannot initialize config file '%s' already exists", ConfigFileName))
	}
	runnrLocalDir := "runnr"
	runnrLocalApp := filepath.Join(runnrLocalDir, "main.go")
	runnrLocalDirFull := filepath.Join(workingDir, runnrLocalDir)
	runnrLocalAppFull := filepath.Join(workingDir, runnrLocalApp)
	if fileExists(runnrLocalAppFull) {
		panic(fmt.Errorf("cannot initialize app already exists '%s'", runnrLocalApp))
	}
	ensureOk(os.MkdirAll(runnrLocalDirFull, 0775))
	ensureOk(ioutil.WriteFile(runnrLocalAppFull, []byte(`
package main

import (
	"bitbucket.org/mdev5000/runnr"
	"context"
	"fmt"
	"github.com/spf13/cobra"
)

func main() {
	ctx := context.Background()
	runner := runnr.NewRunner(ctx)

	runner.AddCommand(&cobra.Command{
		Use: "hello",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello world")
		},
	})

	if err := runner.Run(); err != nil {
		panic(err)
	}
}
`), 0775))
	out, err := yaml.Marshal(&config{
		Outpath: "./.tmpbins/runnr-local",
		Executable:    "./" + runnrLocalApp,
	})
	ensureOk(err)
	ensureOk(ioutil.WriteFile(configPath, out, 0644))
	fmt.Println("Done.")
	return true
}

func shouldRecompile() bool {
	if len(os.Args) < 2 {
	return false
	}
	if os.Args[1] == "-r" || os.Args[1] == "--recompile" {
		return true
	}
	lastArgIndex := len(os.Args) - 1
	if os.Args[lastArgIndex] == "-r" || os.Args[lastArgIndex] == "--recompile" {
		os.Args = os.Args[0:len(os.Args) - 1]
		return true
	}
	return false
}

func validateConfig(config *config) {
	if config.Executable == "" {
		panic("must specify 'execute' variable in runnr config")
	}
	if config.Outpath == "" {
		panic("must specify 'output' variable in runnr config")
	}
	if strings.HasSuffix(config.Outpath, ".go") {
		panic("outpath variable should not end with .go, you might have mixed up execute and outpath variables")
	}
}

func fileExists(path string) bool {
	f, err := os.Lstat(path)
	if err != nil || f.IsDir() {
		return false
	}
	return true
}

func rebuildApp(outPathFull, outPath, exePath string) {
	os.Remove(outPathFull)
	fmt.Println("Recompiling....")
	ensureOk(runCommand("go", "build", "-o", outPath, exePath))
	if !fileExists(outPathFull) {
		panic("failed to build")
	}
	fmt.Println("Recompiling done.")
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

