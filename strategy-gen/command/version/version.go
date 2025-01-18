package version

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tsinghua-cel/strategy-gen/command"
	"github.com/tsinghua-cel/strategy-gen/command/helper"
	"github.com/tsinghua-cel/strategy-gen/version"
)

type VersionResult struct {
	Version string `json:"version"`
	Build   string `json:"build"`
}

func (r *VersionResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[VERSION INFO]\n")
	buffer.WriteString(helper.FormatKV([]string{
		fmt.Sprintf("Release version|%s\n", r.Version),
		fmt.Sprintf("Build version|%s\n", r.Build),
	}))

	return buffer.String()
}

func GetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Returns the current version",
		Args:  cobra.NoArgs,
		Run:   runCommand,
	}
}

func runCommand(cmd *cobra.Command, _ []string) {
	outputter := command.InitializeOutputter(cmd)
	defer outputter.WriteOutput()

	outputter.SetCommandResult(
		&VersionResult{
			Version: version.Version,
			Build:   version.Build,
		},
	)
}
