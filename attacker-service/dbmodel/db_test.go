package dbmodel

import (
	"github.com/tsinghua-cel/attacker-service/config"
	"testing"
)

func init() {
	DbInit(config.MysqlConfig{
		Host:   "127.0.0.1",
		Port:   3306,
		User:   "root",
		Passwd: "12345678",
		DbName: "eth",
	})
}

func TestGetRewardListByValidatorIndex(t *testing.T) {
	list := GetRewardListByValidatorIndex(0)
	t.Log(list)
}
