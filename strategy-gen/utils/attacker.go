package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tsinghua-cel/strategy-gen/types"
	"net/http"
)

func UpdateStrategy(url string, strategy types.Strategy) error {
	d, err := json.Marshal(strategy)
	if err != nil {
		return err
	}

	res, err := http.Post(fmt.Sprintf("http://%s/v1/update-strategy", url), "application/json", bytes.NewReader(d))
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("failed to generate strategy: %s", res.Status)
	}
	return nil
}

func GetCurSlot(url string) (int, error) {
	res, err := http.Get(fmt.Sprintf("http://%s/v1/curslot", url))
	if err != nil {
		return 0, err
	}
	if res.StatusCode != 200 {
		return 0, fmt.Errorf("failed to get slot: %s", res.Status)
	}

	var slot int
	err = json.NewDecoder(res.Body).Decode(&slot)
	if err != nil {
		return 0, err
	}
	return slot, nil
}

func GetHeadSlot(url string) (int, error) {
	res, err := http.Get(fmt.Sprintf("http://%s/v1/slot", url))
	if err != nil {
		return 0, err
	}
	if res.StatusCode != 200 {
		return 0, fmt.Errorf("failed to get slot: %s", res.Status)
	}

	var slot int
	err = json.NewDecoder(res.Body).Decode(&slot)
	if err != nil {
		return 0, err
	}
	return slot, nil
}

func GetEpoch(url string) (int, error) {
	res, err := http.Get(fmt.Sprintf("http://%s/v1/epoch", url))
	if err != nil {
		return 0, err
	}
	if res.StatusCode != 200 {
		return 0, fmt.Errorf("failed to get epoch: %s", res.Status)
	}

	var epoch int
	err = json.NewDecoder(res.Body).Decode(&epoch)
	if err != nil {
		return 0, err
	}
	return epoch, nil
}

type ProposerDuty struct {
	Pubkey         string `json:"pubkey"`
	ValidatorIndex string `json:"validator_index"`
	Slot           string `json:"slot"`
}

func GetEpochDuties(url string, epoch int64) ([]ProposerDuty, error) {
	res, err := http.Get(fmt.Sprintf("http://%s/v1/duties/%d", url, epoch))
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get duties: %s", res.Status)
	}

	var duties []ProposerDuty
	err = json.NewDecoder(res.Body).Decode(&duties)
	if err != nil {
		return nil, err
	}
	return duties, nil
}
