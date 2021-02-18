package mycommands

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
)

type MyCommands struct {
	ctx context.Context
}

func Commands() *MyCommands {
	return &MyCommands{}
}

func (c *MyCommands) GetCommands(ctx context.Context) []*cobra.Command {
	c.ctx = ctx

	something := &cobra.Command{
		Use: "something",
		Short: "Does something pretty cool",
		Run: c.something,
	}

	return []*cobra.Command {something}
}

func (c *MyCommands) something(cmd *cobra.Command, args []string) {
	fmt.Println("something")
}
