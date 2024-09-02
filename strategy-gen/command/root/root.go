package root

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tsinghua-cel/strategy-gen/command/display"
	"github.com/tsinghua-cel/strategy-gen/command/generate"
	"github.com/tsinghua-cel/strategy-gen/command/helper"
	"github.com/tsinghua-cel/strategy-gen/command/runtime"
	"github.com/tsinghua-cel/strategy-gen/command/update"
	"github.com/tsinghua-cel/strategy-gen/command/version"
	"os"
)

type RootCommand struct {
	baseCmd *cobra.Command
}

func NewRootCommand() *RootCommand {
	rootCommand := &RootCommand{
		baseCmd: &cobra.Command{
			Short: "Strategy-gen is a tool to generate strategy file for attacker service",
		},
	}

	helper.RegisterJSONOutputFlag(rootCommand.baseCmd)
	rootCommand.baseCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCommand.registerSubCommands()

	return rootCommand
}

func (rc *RootCommand) registerSubCommands() {
	rc.baseCmd.AddCommand(
		runtime.GetCommand(),
		generate.GetCommand(),
		update.GetCommand(),
		display.GetCommand(),
		version.GetCommand(),
	)
}

func (rc *RootCommand) Execute() {
	if err := rc.baseCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
