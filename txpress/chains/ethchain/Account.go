package ethchain

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"math/big"
	"os"
	"strings"
)

type Account struct {
	Address common.Address    `json:"address"`
	Private string            `json:"private"`
	Nonce   uint64            `json:"nonce"`
	PK      *ecdsa.PrivateKey `json:"-"`
}

type txConfig struct {
	receiver common.Address
	amount   *big.Int
}

func (acc *Account) MakeNormalTx(cfg txConfig, nonce uint64) *types.Transaction {
	txData := &types.LegacyTx{
		Nonce:    nonce,
		To:       &cfg.receiver,
		Value:    cfg.amount,
		Gas:      21000,
		GasPrice: big.NewInt(1000000000),
		Data:     nil,
	}
	return types.NewTx(txData)
}

func (acc *Account) SignTx(chainId *big.Int, tx *types.Transaction) (*types.Transaction, error) {
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), acc.PK)
	if err != nil {
		log.Error("sign tx failed", "err", err)
		return nil, err
	}
	return signedTx, nil
}

func createAccounts(count int) []*Account {
	accs := make([]*Account, 0, count)
	for i := 0; i < count; i++ {
		pk, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(pk.PublicKey)
		private := hexutil.Encode(crypto.FromECDSA(pk))
		accs = append(accs, &Account{Address: addr, Private: private, PK: pk})
	}
	return accs
}

func padding(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat("0", length-len(s)) + s
}

func GetAccountJson(accountFile string) []*Account {
	data, err := os.ReadFile(accountFile)
	if err != nil || len(data) == 0 {
		return []*Account{}
	} else {
		accs := make([]*Account, 0)
		err = json.Unmarshal(data, &accs)
		if err != nil {
			log.Error("unmarshal account failed", "err", err)
		}

		for _, acc := range accs {
			private := acc.Private
			if strings.HasPrefix(private, "0x") {
				private = acc.Private[2:]
			}
			private = padding(private, 64)
			acc.PK, err = crypto.HexToECDSA(private)
			if err != nil {
				log.Error("hex to ecdsa failed", "err", err)
			}
		}
		log.Info("get accounts from json", "len", len(accs))
		return accs
	}
}

func CheckAccountNonce(client *ethclient.Client, account *Account) {
	nonce, err := client.NonceAt(context.Background(), account.Address, nil)
	if err != nil {
		// ignore err
	} else {
		log.Debug("check account nonce", "account", account.Address, "nonce", nonce)
		account.Nonce = nonce
	}

	return
}
