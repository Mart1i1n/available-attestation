package three

import (
	"github.com/tsinghua-cel/strategy-gen/types"
)

func ValidatorStrategy(maxValidatorIndex int, epoch int64) []types.ValidatorStrategy {
	//tool := utils.SlotTool{SlotsPerEpoch: 32}
	//begin := tool.EpochStart(epoch)
	//end := tool.EpochEnd(epoch)
	strategy := make([]types.ValidatorStrategy, 0)
	for idx := 0; idx <= maxValidatorIndex; idx++ {
		strategy = append(strategy, types.ValidatorStrategy{
			ValidatorIndex:    idx,
			AttackerStartSlot: int(0),
			AttackerEndSlot:   int(1000),
		})
	}
	return strategy
}
