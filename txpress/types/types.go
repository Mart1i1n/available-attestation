package types

import (
	"time"
)

type ChainTx interface {
	IsChainTx() bool
}

type BlockInfo struct {
	Timestamp   int64
	Number      int64
	TxCount     int64
	Beneficiary string
}

type ChainPlugin interface {
	CreateTxs(count int, checkNonce bool) ([]ChainTx, error)
	SendTxs(txs []ChainTx) ([]string, error)
	TxReceipt(hash string) error
	TxBlock(hash string) (int, error)
	GetBlockInfo(number int64) (BlockInfo, error)
	LatestBlockInfo() (BlockInfo, error)
	SecondPerBlock() int
	Id() string
}

type RunConfig struct {
	BaseCount     int
	Interval      time.Duration
	Batch         int
	Round         int
	IncRate       int
	BeginToStart  int
	ForceIncrease bool
}

type ChainConfig struct {
	Rpcs          []string `json:"rpc-nodes"`
	Name          string   `json:"chain-name"`
	BaseCount     int      `json:"base-count"`
	Round         int      `json:"round"`
	Interval      int      `json:"interval"`
	Batch         int      `json:"batch"`
	Receiver      string   `json:"receiver"`
	Amount        string   `json:"amount"`
	Accounts      string   `json:"accounts"`
	IncRate       int      `json:"inc-rate"`
	ForceIncrease bool     `json:"force-increase"`
	BeginToStart  int      `json:"begin-to-start"`
}

type Record struct {
	Begin     int
	End       int
	TotalTime int
	TotalTx   int
	Tps       int
}
