package actionset

import (
	"encoding/json"
	"strings"
)

type AType int

const (
	AnyAction AType = iota
	BlockAction
	AttestAction
)

type ActionConfig struct {
	Name                string `json:"name" yaml:"name"`
	Random              bool   `json:"random" yaml:"random"`
	ParamCount          int    `json:"param_count" yaml:"param_count"`
	DefaultParamValue   int    `json:"default_param_value" yaml:"default_param_value"`
	MinRandomParamValue int    `json:"min_random_param_value" yaml:"min_random_param_value"`
	MaxRandomValue      int    `json:"max_random_param_value" yaml:"max_random_param_value"`
}

// implement ActionConfig json.Marshaler and yaml.Marshaler
func (ac ActionConfig) MarshalJSON() ([]byte, error) {
	type simpleConf struct {
		Name string `json:"name"`
	}
	type fullConf struct {
		Name                string `json:"name" yaml:"name"`
		Random              bool   `json:"random" yaml:"random"`
		ParamCount          int    `json:"param_count" yaml:"param_count"`
		DefaultParamValue   int    `json:"default_param_value" yaml:"default_param_value"`
		MinRandomParamValue int    `json:"min_random_param_value" yaml:"min_random_param_value"`
		MaxRandomValue      int    `json:"max_random_param_value" yaml:"max_random_param_value"`
	}
	if ac.ParamCount == 0 {
		return json.MarshalIndent(simpleConf{Name: ac.Name}, "", "  ")
	} else {
		return json.MarshalIndent(fullConf{
			Name:                ac.Name,
			Random:              ac.Random,
			ParamCount:          ac.ParamCount,
			DefaultParamValue:   ac.DefaultParamValue,
			MinRandomParamValue: ac.MinRandomParamValue,
			MaxRandomValue:      ac.MaxRandomValue,
		}, "", "  ")
	}
}

func (ac ActionConfig) MarshalYAML() (interface{}, error) {
	type simpleConf struct {
		Name string `json:"name"`
	}
	type fullConf struct {
		Name                string `json:"name" yaml:"name"`
		Random              bool   `json:"random" yaml:"random"`
		ParamCount          int    `json:"param_count" yaml:"param_count"`
		DefaultParamValue   int    `json:"default_param_value" yaml:"default_param_value"`
		MinRandomParamValue int    `json:"min_random_param_value" yaml:"min_random_param_value"`
		MaxRandomValue      int    `json:"max_random_param_value" yaml:"max_random_param_value"`
	}
	if ac.ParamCount == 0 {
		return simpleConf{Name: ac.Name}, nil
	} else {
		return fullConf{
			Name:                ac.Name,
			Random:              ac.Random,
			ParamCount:          ac.ParamCount,
			DefaultParamValue:   ac.DefaultParamValue,
			MinRandomParamValue: ac.MinRandomParamValue,
			MaxRandomValue:      ac.MaxRandomValue,
		}, nil
	}
}

type Action interface {
	Name() string
	MaxParam() int
	MinParam() int
	ActionType() AType
	DefaultParam() []interface{}
	RandomParam() []interface{}
	Desc() string
	GetConfig() ActionConfig
	WithConfig(ActionConfig) Action
}

var anyAction = []Action{
	NullAction{},
	ReturnAction{},
	ContinueAction{},
	AbortAction{},
	SkipAction{},
	ExitAction{},
	delayWithSecondAction{},
	delayToNextSlotAction{},
	delayToAfterNextSlotAction{},
	delayToNextNEpochStartAction{},
	delayToNextNEpochEndAction{},
	delayToNextNEpochHalfAction{},
	delayToEpochEndAction{},
	delayHalfEpochAction{},
}

var attestAction = []Action{
	storeSignedAttestAction{},
	rePackAttestationAction{},
}

var blockAction = []Action{}

func GetBlockActionSet() []Action {
	a := make([]Action, 0)
	a = append(a, anyAction...)
	a = append(a, blockAction...)

	return a
}

func GetBlockActionNameList() []string {
	a := make([]string, 0)
	for _, action := range anyAction {
		a = append(a, action.Name())
	}

	for _, action := range blockAction {
		a = append(a, action.Name())
	}
	return a
}

func GetAttestActionSet() []Action {
	a := make([]Action, 0)
	a = append(a, anyAction...)
	a = append(a, attestAction...)
	return a
}

func GetAttestActionNameList() []string {
	a := make([]string, 0)
	for _, action := range anyAction {
		a = append(a, action.Name())
	}

	for _, action := range attestAction {
		a = append(a, action.Name())
	}
	return a
}

func GetAllActionSet() []Action {
	all := make([]Action, 0)
	all = append(all, anyAction...)
	all = append(all, blockAction...)
	all = append(all, attestAction...)
	return all
}

func GetActionByName(name string) Action {
	for _, a := range GetAllActionSet() {
		if strings.ToLower(a.Name()) == strings.ToLower(name) {
			return a
		}
	}
	return nil
}

func GetActionByConfig(actionConfig ActionConfig) Action {
	for _, a := range GetAllActionSet() {
		if strings.ToLower(a.Name()) == strings.ToLower(actionConfig.Name) {
			return a
		}
	}
	return nil
}
