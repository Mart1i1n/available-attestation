package actionset

var (
	_ Action = delayWithSecondAction{}
	_ Action = delayToNextSlotAction{}
	_ Action = delayToAfterNextSlotAction{}
	_ Action = delayToNextNEpochStartAction{}
	_ Action = delayToNextNEpochEndAction{}
	_ Action = delayToNextNEpochHalfAction{}
	_ Action = delayToEpochEndAction{}
	_ Action = delayHalfEpochAction{}
)

type delayWithSecondAction struct {
	baseCmdAction
}

func (d delayWithSecondAction) Desc() string { return "# Delay with seconds" }

func (d delayWithSecondAction) Name() string { return "delayWithSecond" }

func (d delayWithSecondAction) ActionType() AType { return AnyAction }

type delayToNextSlotAction struct {
	baseCmdAction
}

func (d delayToNextSlotAction) Desc() string { return "# Delay to next slot" }

func (d delayToNextSlotAction) Name() string { return "delayToNextSlot" }

func (d delayToNextSlotAction) ActionType() AType { return AnyAction }

type delayToAfterNextSlotAction struct {
	baseCmdAction
}

func (d delayToAfterNextSlotAction) Desc() string { return "# Delay to after next slot" }

func (d delayToAfterNextSlotAction) Name() string { return "delayToAfterNextSlot" }

func (d delayToAfterNextSlotAction) ActionType() AType { return AnyAction }

type delayToNextNEpochStartAction struct {
	baseCmdAction
}

func (d delayToNextNEpochStartAction) Desc() string { return "# Delay to next n epoch start" }

func (d delayToNextNEpochStartAction) Name() string { return "delayToNextNEpochStart" }

func (d delayToNextNEpochStartAction) ActionType() AType { return AnyAction }

type delayToNextNEpochEndAction struct {
	baseCmdAction
}

func (d delayToNextNEpochEndAction) Desc() string { return "# Delay to next n epoch end" }

func (d delayToNextNEpochEndAction) Name() string { return "delayToNextNEpochEnd" }

func (d delayToNextNEpochEndAction) ActionType() AType { return AnyAction }

type delayToNextNEpochHalfAction struct {
	baseCmdAction
}

func (d delayToNextNEpochHalfAction) Desc() string { return "# Delay to next n epoch half" }

func (d delayToNextNEpochHalfAction) Name() string { return "delayToNextNEpochHalf" }

func (d delayToNextNEpochHalfAction) ActionType() AType { return AnyAction }

type delayToEpochEndAction struct {
	baseCmdAction
}

func (d delayToEpochEndAction) Desc() string { return "# Delay to epoch end" }

func (d delayToEpochEndAction) Name() string { return "delayToEpochEnd" }

func (d delayToEpochEndAction) ActionType() AType { return AnyAction }

type delayHalfEpochAction struct {
	baseCmdAction
}

func (d delayHalfEpochAction) Desc() string { return "# Delay half epoch" }

func (d delayHalfEpochAction) Name() string { return "delayHalfEpoch" }

func (d delayHalfEpochAction) ActionType() AType { return AnyAction }
