package chains

import (
	log "github.com/sirupsen/logrus"
	"github.com/xueqianLu/txpress/chains/ethchain"
	"github.com/xueqianLu/txpress/types"
	"time"
)

func NewChains(config types.ChainConfig) []types.ChainPlugin {
	chains := make([]types.ChainPlugin, 0)

	var createFunc func(rpc string, index int, config types.ChainConfig) (types.ChainPlugin, error)
	switch config.Name {
	case "eth":
		createFunc = ethchain.NewEthChain
	default:
		log.Errorf("unsupport chain %s", config.Name)
		return nil
	}

	for i, rpc := range config.Rpcs {
		chain, err := createFunc(rpc, i, config)
		if err != nil {
			log.Errorf("create chain %s for with rpc(%s) failed", config.Name, rpc)
			continue
		}
		chains = append(chains, chain)
	}
	return chains
}

func CalcTps(chain types.ChainPlugin, minBlock, maxBlock int, duration time.Duration) types.Record {
	txCount := int64(0)
	for i := minBlock; i <= maxBlock; i++ {
		block, err := chain.GetBlockInfo(int64(i))
		if err != nil {
			log.Errorf("get block info failed: %s", err)
			continue
		}
		txCount += block.TxCount
	}
	record := types.Record{
		Begin:     int(minBlock),
		End:       int(maxBlock),
		TotalTime: int(duration.Seconds()),
		TotalTx:   int(txCount),
	}
	if record.TotalTime > 0 {
		record.Tps = int(txCount) / (record.TotalTime)
	}

	return record
}
