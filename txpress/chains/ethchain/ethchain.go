package ethchain

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"github.com/xueqianLu/txpress/types"
	"math/big"
)

var _ types.ChainPlugin = &EthChain{}

type EthChain struct {
	ctx      context.Context
	chainId  *big.Int
	rpc      string
	index    int
	accounts []*Account
	client   *ethclient.Client
	config   types.ChainConfig
}

func (e EthChain) LatestBlockInfo() (types.BlockInfo, error) {
	latest, err := e.client.BlockNumber(e.ctx)
	if err != nil {
		return types.BlockInfo{}, err
	}
	return e.GetBlockInfo(int64(latest))
}

func (e EthChain) CreateTxs(count int, checkNonce bool) ([]types.ChainTx, error) {
	txs := make([]types.ChainTx, 0)
	updated := make(map[common.Address]bool)

	amount, _ := new(big.Int).SetString(e.config.Amount, 10)
	txcfg := txConfig{
		receiver: common.HexToAddress(e.config.Receiver),
		amount:   amount,
	}
	for i := 0; i < count; i++ {
		acc := e.accounts[i%len(e.accounts)]
		if _, exist := updated[acc.Address]; !exist && checkNonce {
			acc.Nonce, _ = e.client.NonceAt(e.ctx, acc.Address, nil)
			log.WithFields(log.Fields{
				"acc":   acc.Address.String(),
				"nonce": acc.Nonce,
			}).Info("update account nonce")
			updated[acc.Address] = true
		}
		tx := acc.MakeNormalTx(txcfg, acc.Nonce)
		stx, _ := acc.SignTx(e.chainId, tx)
		txs = append(txs, EthTx{stx})

		acc.Nonce++
	}
	return txs, nil
}

func (e EthChain) SendTxs(txs []types.ChainTx) ([]string, error) {
	hashes := make([]string, 0)
	for _, tx := range txs {
		etx := tx.(EthTx)
		err := e.client.SendTransaction(e.ctx, etx.Transaction)
		if err != nil {
			log.WithFields(log.Fields{
				"chain": e.config.Name,
				"rpc":   e.rpc,
				"index": e.index,
				"tx":    etx.Transaction.Hash().String(),
				"err":   err,
			}).Error("send tx failed")
			continue
		} else {
			//log.WithFields(log.Fields{
			//	"chain": e.config.Name,
			//	"rpc":   e.rpc,
			//	"index": e.index,
			//	"tx":    etx.Transaction.Hash().String(),
			//}).Info("send tx success")
		}
		hash := etx.Transaction.Hash()
		hashes = append(hashes, hash.String())

	}
	return hashes, nil
}

func (e EthChain) TxReceipt(hash string) error {
	_, err := e.client.TransactionReceipt(e.ctx, common.HexToHash(hash))
	if err != nil {
		log.WithFields(log.Fields{
			"chain": e.config.Name,
			"rpc":   e.rpc,
			"index": e.index,
			"err":   err,
		}).Error("get tx receipt failed")
		return err
	}
	return err
}

func (e EthChain) TxBlock(hash string) (int, error) {
	receipt, err := e.client.TransactionReceipt(e.ctx, common.HexToHash(hash))
	if err != nil {
		log.WithFields(log.Fields{
			"chain": e.config.Name,
			"rpc":   e.rpc,
			"index": e.index,
			"err":   err,
			"hash":  hash,
		}).Error("get tx receipt failed")
		return 0, err
	}
	return int(receipt.BlockNumber.Int64()), nil
}

func (e EthChain) GetBlockInfo(number int64) (types.BlockInfo, error) {
	blk := new(big.Int).SetInt64(number)
	info, err := e.client.BlockByNumber(e.ctx, blk)
	if err != nil {
		log.WithFields(log.Fields{
			"chain": e.config.Name,
			"rpc":   e.rpc,
			"index": e.index,
			"err":   err,
		}).Error("get block info failed")
		return types.BlockInfo{}, err
	}
	return types.BlockInfo{
		Number:      info.Number().Int64(),
		Timestamp:   int64(info.Time()),
		TxCount:     int64(len(info.Transactions())),
		Beneficiary: info.Coinbase().String(),
	}, nil
}

func (e EthChain) Id() string {
	return fmt.Sprintf("%s-%d", e.config.Name, e.index)
}

func (e EthChain) SecondPerBlock() int {
	return 12
}

var (
	totalAccounts []*Account
)

func NewEthChain(rpc string, index int, config types.ChainConfig) (types.ChainPlugin, error) {
	ctx := context.TODO()
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}
	chainId, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	total := len(config.Rpcs)

	if totalAccounts == nil {
		totalAccounts = GetAccountJson(config.Accounts)
	}
	chainAccounts := make([]*Account, 0)
	for i := 0; i < len(totalAccounts); i++ {
		if i%total == index {
			chainAccounts = append(chainAccounts, totalAccounts[i])
		}
	}

	return &EthChain{
		ctx:      ctx,
		rpc:      rpc,
		index:    index,
		accounts: chainAccounts,
		client:   client,
		chainId:  chainId,
		config:   config,
	}, nil

}
