package ethchain

import "github.com/ethereum/go-ethereum/core/types"

type EthTx struct {
	*types.Transaction
}

func (e EthTx) IsChainTx() bool {
	return true
}
