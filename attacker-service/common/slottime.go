package common

import "github.com/tsinghua-cel/attacker-service/types"

type slotTimeTool struct {
	SecondsPerSlot int
	SlotsPerEpoch  int64
	GenesisTime    int64
}

var (
	// SlotTimeTool is a global slot time tool
	tool     *slotTimeTool
	baseInfo *types.ChainBaseInfo
)

func InitSlotTool(secondsPerSlot int, slotsPerEpoch int64, genesisTime int64) {
	if tool == nil {
		tool = &slotTimeTool{
			SecondsPerSlot: secondsPerSlot,
			SlotsPerEpoch:  slotsPerEpoch,
			GenesisTime:    genesisTime,
		}
	}
	if baseInfo == nil {
		baseInfo = &types.ChainBaseInfo{
			SecondsPerSlot: secondsPerSlot,
			SlotsPerEpoch:  int(slotsPerEpoch),
			GenesisTime:    genesisTime,
		}
	}
}

func SlotToEpoch(slot int64) int64 {
	return slot / int64(tool.SlotsPerEpoch)
}

func EpochEnd(epoch int64) int64 {
	return (epoch+1)*int64(tool.SlotsPerEpoch) - 1
}

func EpochStart(epoch int64) int64 {
	return epoch * int64(tool.SlotsPerEpoch)
}

func TimeToSlot(slot int64) int64 {
	return tool.GenesisTime + int64(slot*int64(tool.SecondsPerSlot))
}

func GetChainBaseInfo() types.ChainBaseInfo {
	return *baseInfo
}
