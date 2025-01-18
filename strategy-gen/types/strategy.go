package types

import (
	"encoding/json"
	"os"
	"strconv"
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
	Uid        string              `json:"uid"`
	Slots      []SlotStrategy      `json:"slots"`
	Validators []ValidatorStrategy `json:"validator"`
}

func (s Strategy) ToFile(name string) error {
	d, _ := json.MarshalIndent(s, "", "  ")
	return os.WriteFile(name, d, 0644)
}

func (s Strategy) String() string {
	d, _ := json.MarshalIndent(s, "", "  ")
	return string(d)
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
	MinValidatorIndex int
}

func (p LibraryParams) GetLatestHackerSlot(duties []ProposerDuty) int {
	latest, _ := strconv.Atoi(duties[0].Slot)
	for _, duty := range duties {
		idx, _ := strconv.Atoi(duty.ValidatorIndex)
		slot, _ := strconv.Atoi(duty.Slot)
		if !p.IsHackValidator(idx) {
			continue
		}
		if slot > latest {
			latest = slot
		}
	}
	return latest
}

func (p LibraryParams) IsHackValidator(valIdx int) bool {
	return valIdx >= p.MinValidatorIndex && valIdx <= p.MaxValidatorIndex
}

func (p LibraryParams) FillterHackDuties(duties []ProposerDuty) []ProposerDuty {
	res := make([]ProposerDuty, 0)
	for _, duty := range duties {
		idx, _ := strconv.Atoi(duty.ValidatorIndex)
		if p.IsHackValidator(idx) {
			res = append(res, duty)
		}
	}
	return res
}

func IsHackValidator(valIdx int, params LibraryParams) bool {
	return valIdx >= params.MinValidatorIndex && valIdx <= params.MaxValidatorIndex
}

type ProposerDuty struct {
	Pubkey         string `json:"pubkey"`
	ValidatorIndex string `json:"validator_index"`
	Slot           string `json:"slot"`
}
