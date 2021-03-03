package main

import (
	"context"
	"github.com/mdev5000/runnr"
	"github.com/mdev5000/runnr/examples/example1/mycommands"
)

func main() {
	runner := runnr.NewRunner()
	runner.Register(mycommands.Commands)
	ctx := context.Background()
	if err := runner.Run(ctx); err != nil {
		panic(err)
	}
}
