package config

import (
	"encoding/json"
	"fmt"
	"github.com/tsinghua-cel/strategy-gen/actionset"
	"github.com/tsinghua-cel/strategy-gen/pointset"
	"github.com/tsinghua-cel/strategy-gen/types"
	"gopkg.in/yaml.v3"
	"log"
	"math/rand"
	"os"
	"strings"
)

// Config defines the server configuration params
type Config struct {
	ValidatorCount     int                               `json:"validator_count" yaml:"validator_count"`
	StartSlot          int                               `json:"start_slot" yaml:"start_slot"`
	EndSlot            int                               `json:"end_slot" yaml:"end_slot"`
	EnableBlockPoints  string                            `json:"enable_block_points" yaml:"enable_block_points"`
	EnableBlockActions string                            `json:"enable_block_actions" yaml:"enable_block_actions"`
	EnableAttPoints    string                            `json:"enable_att_points" yaml:"enable_att_points"`
	EnableAttActions   string                            `json:"enable_att_actions" yaml:"enable_att_actions"`
	ActionsConfig      map[string]actionset.ActionConfig `json:"actions_config" yaml:"actions_config"`
}

const (
	DefaultValidatorCount     = 21
	DefaultStartSlot          = 0
	DefaultEndSlot            = 10000
	DefaultEnableBlockPoints  = "BlockDelayForReceiveBlock,BlockBeforeBroadCast"
	DefaultEnableAttPoints    = "AttestBeforeBroadCast"
	DefaultEnableBlockActions = "null,delayWithSecond,delayToAfterNextSlot,delayToNextNEpochStart,delayToNextNEpochHalf,delayToEpochEnd,return"
	DefaultEnableAttActions   = "null,delayWithSecond,delayToAfterNextSlot,return"
)

// DefaultConfig returns the default server configuration
func DefaultConfig() *Config {
	defaultConfigInit()

	conf := &Config{
		ValidatorCount:     DefaultValidatorCount,
		StartSlot:          DefaultStartSlot,
		EndSlot:            DefaultEndSlot,
		EnableBlockPoints:  DefaultEnableBlockPoints,
		EnableBlockActions: DefaultEnableBlockActions,
		EnableAttPoints:    DefaultEnableAttPoints,
		EnableAttActions:   DefaultEnableAttActions,
		ActionsConfig:      make(map[string]actionset.ActionConfig),
	}
	for _, action := range defaultActionConfigs {
		conf.ActionsConfig[action.Name] = action
	}

	return conf
}

// ReadConfigFile reads the config file from the specified path, builds a Config object
// and returns it.
//
// Supported file types: .json, .hcl, .yaml, .yml
func ReadConfigFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var unmarshalFunc func([]byte, interface{}) error

	switch {
	case strings.HasSuffix(path, ".json"):
		unmarshalFunc = json.Unmarshal
	case strings.HasSuffix(path, ".yaml"), strings.HasSuffix(path, ".yml"):
		unmarshalFunc = yaml.Unmarshal
	default:
		return nil, fmt.Errorf("suffix of %s is neither hcl, json, yaml nor yml", path)
	}

	config := DefaultConfig()
	if err := unmarshalFunc(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

func ConfigToStrategy(mode int, conf Config) types.Strategy {
	strategy := types.Strategy{}
	valEndIndex := conf.ValidatorCount - 1
	validators := types.GetValidatorStrategy(0, valEndIndex, conf.StartSlot, conf.EndSlot)
	attActions := make([]actionset.Action, 0)

	enableAttPoints := make([]string, 0)
	{
		points := strings.Split(conf.EnableAttPoints, ",")
		for _, point := range points {
			p := pointset.GetPointByName(point)
			if p != "" {
				enableAttPoints = append(enableAttPoints, p)
			}
		}
	}
	enableBlockPoints := make([]string, 0)
	{
		points := strings.Split(conf.EnableBlockPoints, ",")
		for _, point := range points {
			p := pointset.GetPointByName(point)
			if p != "" {
				enableBlockPoints = append(enableBlockPoints, p)
			}
		}
	}

	enables := strings.Split(conf.EnableAttActions, ",")
	for _, enable := range enables {
		actionConfig, exist := conf.ActionsConfig[enable]
		if !exist {
			log.Printf("action %s not exist in config\n", enable)
			continue
		}
		action := actionset.GetActionByName(enable).WithConfig(actionConfig)
		if action.ActionType() != actionset.AttestAction && action.ActionType() != actionset.AnyAction {
			log.Printf("action %s is not block action\n", enable)
			continue
		}
		attActions = append(attActions, action)
	}

	blockActions := make([]actionset.Action, 0)
	enables = strings.Split(conf.EnableBlockActions, ",")
	for _, enable := range enables {
		actionConfig, exist := conf.ActionsConfig[enable]
		if !exist {
			log.Printf("block action %s not exist in config\n", enable)
			continue
		}
		action := actionset.GetActionByName(enable).WithConfig(actionConfig)
		if action.ActionType() != actionset.BlockAction && action.ActionType() != actionset.AnyAction {
			log.Printf("action %s is not block action\n", enable)
			continue
		}
		blockActions = append(blockActions, action)
	}

	slotStrategy := make([]types.SlotStrategy, 0)
	for slot := conf.StartSlot; slot <= conf.EndSlot; slot++ {
		strate := types.SlotStrategy{Actions: make(map[string]string)}
		strate.Slot = fmt.Sprintf("%d", slot)
		strate.Level = 0
		for _, point := range enableBlockPoints {
			action := randomAction(blockActions)
			astr := action.Name()
			if action.GetConfig().ParamCount > 0 {
				param := []interface{}{}
				if mode == 1 {
					param = action.RandomParam()
				} else {
					param = action.DefaultParam()
				}
				for _, p := range param {
					astr = fmt.Sprintf("%s:%v", astr, p)
				}
			}
			strate.Actions[point] = astr
		}
		for _, point := range enableAttPoints {
			action := randomAction(attActions)
			astr := action.Name()
			if action.GetConfig().ParamCount > 0 {
				param := []interface{}{}
				if mode == 1 {
					param = action.RandomParam()
				} else {
					param = action.DefaultParam()
				}
				for _, p := range param {
					astr = fmt.Sprintf("%s:%v", astr, p)
				}
			}
			strate.Actions[point] = astr
		}
		slotStrategy = append(slotStrategy, strate)
	}
	strategy.Slots = slotStrategy
	strategy.Validators = validators

	return strategy
}

func randomAction(actions []actionset.Action) actionset.Action {
	if len(actions) == 0 {
		return nil
	}
	return actions[rand.Intn(len(actions))]
}
