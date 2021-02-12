package runnr_test

import (
	"bytes"
	"github.com/mdev5000/runnr"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func runCommand(out io.Writer, cmd string, args ...string) error {
	pre := exec.Command(cmd, args...)
	pre.Stdin = os.Stdin
	pre.Stdout = out
	pre.Stderr = os.Stderr
	if err := pre.Start(); err != nil {
		return err
	}
	if err := pre.Wait(); err != nil {
		return err
	}
	return nil
}

func runRunnrCommand(t *testing.T, relativePathToRunnr string, runnrArgs ...string) string {
	var runnrExe string
	if relativePathToRunnr != "" {
		runnrExe = filepath.Join(relativePathToRunnr, "cmd", "runnr", "runnr.go")
	} else {
		runnrExe = "./" + filepath.Join("cmd", "runnr", "runnr.go")
	}
	b := bytes.NewBufferString("")
	args := []string{"run", runnrExe}
	for _, arg := range runnrArgs {
		args = append(args, arg)
	}
	require.Nil(t, runCommand(b, "go", args...))
	return b.String()
}

func resetAndMoveIntoTmpDir(t *testing.T) string {
	runnr.RunCommand("rm", "-rf", "_tmp")
	require.Nil(t, runnr.RunCommand("mkdir", "_tmp"))
	os.Chdir("_tmp")
	wd, err := os.Getwd()
	require.Nil(t, err)
	return wd
}

func Test_SettingUpAndRunningNewProjects(t *testing.T) {
	runnrPath := ".."
	resetAndMoveIntoTmpDir(t)

	// Test initialization works.
	output := runRunnrCommand(t, runnrPath, "g", "new")
	require.Equal(t, output, strings.TrimSpace(`
Initializing...
Done.
You can get started by running:

  runnr hello
`)+"\n\n")

	t.Run("help command works", func(t *testing.T) {
		output := runRunnrCommand(t, runnrPath, "-h")
		require.Equal(t, output, strings.TrimSpace(`
Recompiling....
Recompiling done.
Usage:
  runnr [command]

Available Commands:
  hello       
  help        Help about any command

Flags:
  -h, --help        help for runnr
  -n, --new         Create a instance
  -r, --recompile   Recompile the tasks.

Use "runnr [command] --help" for more information about a command.
`)+"\n")
	})

	t.Run("recompiles only when -r flat is appended to the end (app must still exists though)", func(t *testing.T) {
		outputRecompile := runRunnrCommand(t, runnrPath, "hello", "-r")
		require.Equal(t, outputRecompile, strings.TrimSpace(`
Recompiling....
Recompiling done.
hello world
`)+"\n")

		output := runRunnrCommand(t, runnrPath, "hello")
		require.Equal(t, output, strings.TrimSpace(`
hello world
`)+"\n")
	})

	// Next set of tests run from root in of the _tmp dir.

	os.Chdir("..")

	t.Run("can run statically", func(t *testing.T) {
		staticPath := "./_tmp/internal/cmd/runnr_local/main.go"
		staticExe := "./_tmp/statictest"
		output := runRunnrCommand(t, "",
			"g", "s", "run", staticPath, staticExe, "--", "hello", "-r")
		require.Equal(t, output, strings.TrimSpace(`
Recompiling....
Recompiling done.
hello world
`)+"\n")
	})
}

func Test_canRunExample1(t *testing.T) {
	require.Nil(t, os.Chdir("examples/example1"))
	output := runRunnrCommand(t, "../../", "something", "-r")
	require.Equal(t, output, strings.TrimSpace(`
Recompiling....
Recompiling done.
working dir:  /Users/matt/devtmp/go/runnr/examples/example1
something
`)+"\n")
}
