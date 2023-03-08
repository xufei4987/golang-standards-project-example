package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	my_error "golang-standards-project-example/pkg/errors"
	"golang-standards-project-example/pkg/term"
	"golang-standards-project-example/pkg/version"
	"golang-standards-project-example/pkg/version/verflag"
	"log"
	"os"
)

var progressMessage = color.GreenString("==>")

type App struct {
	name        string
	basename    string
	description string
	options     CliOptions
	noVersion   bool
	noConfig    bool
	silence     bool
	runFunc     RunFunc
	commands    []*Command //子命令
	cmd         *cobra.Command
}

type Option func(*App)

func WithOptions(opt CliOptions) Option {
	return func(app *App) {
		app.options = opt
	}
}

type RunFunc func(basename string) error

func WithRunFunc(run RunFunc) Option {
	return func(app *App) {
		app.runFunc = run
	}
}

func WithDescription(desc string) Option {
	return func(app *App) {
		app.description = desc
	}
}

func WithNoVersion(noVersion bool) Option {
	return func(app *App) {
		app.noVersion = noVersion
	}
}

func WithNoConfig(noConfig bool) Option {
	return func(app *App) {
		app.noConfig = noConfig
	}
}

func NewApp(name, baseName string, opts ...Option) *App {
	app := &App{
		name:     name,
		basename: baseName,
	}
	for _, opt := range opts {
		opt(app)
	}
	app.buildCommond()
	return app
}

func (a *App) buildCommond() {
	cmd := cobra.Command{
		Use:           FormatBaseName(a.basename),
		Short:         a.name,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true
	InitFlags(cmd.Flags())
	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		cmd.SetHelpCommand(helpCommand(FormatBaseName(a.basename)))
	}
	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}
	var namedFlagSets NamedFlagSets
	//添加自定义flag
	if a.options != nil {
		namedFlagSets = a.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f)
		}
	}
	if !a.noVersion {
		verflag.AddFlags(namedFlagSets.FlagSet("global"))
	}
	//添加全局flag
	if !a.noConfig {
		addConfigFlag(a.basename, namedFlagSets.FlagSet("global"))
	}
	AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name())
	// add new global flagset to cmd FlagSet
	cmd.Flags().AddFlagSet(namedFlagSets.FlagSet("global"))

	addCmdTemplate(&cmd, namedFlagSets)
	a.cmd = &cmd
}

// Run is used to launch the application.
func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	printWorkingDir()
	PrintFlags(cmd.Flags())
	if !a.noVersion {
		// display application version information
		verflag.PrintAndExitIfRequested()
	}

	if !a.noConfig {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}

		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}

	if !a.silence {
		log.Printf("%v Starting %s ...\n", progressMessage, a.name)
		if !a.noVersion {
			log.Printf("%v Version: `%s`\n", progressMessage, version.Get().ToJSON())
		}
		if !a.noConfig {
			log.Printf("%v Config file used: `%s`\n", progressMessage, viper.ConfigFileUsed())
		}
	}
	if a.options != nil {
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}
	// run application
	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}

	return nil
}

func (a *App) applyOptionRules() error {
	if completeableOptions, ok := a.options.(CompleteableOptions); ok {
		if err := completeableOptions.Complete(); err != nil {
			return err
		}
	}

	if errs := a.options.Validate(); len(errs) != 0 {
		return my_error.NewAggregate(errs)
	}

	if printableOptions, ok := a.options.(PrintableOptions); ok && !a.silence {
		log.Printf("%v Config: `%s`\n", progressMessage, printableOptions.String())
	}

	return nil
}

func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Printf("%v WorkingDir: %s\n", progressMessage, wd)
}

func addCmdTemplate(cmd *cobra.Command, namedFlagSets NamedFlagSets) {
	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})
}
