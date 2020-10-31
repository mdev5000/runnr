package runnr

import "github.com/spf13/cobra"

type filterFunc = func(name string) bool

type Registration interface {
	UnderParent(parentCmdName string) Registration
	Only(commandsNames ...string) Registration
	Exclude(commandsNames ...string) Registration
}

type registration struct {
	filterFuncs []filterFunc
	commands []*cobra.Command
	parentName string
}


var _ Registration = &registration{}


func newRegistration(cmds []*cobra.Command) *registration {
	return &registration{
		commands: cmds,
	}
}

func (r *registration) UnderParent(parentCmdName string) Registration {
	r.parentName = parentCmdName
	return r
}

func (r *registration) Only(commandsNames ...string) Registration {
	r.filterFuncs = append(r.filterFuncs, func(cmdName string) bool {
		for _, name := range commandsNames {
			if cmdName == name {
				return true
			}
		}
		return false
	})
	return r
}

func (r *registration) Exclude(commandsNames ...string) Registration {
	r.filterFuncs = append(r.filterFuncs, func(cmdName string) bool {
		for _, name := range commandsNames {
			if cmdName == name {
				return false
			}
		}
		return true
	})
	return r
}

func (r *registration) finalize(rootCmd *cobra.Command) {
	var rootToUse *cobra.Command
	if r.parentName != "" {
		rootToUse = &cobra.Command{
			Use: r.parentName,
		}
		rootCmd.AddCommand(rootToUse)
	} else {
		rootToUse = rootCmd
	}
	for _, cmd := range r.commands {
		if r.passesFilters(cmd.Name()) {
			rootToUse.AddCommand(cmd)
		}
	}
}

func (r *registration) passesFilters(cmdName string) bool {
	for _, filter := range r.filterFuncs {
		if !filter(cmdName) {
			return false
		}
	}
	return true
}
