package three

import (
	"github.com/tsinghua-cel/strategy-gen/types"
	"strconv"
)

func CheckDuties(maxValidatorIndex int, duties []types.ProposerDuty) ([]interface{}, bool) {
	result := make([]interface{}, 0)

	tmpsub := make([]types.ProposerDuty, 0)

	firstproposerindex, _ := strconv.Atoi(duties[0].ValidatorIndex)
	if firstproposerindex > maxValidatorIndex {
		return nil, false
	}
	duty := duties[0]
	tmpsub = append(tmpsub, duty)
	result = append(result, tmpsub)
	return result, true
}
