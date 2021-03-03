# Simple task runner

`runnr` is simple task runner for go (similar to packages like rake). It's a small wrapper around a `cobra` cli
application, to create more flexible application build commands.

## Getting started

```bash
go get github.com/mdev5000/runnr
go get github.com/mdev5000/runnr/runnr

# Create a new runnr project.
runnr g new 

# Run the example.
runnr hello
```

By default `runnr` creates a local exe at `internal/cmd/runnr_local/main.go`. You can move and configure this in
`runnr.yml`.

## Registering commands 

Different module commands can be registered. You can also limit which
module commands get registered.

Also you can register a `cobra.Command` directly.

`runnr/main.go`
```go
package main

import (
	"github.com/mdev5000/runnr"
	"github.com/spf13/cobra"
	"somecommands"
	"somecommands2"
	"somecommands3"
	"context"
)

func main() {
	runner := runnr.NewRunner()

	runner.Register(somecommands.Commands)
	runner.Register(somecommands2.Commands).Only("first", "second")
	runner.Register(somecommands2.Commands).Exclude("third", "fourth")

	runner.AddCommand(&cobra.Command{
		Use: "hello",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello world")
		},
	})

	ctx := context.Background()
	if err := runner.Run(ctx); err != nil {
		panic(err)
	}
}
```

Given the example above you could run the **hello** command like you would
any `cobra` subcommand.

```bash
runnr hello
```

You can recompile at anytime by appending `-r` after the command.

```bash
runnr hello -r
```

For more examples of modules and packages registering commands see the
**example/** folder.

---

## Running static commands

You can run a static project using the static command

```bash
alias myapp="runnr g s run ~/path/to/main.go ~/path/to/myapp --"
```

Then you can run `myapp` like a `runnr` command:

```bash
# Recompile your application.
myapp -r

# Show the help for your application (if it exists).
myapp -h
```