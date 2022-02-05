package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

type TimeTracker interface {
	StopTracking() (string, error)
	StartTracking(key string) (string, error)
	Export(format string) (string, error)
}

var timeTracker TimeTracker
var rootCmd *cobra.Command

func Execute(tracker TimeTracker) {
	timeTracker = tracker

	cmd, _, err := rootCmd.Find(os.Args[1:])

	if err != nil && cmd.Use == rootCmd.Use && cmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp {
		args := append([]string{"start"}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func stop(_ *cobra.Command, _ []string) (string, error) {
	return timeTracker.StopTracking()
}

func export(cmd *cobra.Command, _ []string) (string, error) {
	return timeTracker.Export(cmd.Flag("format").Value.String())
}

func start(_ *cobra.Command, args []string) (string, error) {
	if len(args) == 0 {
		return timeTracker.StartTracking("default")
	}
	return timeTracker.StartTracking(strings.Join(args, " "))
}

func init() {
	rootCmd = &cobra.Command{
		Use:   "tk",
		Short: "timekeeper is a small app to help you keeping track of the time you spend on your projects",
	}

	rootCmd.AddCommand(newCmd("stop", "stops time tracking", stop, nil))
	rootCmd.AddCommand(newCmd("start", "start time tracking with specified key", start, nil))
	rootCmd.AddCommand(newCmd("export", "export timetable", export, []stringFlag{
		{"format", "f", "json", "output format [csv, json]"},
	}))
}

type stringFlag struct {
	name  string
	short string
	value string
	usage string
}

func newCmd(use, short string, run func(_ *cobra.Command, args []string) (string, error), flags []stringFlag) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := run(cmd, args)
			if err != nil {
				return err
			}
			fmt.Println(out)
			return nil
		},
	}
	for _, flag := range flags {
		cmd.Flags().StringP(flag.name, flag.short, flag.value, flag.usage)
	}
	return cmd
}
