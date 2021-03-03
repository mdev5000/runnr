package runnr

import (
	"context"
	"github.com/spf13/cobra"
)

type Runner struct {
	InstanceName  string
	rootCmd       *cobra.Command
	registrations []*registration
}

// Interface implemented by third parties to register tasks (tasks are basically just cobra comands).
type CommandRegisterer = func() []*cobra.Command

// Create a new runnr application.
func NewRunner() *Runner {
	runnrCmd := &cobra.Command{
		Use: "runnr",
	}
	runnrCmd.PersistentFlags().BoolP("recompile", "r", false, "Recompile the tasks.")
	runnrCmd.PersistentFlags().BoolP("new", "n", false, "Create a instance")
	return &Runner{
		InstanceName:  "default",
		rootCmd:       runnrCmd,
		registrations: nil,
	}
}

func (r *Runner) AddCommand(cmd *cobra.Command) {
	r.rootCmd.AddCommand(cmd)
}

func (r *Runner) Register(registration CommandRegisterer) Registration {
	reg := newRegistration(registration())
	r.registrations = append(r.registrations, reg)
	return reg
}

func (r *Runner) ToCommand() *cobra.Command {
	r.finalizeRegistration()
	return r.rootCmd
}

func (r *Runner) Run(ctx context.Context) error {
	return r.ToCommand().ExecuteContext(ctx)
}

func (r *Runner) finalizeRegistration() {
	for _, reg := range r.registrations {
		reg.finalize(r.rootCmd)
	}
}
