package ethchain

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"testing"
)

func TestChainValue(t *testing.T) {
	var account = []string{
		"0x122771410575cf99a131eBEF6B7AA3f8CC6c2d76",
		"0x572379d1D5e38500806D7C4232f8Abc9697116f8",
		"0x1A1c2323ac70791eA77735F9FEF3504d27621E0D",
		"0x2B32D30a38C19EBFcFb45723AD70B1401BEC21D6",
	}
	rpc := "http://13.41.176.56:27658"
	client, err := ethclient.Dial(rpc)
	if err != nil {
		t.Error("failed to dial rpc", err)
	}

	for _, addr := range account {
		bal, err := client.BalanceAt(context.Background(), common.HexToAddress(addr), nil)
		if err != nil {
			t.Error("failed to get balance", err)
		}
		floatBal := big.NewFloat(0).SetInt(bal)
		floatBal.Quo(floatBal, big.NewFloat(1e18))
		fmt.Printf("account %s balance: %s\n", addr, floatBal.String())
	}

}
