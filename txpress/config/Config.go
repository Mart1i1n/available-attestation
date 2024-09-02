package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/xueqianLu/txpress/types"
	"io/ioutil"
)

func ParseConfig(path string) (types.ChainConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("get config failed", "err", err)
		return types.ChainConfig{}, err
	}
	conf := types.ChainConfig{}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Error("unmarshal config failed", "err", err)
		return types.ChainConfig{}, err
	}
	return conf, nil
}
