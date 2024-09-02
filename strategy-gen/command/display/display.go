package display

import (
	"github.com/spf13/cobra"
	"github.com/tsinghua-cel/strategy-gen/actionset"
	"github.com/tsinghua-cel/strategy-gen/command"
	"github.com/tsinghua-cel/strategy-gen/pointset"
)

func GetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "display",
		Short: "Display all action and action-point information",
		Args:  cobra.NoArgs,
		Run:   runCommand,
	}
}

func runCommand(cmd *cobra.Command, _ []string) {
	outputter := command.InitializeOutputter(cmd)
	defer outputter.WriteOutput()

	outputter.SetCommandResult(DisplayResult{
		ActionPointResult: ActionPointResult{
			BlockPoints: pointset.BlockPointSet,
			AttPoints:   pointset.AttestPointSet,
		},
		ActionResult: ActionResult{
			BlockActionList: actionset.GetBlockActionNameList(),
			AttActionList:   actionset.GetAttestActionNameList(),
		},
	})
}
