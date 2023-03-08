package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"runtime"
	"strings"
)

type Command struct {
	usage    string
	desc     string
	options  CliOptions
	commands []*Command
	runFunc  RunCommandFunc
}

type CommandOption func(*Command)

type RunCommandFunc func(args []string) error

func WithCommandOptions(opt CliOptions) CommandOption {
	return func(command *Command) {
		command.options = opt
	}
}

func WithCommandRunFunc(run RunCommandFunc) CommandOption {
	return func(command *Command) {
		command.runFunc = run
	}
}

func NewCommand(usage string, desc string, opts ...CommandOption) *Command {
	cmd := &Command{
		usage: usage,
		desc:  desc,
	}

	for _, opt := range opts {
		opt(cmd)
	}

	return cmd
}

// AddCommand adds sub command to the current command.
func (c *Command) AddCommand(cmd *Command) {
	c.commands = append(c.commands, cmd)
}

// AddCommands adds multiple sub commands to the current command.
func (c *Command) AddCommands(cmds ...*Command) {
	c.commands = append(c.commands, cmds...)
}

func (c *Command) runCommand(_ *cobra.Command, args []string) {
	if c.runFunc != nil {
		if err := c.runFunc(args); err != nil {
			fmt.Printf("%v %v\n", color.RedString("Error:"), err)
			os.Exit(1)
		}
	}
}

func (c *Command) cobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.usage,
		Short: c.desc,
	}
	cmd.SetOut(os.Stdout)
	cmd.Flags().SortFlags = false
	if len(c.commands) > 0 {
		for _, command := range c.commands {
			cmd.AddCommand(command.cobraCommand())
		}
	}
	if c.runFunc != nil {
		cmd.Run = c.runCommand
	}
	if c.options != nil {
		for _, f := range c.options.Flags().FlagSets {
			cmd.Flags().AddFlagSet(f)
		}
		// c.options.AddFlags(cmd.Flags())
	}
	addHelpCommandFlag(c.usage, cmd.Flags())

	return cmd
}

func FormatBaseName(baseName string) string {
	if runtime.GOOS == "windows" {
		baseName = strings.ToLower(baseName)
		baseName = strings.TrimSuffix(baseName, ".exe")
	}
	return baseName
}
