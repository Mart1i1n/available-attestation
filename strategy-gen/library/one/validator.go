package one

import (
	"github.com/tsinghua-cel/strategy-gen/types"
	"github.com/tsinghua-cel/strategy-gen/utils"
	"strconv"
)

func ValidatorStrategy(allHacks []interface{}) []types.ValidatorStrategy {
	strategy := make([]types.ValidatorStrategy, 0)
	for _, subduties := range allHacks {
		duties := subduties.([]utils.ProposerDuty)

		begin, _ := strconv.Atoi(duties[0].Slot)
		end, _ := strconv.Atoi(duties[len(duties)-1].Slot)

		for i := 0; i < len(duties); i++ {
			idx, _ := strconv.Atoi(duties[i].ValidatorIndex)
			strategy = append(strategy, types.ValidatorStrategy{
				ValidatorIndex:    idx,
				AttackerStartSlot: begin,
				AttackerEndSlot:   end,
			})
		}
	}

	return strategy
}
