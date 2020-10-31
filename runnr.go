package runnr

import (
	"context"
	"github.com/spf13/cobra"
)

type Runner struct {
	InstanceName  string
	rootCmd       *cobra.Command
	registrations []*registration
	ctx           context.Context
}

type CommandRegisterer interface {
	GetCommands(ctx context.Context) []*cobra.Command
}

func NewRunner(ctx context.Context) *Runner {
	runnrCmd := &cobra.Command{
		Use: "runnr",
	}
	runnrCmd.PersistentFlags().BoolP("recompile", "r", false, "Recompile the tasks.")
	runnrCmd.PersistentFlags().BoolP("new", "n", false, "Create a instance")
	return &Runner{
		InstanceName: "default",
		rootCmd: runnrCmd,
		registrations: nil,
		ctx:           nil,
	}
}

func (r *Runner) Register(c CommandRegisterer) Registration {
	reg := newRegistration(c.GetCommands(r.ctx))
	r.registrations = append(r.registrations, reg)
	return reg
}

func (r *Runner) Run() error {
	r.finalizeRegistration()
	return r.rootCmd.ExecuteContext(r.ctx)
}

func (r *Runner) finalizeRegistration() {
	for _, reg := range r.registrations {
		reg.finalize(r.rootCmd)
	}
}
