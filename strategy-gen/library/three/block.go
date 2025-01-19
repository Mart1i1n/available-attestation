package three

import (
	"fmt"
	"github.com/tsinghua-cel/strategy-gen/types"
	"strconv"
)

func GenSlotStrategy(duties []interface{}) []types.SlotStrategy {
	strategys := make([]types.SlotStrategy, 0)
	duty := duties[0].([]types.ProposerDuty)
	slotStrategy := types.SlotStrategy{
		Slot:    fmt.Sprintf("%s", duty[0].Slot),
		Level:   1,
		Actions: make(map[string]string),
	}
	currentslot, _ := strconv.Atoi(duty[0].Slot)
	slotStrategy.Actions["BlockGetNewParentRoot"] = fmt.Sprintf("modifyParentRoot:%d", strconv.Itoa(currentslot-10))
	strategys = append(strategys, slotStrategy)

	return strategys

}
