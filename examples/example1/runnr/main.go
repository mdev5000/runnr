package main

import (
	"github.com/mdev5000/runnr"
	"github.com/mdev5000/runnr/examples/example1/mycommands"
	"context"
)

func main() {
	ctx := context.Background()
	runner := runnr.NewRunner(ctx)
	runner.Register(mycommands.Commands())
	if err := runner.Run(); err != nil {
		panic(err)
	}
}
