package four

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
			// and then delay 1.5 epoch.
			totalSlots := 32 / 2 * 3
			totalSeconds := 12 * totalSlots
			strategy.Actions["BlockDelayForReceiveBlock"] = fmt.Sprintf("delayWithSecond:%d", totalSeconds)

			// delay half epoch first.
			strategy.Actions["BlockBeforeBroadCast"] = fmt.Sprintf("delayHalfEpoch")

			strategys = append(strategys, strategy)
		}

	}
	return strategys

}
