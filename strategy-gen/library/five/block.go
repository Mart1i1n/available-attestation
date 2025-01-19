package five

import (
	"fmt"
	"github.com/tsinghua-cel/strategy-gen/types"
)

/*
sandwich reorg attack
*/

func GenSlotStrategy(duties []interface{}) []types.SlotStrategy {
	strategys := make([]types.SlotStrategy, 0)
	for i := 0; i < len(duties); i++ {
		duty := duties[i].([]types.ProposerDuty)
		if len(duty) != 3 {
			continue
		}
		a := duty[0]
		//b := duty[1]
		c := duty[2]
		{
			slotStrategy := types.SlotStrategy{
				Slot:    fmt.Sprintf("%s", a.Slot),
				Level:   1,
				Actions: make(map[string]string),
			}
			slotStrategy.Actions["BlockBeforeBroadCast"] = fmt.Sprintf("%s:%d", "delayWithSecond", 2*12)
			strategys = append(strategys, slotStrategy)
		}
		{
			slotStrategy := types.SlotStrategy{
				Slot:    fmt.Sprintf("%s", c.Slot),
				Level:   1,
				Actions: make(map[string]string),
			}
			slotStrategy.Actions["BlockGetNewParentRoot"] = fmt.Sprintf("modifyParentRoot:%d", a.Slot)
			strategys = append(strategys, slotStrategy)
		}

	}

	return strategys

}
