package cli

import (
	"fmt"
	"github.com/blang/vfs"
	"github.com/mdev5000/runnr"
	"github.com/spf13/cobra"
	"os"
)

func checkForGlobal(fs vfs.Filesystem, workingDir string) (bool, error) {
	if len(os.Args) < 2 {
		return false, nil
	}
	if !(os.Args[1] == "g" || os.Args[1] == "global") {
		return false, nil
	}
	err := runGlobal(fs, workingDir)
	return true, err
}

func runGlobal(fs vfs.Filesystem, workingDir string) error {
	cmdGlobalRoot := cobra.Command{
		Use:     "global",
		Short:   "Global runnr commands",
		Aliases: []string{"g"},
	}

	newCmd := &cobra.Command{
		Use:   "new",
		Short: "generate a new local project",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Initializing...")
			settings := runnr.NewProjectSettings{}
			if err := runnr.SetupNewProject(fs, workingDir, settings); err != nil {
				return err
			}
			fmt.Println("Done.")
			fmt.Println("You can get started by running:")
			fmt.Println("")
			fmt.Println("  runnr hello")
			fmt.Println("")
			return nil
		},
	}

	cmdGlobalRoot.AddCommand(newCmd)
	s := &staticCmd{}
	cmdGlobalRoot.AddCommand(s.Cmd())

	rootCmd := cobra.Command{
		Use: "runnr",
	}
	rootCmd.AddCommand(&cmdGlobalRoot)

	return cmdGlobalRoot.Execute()
}

type staticCmd struct{}

func (s *staticCmd) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "static",
		Short:   "Commands for running static projects",
		Aliases: []string{"s"},
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "run [gofile] [outpath] --",
		Short: "Run a static command",
		Args:  cobra.MinimumNArgs(2),
		RunE:  s.runStatic,
	})
	return cmd
}

func (s *staticCmd) runStatic(cmd *cobra.Command, args []string) error {
	pathToApp := args[0]
	outPath := args[1]
	appArgs := os.Args[6:] // ignore the: runnr g static run [gofile] [outpath] --
	recompile, newAppArgs := shouldRecompile(appArgs)
	if recompile {
		if err := rebuildApp(outPath, outPath, pathToApp); err != nil {
			return err
		}
	}
	if !fileExists(outPath) {
		if err := rebuildApp(outPath, outPath, pathToApp); err != nil {
			return err
		}
	}
	return runnr.RunExe(outPath, newAppArgs)
}
