package runnr

import (
	"fmt"
	"github.com/blang/vfs"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"strings"
)

type NewProjectSettings struct {
	LocalLocation string
}

func SetupNewProject(fs vfs.Filesystem, workingDir string, settings NewProjectSettings) error {
	configFile := filepath.Join(workingDir, ConfigFileName)
	if fileExists(fs, configFile) {
		return fmt.Errorf("cannot initialize config file '%s' already exists", ConfigFileName)
	}
	runnrLocalApp := filepath.Join("internal", "cmd", "runnr_local", "main.go")
	if settings.LocalLocation != "" {
		runnrLocalApp = settings.LocalLocation
	}

	// @todo add cleaning
	runnerLocalAppFull := filepath.Join(workingDir, runnrLocalApp)
	//if err != nil {
	//	return fmt.Errorf("runnr local path is invalid ('%s'), %w", runnrLocalApp, err)
	//}
	if !strings.HasPrefix(runnerLocalAppFull, workingDir) {
		return fmt.Errorf("cannot initialize runnr outside of the working directory")
	}
	runnrLocalAppFull := filepath.Join(workingDir, runnrLocalApp)
	if fileExists(fs, runnrLocalAppFull) {
		panic(fmt.Errorf("cannot initialize app already exists '%s'", runnrLocalApp))
	}
	runnrLocalAppDirFull := filepath.Dir(runnerLocalAppFull)
	if err := vfs.MkdirAll(fs, runnrLocalAppDirFull, 0775); err != nil {
		return err
	}

	mainFileContents := `
package main

import (
	"context"
	"fmt"
	"github.com/mdev5000/runnr"
	"github.com/spf13/cobra"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	ctx := context.Background()
	runner := runnr.NewRunner(ctx)

	runner.AddCommand(&cobra.Command{
		Use: "hello",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello world")
		},
	})

	return runner.Run()
}
`
	if err := vfs.WriteFile(fs, runnerLocalAppFull, []byte(mainFileContents), 0775); err != nil {
		return err
	}
	out, err := yaml.Marshal(&Config{
		Outpath:    "./.tmpbins/runnr-local",
		Executable: "./" + runnrLocalApp,
	})
	if err != nil {
		return err
	}
	return vfs.WriteFile(fs, configFile, out, 0644)
}
