package mycommands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func Commands() []*cobra.Command {
	something := &cobra.Command{
		Use:   "something",
		Short: "Does something pretty cool",
		RunE:  something,
	}

	return []*cobra.Command{something}
}

func something(cmd *cobra.Command, args []string) error {
	fmt.Println("something")
	return nil
}
