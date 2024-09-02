package config

import "github.com/tsinghua-cel/strategy-gen/actionset"

var (
	defaultActionConfigs = make([]actionset.ActionConfig, 0)
)

func defaultConfigInit() {
	if len(defaultActionConfigs) > 0 {
		return
	}

	defaultActionConfigs = append(defaultActionConfigs, actionset.ActionConfig{
		Name:                "null",
		Random:              false,
		ParamCount:          0,
		DefaultParamValue:   0,
		MinRandomParamValue: 0,
		MaxRandomValue:      0,
	})
	defaultActionConfigs = append(defaultActionConfigs, actionset.ActionConfig{
		Name:                "return",
		Random:              false,
		ParamCount:          0,
		DefaultParamValue:   0,
		MinRandomParamValue: 0,
		MaxRandomValue:      0,
	})
	defaultActionConfigs = append(defaultActionConfigs, actionset.ActionConfig{
		Name:                "continue",
		Random:              false,
		ParamCount:          0,
		DefaultParamValue:   0,
		MinRandomParamValue: 0,
		MaxRandomValue:      0,
	})
	defaultActionConfigs = append(defaultActionConfigs, actionset.ActionConfig{
		Name:                "abort",
		Random:              false,
		ParamCount:          0,
		DefaultParamValue:   0,
		MinRandomParamValue: 0,
		MaxRandomValue:      0,
	})
	defaultActionConfigs = append(defaultActionConfigs, actionset.ActionConfig{
		Name:                "skip",
		Random:              false,
		ParamCount:          0,
		DefaultParamValue:   0,
		MinRandomParamValue: 0,
		MaxRandomValue:      0,
	})
	defaultActionConfigs = append(defaultActionConfigs, actionset.ActionConfig{
		Name:                "exit",
		Random:              false,
		ParamCount:          0,
		DefaultParamValue:   0,
		MinRandomParamValue: 0,
		MaxRandomValue:      0,
	})
	defaultActionConfigs = append(defaultActionConfigs, actionset.ActionConfig{
		Name:                "delayWithSecond",
		Random:              false,
		ParamCount:          1,
		DefaultParamValue:   4,
		MinRandomParamValue: 1,
		MaxRandomValue:      15,
	})
	defaultActionConfigs = append(defaultActionConfigs, actionset.ActionConfig{
		Name:                "delayToNextSlot",
		Random:              false,
		ParamCount:          0,
		DefaultParamValue:   0,
		MinRandomParamValue: 0,
		MaxRandomValue:      0,
	})
	defaultActionConfigs = append(defaultActionConfigs, actionset.ActionConfig{
		Name:                "delayToAfterNextSlot",
		Random:              false,
		ParamCount:          1,
		DefaultParamValue:   4,
		MinRandomParamValue: 1,
		MaxRandomValue:      15,
	})
	defaultActionConfigs = append(defaultActionConfigs, actionset.ActionConfig{
		Name:                "delayToNextNEpochStart",
		Random:              false,
		ParamCount:          1,
		DefaultParamValue:   0,
		MinRandomParamValue: 0,
		MaxRandomValue:      5,
	})
	defaultActionConfigs = append(defaultActionConfigs, actionset.ActionConfig{
		Name:                "delayToNextNEpochEnd",
		Random:              false,
		ParamCount:          1,
		DefaultParamValue:   0,
		MinRandomParamValue: 0,
		MaxRandomValue:      5,
	})
	defaultActionConfigs = append(defaultActionConfigs, actionset.ActionConfig{
		Name:                "delayToNextNEpochHalf",
		Random:              false,
		ParamCount:          1,
		DefaultParamValue:   0,
		MinRandomParamValue: 0,
		MaxRandomValue:      5,
	})
	defaultActionConfigs = append(defaultActionConfigs, actionset.ActionConfig{
		Name:                "delayToEpochEnd",
		Random:              false,
		ParamCount:          0,
		DefaultParamValue:   0,
		MinRandomParamValue: 0,
		MaxRandomValue:      0,
	})
	defaultActionConfigs = append(defaultActionConfigs, actionset.ActionConfig{
		Name:                "delayHalfEpoch",
		Random:              false,
		ParamCount:          1,
		DefaultParamValue:   0,
		MinRandomParamValue: 0,
		MaxRandomValue:      0,
	})

	defaultActionConfigs = append(defaultActionConfigs, actionset.ActionConfig{
		Name:                "storeSignedAttest",
		Random:              false,
		ParamCount:          0,
		DefaultParamValue:   0,
		MinRandomParamValue: 0,
		MaxRandomValue:      0,
	})
	defaultActionConfigs = append(defaultActionConfigs, actionset.ActionConfig{
		Name:                "rePackAttestation",
		Random:              false,
		ParamCount:          0,
		DefaultParamValue:   0,
		MinRandomParamValue: 0,
		MaxRandomValue:      0,
	})
}
