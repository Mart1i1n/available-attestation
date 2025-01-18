package one

import (
	"fmt"
	"github.com/tsinghua-cel/strategy-gen/globalinfo"
	"github.com/tsinghua-cel/strategy-gen/pointset"
	"github.com/tsinghua-cel/strategy-gen/types"
	"strconv"
)

func BlockStrategy(cur, end int, actions map[string]string) {
	secondPerSlot := globalinfo.ChainBaseInfo().SecondsPerSlot
	point := pointset.GetPointByName("BlockBeforeBroadCast")
	actions[point] = fmt.Sprintf("%s:%d", "delayWithSecond", (end+1-cur)*secondPerSlot)
}

func GenSlotStrategy(allHacks []interface{}) []types.SlotStrategy {
	strategys := make([]types.SlotStrategy, 0)
	for _, subduties := range allHacks {
		duties := subduties.([]types.ProposerDuty)
		//begin, _ := strconv.Atoi(duties[0].Slot)
		end, _ := strconv.Atoi(duties[len(duties)-1].Slot)

		for i := 0; i < len(duties); i++ {
			slot, _ := strconv.Atoi(duties[i].Slot)
			//idx, _ := strconv.Atoi(duties[i].ValidatorIndex)
			strategy := types.SlotStrategy{
				Slot:    duties[i].Slot,
				Level:   0,
				Actions: make(map[string]string),
			}
			BlockStrategy(slot, end, strategy.Actions)
			strategys = append(strategys, strategy)
		}
	}

	return strategys

}
