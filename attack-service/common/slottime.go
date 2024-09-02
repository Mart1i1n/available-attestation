package common

type slotTimeTool struct {
	SecondsPerSlot int
	GenesisTime    int64
}

var (
	// SlotTimeTool is a global slot time tool
	tool *slotTimeTool
)

func InitSlotTime(secondsPerSlot int, genesisTime int64) {
	if tool == nil {
		tool = &slotTimeTool{
			SecondsPerSlot: secondsPerSlot,
			GenesisTime:    genesisTime,
		}
	}
}

func TimeToSlot(slot int64) int64 {
	return tool.GenesisTime + int64(slot*int64(tool.SecondsPerSlot))
}
