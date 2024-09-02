package two

import (
	"fmt"
	"github.com/tsinghua-cel/strategy-gen/types"
)

func GenSlotStrategy(latestHackDutySlot int, epoch int64) []types.SlotStrategy {
	strategys := make([]types.SlotStrategy, 0)
	switch epoch % 3 {
	case 0, 1:
		strategy := types.SlotStrategy{
			Slot:    "every",
			Level:   0,
			Actions: make(map[string]string),
		}
		strategy.Actions["BlockBeforeSign"] = "return"
		strategy.Actions["AttestBeforeSign"] = fmt.Sprintf("return")
		strategys = append(strategys, strategy)

	case 2:
		{
			strategy := types.SlotStrategy{
				Slot:    "every",
				Level:   0,
				Actions: make(map[string]string),
			}
			strategy.Actions["BlockBeforeSign"] = "return"
			strategy.Actions["AttestAfterSign"] = fmt.Sprintf("addAttestToPool")
			strategy.Actions["AttestBeforeBroadCast"] = fmt.Sprintf("return")
			strategys = append(strategys, strategy)
		}
		{
			strategy := types.SlotStrategy{
				Slot:    fmt.Sprintf("%d", latestHackDutySlot),
				Level:   1,
				Actions: make(map[string]string),
			}
			strategy.Actions["AttestBeforeSign"] = fmt.Sprintf("return")

			strategy.Actions["BlockBeforeSign"] = "packPooledAttest"
			strategy.Actions["BlockDelayForReceiveBlock"] = fmt.Sprintf("delayHalfEpoch")
			// block delay to next epoch-end slot
			strategy.Actions["BlockBeforeBroadCast"] = fmt.Sprintf("delayHalfEpoch")

			strategys = append(strategys, strategy)
		}

	}
	return strategys

}
