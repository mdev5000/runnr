package main

import (
	"fmt"
	"github.com/spf13/cobra"
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
	Outpath    string `yml:"outpath"`
	Executable string `yml:"executable"`
}

func ensureOk(err error) {
	if err != nil {
		panic(err)
	}
}

var configFile, workingDir string

func main() {
	var err error
	workingDir, err = os.Getwd()
	ensureOk(err)
	configFile = filepath.Join(workingDir, ConfigFileName)

	if checkForGlobal() {
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
	fullOutPath := filepath.Join(workingDir, outPath)

	recompile, _ := shouldRecompile(os.Args)
	if recompile {
		rebuildApp(fullOutPath, outPath, exePath)
	}

	if !fileExists(fullOutPath) {
		rebuildApp(fullOutPath, outPath, exePath)
	}

	//ensureOk(syscall.Exec(outPath, os.Args, os.Environ()))
	runExe(outPath, os.Args)
}

func checkForGlobal() bool {
	if len(os.Args) < 2 {
		return false
	}
	if !(os.Args[1] == "g" || os.Args[1] == "global") {
		return false
	}
	runGlobal()
	return true
}

func runGlobal() {
	cmdGlobalRoot := cobra.Command{
		Use:   "global",
		Short: "Global runnr commands",
		Aliases: []string{"g"},
	}

	newCmd := &cobra.Command{
		Use:   "new",
		Short: "generate a new local project",
		Run: func(cmd *cobra.Command, args []string) {
			newProject()
		},
	}

	cmdGlobalRoot.AddCommand(newCmd)
	s := &staticCmd{}
	cmdGlobalRoot.AddCommand(s.Cmd())

	rootCmd := cobra.Command{
		Use: "runnr",
	}
	rootCmd.AddCommand(&cmdGlobalRoot)

	ensureOk(cmdGlobalRoot.Execute())
}

func newProject() {
	fmt.Println("Initializing...")
	if fileExists(configFile) {
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
		Outpath:    "./.tmpbins/runnr-local",
		Executable: "./" + runnrLocalApp,
	})
	ensureOk(err)
	ensureOk(ioutil.WriteFile(configFile, out, 0644))
	fmt.Println("Done.")
}

func shouldRecompile(args []string) (bool, []string) {
	if len(args) < 2 {
		return false, args
	}
	//if args[1] == "-r" || args[1] == "--recompile" {
	//	return true, args
	//}
	lastArgIndex := len(args) - 1
	if args[lastArgIndex] == "-r" || args[lastArgIndex] == "--recompile" {
		args = args[0 : len(args)-1]
		return true, args
	}
	return false, args
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

func runExe(outPath string, args []string) {
	ensureOk(syscall.Exec(outPath, args, os.Environ()))
}

// Static

type staticCmd struct {}

func (s *staticCmd) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "static",
		Short: "Commands for running static projects",
		Aliases: []string{"s"},
	}
	cmd.AddCommand(&cobra.Command{
		Use: "run [gofile] [outpath]",
		Short: "Run a static command",
		Args: cobra.MinimumNArgs(2),
		Run: s.runStatic,
	})
	return cmd
}

func (s *staticCmd) runStatic(cmd *cobra.Command, args []string) {
	pathToApp := args[0]
	outPath := args[1]
	appArgs := os.Args[6:] // ignore the: runnr g static run [gofile] [outpath] --
	recompile, newAppArgs := shouldRecompile(appArgs)
	if recompile {
		rebuildApp(outPath, outPath, pathToApp)
	}
	if !fileExists(outPath) {
		rebuildApp(outPath, outPath, pathToApp)
	}
	runExe(outPath, newAppArgs)
}
