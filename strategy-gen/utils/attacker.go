package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tsinghua-cel/strategy-gen/globalinfo"
	"github.com/tsinghua-cel/strategy-gen/types"
	"net"
	"net/http"
	"time"
)

var (
	hclient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
			DisableCompression:  true,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		},
		Timeout: 30 * time.Second,
	}
)

func UpdateStrategy(url string, strategy types.Strategy) error {
	d, err := json.Marshal(strategy)
	if err != nil {
		return err
	}

	res, err := hclient.Post(fmt.Sprintf("http://%s/v1/update-strategy", url), "application/json", bytes.NewReader(d))
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		// get response string message from body
		var msg string
		err = json.NewDecoder(res.Body).Decode(&msg)

		return fmt.Errorf("failed to update strategy: %s", msg)
	}
	return nil
}

func GetStrategyFeedback(url string, uid string) (types.FeedBackInfo, error) {
	res, err := hclient.Get(fmt.Sprintf("http://%s/v1/strategy-feedback/%s", url, uid))
	if err != nil {
		return types.FeedBackInfo{}, err
	}
	if res.StatusCode != 200 {
		return types.FeedBackInfo{}, fmt.Errorf("failed to get strategy feedback: %s", res.Status)
	}
	// parse response body to types.FeedBackInfo
	var feedback types.FeedBackInfo
	err = json.NewDecoder(res.Body).Decode(&feedback)
	if err != nil {
		return types.FeedBackInfo{}, err
	}
	// do something with feedback
	return feedback, nil
}

func GetChainBaseInfo(url string) (globalinfo.ChainBaseConfig, error) {
	res, err := hclient.Get(fmt.Sprintf("http://%s/v1/chain-base-info", url))
	if err != nil {
		return globalinfo.ChainBaseInfo(), err
	}
	if res.StatusCode != 200 {
		return globalinfo.ChainBaseInfo(), fmt.Errorf("failed to get second per slot: %s", res.Status)
	}

	var chainBaseInfo globalinfo.ChainBaseConfig
	err = json.NewDecoder(res.Body).Decode(&chainBaseInfo)
	if err != nil {
		return globalinfo.ChainBaseInfo(), err
	}
	return chainBaseInfo, nil
}

func GetCurSlot(url string) (int, error) {
	res, err := hclient.Get(fmt.Sprintf("http://%s/v1/curslot", url))
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
	res, err := hclient.Get(fmt.Sprintf("http://%s/v1/slot", url))
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
	res, err := hclient.Get(fmt.Sprintf("http://%s/v1/epoch", url))
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

func GetEpochDuties(url string, epoch int64) ([]types.ProposerDuty, error) {
	res, err := hclient.Get(fmt.Sprintf("http://%s/v1/duties/%d", url, epoch))
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get duties: %s", res.Status)
	}

	var duties []types.ProposerDuty
	err = json.NewDecoder(res.Body).Decode(&duties)
	if err != nil {
		return nil, err
	}
	return duties, nil
}
