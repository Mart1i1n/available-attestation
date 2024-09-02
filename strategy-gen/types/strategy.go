package types

import (
	"encoding/json"
	"os"
)

type ValidatorStrategy struct {
	ValidatorIndex    int `json:"validator_index"`
	AttackerStartSlot int `json:"attacker_start_slot"`
	AttackerEndSlot   int `json:"attacker_end_slot"`
}

type SlotStrategy struct {
	Slot    string            `json:"slot"`
	Level   int               `json:"level"`
	Actions map[string]string `json:"actions"`
}

type Strategy struct {
	Slots      []SlotStrategy      `json:"slots"`
	Validators []ValidatorStrategy `json:"validator"`
}

func (s Strategy) ToFile(name string) error {
	d, _ := json.MarshalIndent(s, "", "  ")
	return os.WriteFile(name, d, 0644)
}

func GetValidatorStrategy(startIndex, endIndex int, startSlot, endSlot int) []ValidatorStrategy {
	res := make([]ValidatorStrategy, 0)
	for i := startIndex; i <= endIndex; i++ {
		res = append(res, ValidatorStrategy{
			ValidatorIndex:    i,
			AttackerStartSlot: startSlot,
			AttackerEndSlot:   endSlot,
		})
	}
	return res
}

type LibraryParams struct {
	Attacker          string
	MaxValidatorIndex int
}
