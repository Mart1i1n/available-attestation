package four

import (
	"github.com/tsinghua-cel/strategy-gen/types"
	"strconv"
)

func CheckDuties(maxValidatorIndex int, duties []types.ProposerDuty) ([]interface{}, bool) {
	result := make([]interface{}, 0)
	endSlot, _ := strconv.Atoi(duties[len(duties)-1].Slot)

	tmpsub := make([]types.ProposerDuty, 0)
	for _, duty := range duties {
		valIdx, _ := strconv.Atoi(duty.ValidatorIndex)

		if valIdx <= maxValidatorIndex {
			tmpsub = append(tmpsub, duty)
		} else {
			if len(tmpsub) > 9 {
				startSlot, _ := strconv.Atoi(tmpsub[0].Slot)
				if (endSlot - startSlot) <= 15 {
					result = append(result, tmpsub)
				}
			}
			tmpsub = make([]types.ProposerDuty, 0)
		}
	}
	if len(tmpsub) > 9 {
		startSlot, _ := strconv.Atoi(tmpsub[0].Slot)
		if (endSlot - startSlot) <= 15 {
			result = append(result, tmpsub)
		}
	}

	if len(result) > 0 {
		return result, true
	}

	return nil, false
}
