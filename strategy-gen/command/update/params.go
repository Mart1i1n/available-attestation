package update

const (
	intervalFlag  = "interval"
	attackerFlag  = "attacker"
	sliceFlag     = "slice"
	modeFlag      = "mode"       // 0: sorted loop 1: random loop
	loopCountFlag = "loop-count" // 0 is never stop.
)

type updateParam struct {
	interval   int
	fileSlice  []string
	updateMode int
	loopCount  int
	attacker   string
}

var (
	params = &updateParam{
		interval:   1200,
		updateMode: 0,
		loopCount:  0,
		fileSlice:  []string{},
		attacker:   "",
	}
)
