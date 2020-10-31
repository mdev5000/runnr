# runnr

`runnr` is simple task runner for go (similar to packages like rake).

Basically a small wrapper around a `cobra` cli application.

# Getting started

```bash
go get bitbucket.org/mdev5000/runnr
go get bitbucket.org/mdev5000/runnr/runnr

# create a new runnr project
runnr -n
```

By default `runnr` creates a `runnr/main.go` and a `runnr.yml`.

# Registering commands 

Different module commands can be registered. You can also limit which
module commands get registered.

Also you can register a `cobra.Command` directly.

`runnr/main.go`
```go
package main

import (
	"bitbucket.org/mdev5000/runnr"
	"github.com/spf13/cobra"
	"somecommands"
	"somecommands2"
	"somecommands3"
	"context"
)

func main() {
	ctx := context.Background()
	runner := runnr.NewRunner(ctx)

	runner.Register(somecommands.Commands())
	runner.Register(somecommands2.Commands()).Only("first", "second")
	runner.Register(somecommands2.Commands()).Exclude("third", "fourth")

	runner.AddCommand(&cobra.Command{
		Use: "hello",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello world")
		},
	})

	if err := runner.Run(); err != nil {
		panic(err)
	}
}
```

Given the example above you could run the **hello** command like you would
any `cobra` subcommand.

```bash
runnr hello
```

If you wanted to recompile and then run it you could do:

```bash
runnr hello -r
```

For more examples of modules and packages registering commands see the
**example/** folder.

---

# Development

## Testing

```bash
cd examples/example1
go run ../../runnr/runnr.go -h
```