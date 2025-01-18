package runtime

const (
	attackerFlag          = "attacker"
	maxValidatorIndexFlag = "max-validator-index"
	minValidatorIndexFlag = "min-validator-index"
	strategyFlag          = "strategy"
	listFlag              = "list"
	logFlag               = "log"
	durationFlag          = "strategy-duration"
)

type updateParam struct {
	attacker          string
	maxValidatorIndex int
	minValidatorIndex int
	strategy          string
	listLibrary       bool
	logPath           string
	duration          int
}

var (
	params = &updateParam{
		attacker:          "",
		strategy:          "",
		logPath:           "",
		maxValidatorIndex: -1,
		minValidatorIndex: 0,
		duration:          60,
		listLibrary:       false,
	}
)
