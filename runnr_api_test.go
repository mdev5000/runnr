package runnr_test

import (
	"bytes"
	"fmt"
	"github.com/mdev5000/runnr"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"testing"
)

func myCommands() runnr.CommandRegisterer {
	return func() []*cobra.Command {
		return []*cobra.Command{
			{
				Use:  "first",
				Args: cobra.MinimumNArgs(1),
				RunE: func(cmd *cobra.Command, args []string) error {
					out := cmd.OutOrStdout()
					fmt.Fprintln(out, args[0])
					return nil
				},
			},
		}
	}
}

func TestRunnr_CanRun(t *testing.T) {
	out := bytes.NewBufferString("")
	r := runnr.NewRunner()
	r.Register(myCommands())
	cmd := r.ToCommand()
	cmd.SetArgs([]string{"first", "hello"})
	cmd.SetOut(out)
	require.Nil(t, cmd.Execute())
	require.Equal(t, out.String(), "hello\n")
}
