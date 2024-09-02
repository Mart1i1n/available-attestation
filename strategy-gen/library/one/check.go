package one

import (
	"github.com/tsinghua-cel/strategy-gen/utils"
	"strconv"
)

func CheckDuties(maxValidatorIndex int, duties []utils.ProposerDuty) ([]interface{}, bool) {
	result := make([]interface{}, 0)

	tmpsub := make([]utils.ProposerDuty, 0)
	for _, duty := range duties {
		valIdx, _ := strconv.Atoi(duty.ValidatorIndex)

		if valIdx <= maxValidatorIndex {
			tmpsub = append(tmpsub, duty)
		} else {
			if len(tmpsub) > 2 {
				result = append(result, tmpsub)
			}
			tmpsub = make([]utils.ProposerDuty, 0)
		}
	}
	if len(tmpsub) > 2 {
		result = append(result, tmpsub)
	}

	if len(result) > 0 {
		return result, true
	}

	return nil, false
}
