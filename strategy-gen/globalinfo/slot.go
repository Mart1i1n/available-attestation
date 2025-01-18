package globalinfo

type ChainBaseConfig struct {
	SecondsPerSlot int   `json:"secondsPerSlot"`
	SlotsPerEpoch  int   `json:"slotsPerEpoch"`
	GenesisTime    int64 `json:"genesisTime"`
}

var (
	gInfo = ChainBaseConfig{}
)

func Init(g ChainBaseConfig) {
	gInfo = g
}

func ChainBaseInfo() ChainBaseConfig {
	return gInfo
}
