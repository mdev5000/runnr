package main

import (
	"bitbucket.org/mdev5000/runnr"
	"bitbucket.org/mdev5000/runnr/examples/example1/mycommands"
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
