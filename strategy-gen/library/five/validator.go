package five

import (
	"github.com/tsinghua-cel/strategy-gen/types"
	"strconv"
)

func ValidatorStrategy(duties []types.ProposerDuty) []types.ValidatorStrategy {
	begin, _ := strconv.Atoi(duties[0].Slot)
	end, _ := strconv.Atoi(duties[len(duties)-1].Slot)
	strategy := make([]types.ValidatorStrategy, 0)
	for i := 0; i < len(duties); i++ {
		idx, _ := strconv.Atoi(duties[i].ValidatorIndex)
		strategy = append(strategy, types.ValidatorStrategy{
			ValidatorIndex:    idx,
			AttackerStartSlot: begin,
			AttackerEndSlot:   end,
		})
	}
	return strategy
}
