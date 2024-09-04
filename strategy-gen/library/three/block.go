package three

import (
	"fmt"
	"github.com/tsinghua-cel/strategy-gen/types"
	"github.com/tsinghua-cel/strategy-gen/utils"
	"strconv"
)

func GenSlotStrategy(duties []interface{}) []types.SlotStrategy {
	strategys := make([]types.SlotStrategy, 0)
	duty := duties[0].([]utils.ProposerDuty)
	slotStrategy := types.SlotStrategy{
		Slot:    fmt.Sprintf("%d", duty[0].Slot),
		Level:   1,
		Actions: make(map[string]string),
	}
	currentslot, _ := strconv.Atoi(duty[0].Slot)
	slotStrategy.Actions["BlockGetNewParentRoot"] = fmt.Sprintf("modifyParentRoot:%d", strconv.Itoa(currentslot-10))
	strategys = append(strategys, slotStrategy)

	return strategys

}
