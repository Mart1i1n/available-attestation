package five

import (
	"github.com/tsinghua-cel/strategy-gen/utils"
	"strconv"
)


func CheckDuties(maxValidatorIndex int, duties []utils.ProposerDuty) ([]interface{}, bool) {
	result := make([]interface{}, 0)
	for i := 0; i < len(duties)-2; {
		a := duties[i]
		b := duties[i+1]
		c := duties[i+2]
		aValidatorIndex, _ := strconv.ParseInt(a.ValidatorIndex, 10, 64)
		bValidatorIndex, _ := strconv.ParseInt(b.ValidatorIndex, 10, 64)
		cValidatorIndex, _ := strconv.ParseInt(c.ValidatorIndex, 10, 64)
		if aValidatorIndex <= int64(maxValidatorIndex) && bValidatorIndex > int64(maxValidatorIndex) && cValidatorIndex <= int64(maxValidatorIndex) {
			result = append(result, []utils.ProposerDuty{a, b, c})
			i += 3
		} else {
			i++
		}
	}
	if len(result) > 0 {
		return result, true
	}
	return nil, false
}
