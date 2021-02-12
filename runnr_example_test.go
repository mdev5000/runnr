package runnr_test

import (
	"context"
	"fmt"
	"github.com/mdev5000/runnr"
	"github.com/spf13/cobra"
)

func Example_settingsUpApplication() {
	ctx := context.Background()
	runner := runnr.NewRunner(ctx)

	runner.AddCommand(&cobra.Command{
		Use: "hello",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("hello world")
			return nil
		},
	})

	if err := runner.Run(); err != nil {
		panic(err)
	}
}

type ThirdParty struct{}

func (t *ThirdParty) GetCommands(ctx context.Context) []*cobra.Command {
	first := &cobra.Command{
		Use: "first",
	}
	second := &cobra.Command{
		Use: "second",
	}
	return []*cobra.Command{first, second}
}
func NewThirdParty() runnr.CommandRegisterer {
	return &ThirdParty{}
}

func Example_registerThirdPartyCommands() {
	ctx := context.Background()
	runner := runnr.NewRunner(ctx)

	// Register the third party commands under a root task called tp and exclude the "second" command.
	//
	// For example you would run the "first" command as follows:
	//   runnr tp first
	//
	runner.Register(NewThirdParty()).UnderParent("tp").Exclude("second")

	if err := runner.Run(); err != nil {
		panic(err)
	}
}
