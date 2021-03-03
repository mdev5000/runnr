package runnr_test

import (
	"context"
	"fmt"
	"github.com/mdev5000/runnr"
	"github.com/spf13/cobra"
)

func Example_settingsUpApplication() {
	runner := runnr.NewRunner()

	runner.AddCommand(&cobra.Command{
		Use: "hello",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("hello world")
			return nil
		},
	})

	ctx := context.Background()
	if err := runner.Run(ctx); err != nil {
		panic(err)
	}
}

func GetCommands() []*cobra.Command {
	first := &cobra.Command{
		Use: "first",
	}
	second := &cobra.Command{
		Use: "second",
	}
	return []*cobra.Command{first, second}
}

func Example_registerThirdPartyCommands() {
	runner := runnr.NewRunner()

	// Register the third party commands under a root task called tp and exclude the "second" command.
	//
	// For example you would run the "first" command as follows:
	//   runnr tp first
	//
	runner.Register(GetCommands).UnderParent("tp").Exclude("second")

	ctx := context.Background()
	if err := runner.Run(ctx); err != nil {
		panic(err)
	}
}
