package actionset

import "math/rand"

var (
	_ Action = NullAction{}
	_ Action = ReturnAction{}
	_ Action = ContinueAction{}
	_ Action = AbortAction{}
	_ Action = SkipAction{}
	_ Action = ExitAction{}
)

type baseCmdAction struct {
	config ActionConfig
}

func (b baseCmdAction) DefaultParam() []interface{} {
	if b.config.ParamCount > 0 {
		return []interface{}{b.config.DefaultParamValue}
	}
	return []interface{}{}
}

func (b baseCmdAction) RandomParam() []interface{} {
	if b.config.ParamCount <= 0 {
		return []interface{}{}
	}

	params := make([]interface{}, b.config.ParamCount)
	for i := 0; i < b.config.ParamCount; i++ {
		params[i] = rand.Intn(b.config.MaxRandomValue-b.config.MinRandomParamValue) + b.config.MinRandomParamValue
	}
	return params
}

func (b baseCmdAction) Desc() string { return "Should be replaced" }

func (b baseCmdAction) Name() string { return b.config.Name }

func (b baseCmdAction) MaxParam() int { return b.config.MaxRandomValue }

func (b baseCmdAction) MinParam() int { return b.config.MinRandomParamValue }

func (b baseCmdAction) GetConfig() ActionConfig { return b.config }

func (b baseCmdAction) WithConfig(config ActionConfig) Action {
	b.config = config
	return b
}

func (b baseCmdAction) ActionType() AType { return AnyAction }

type NullAction struct {
	baseCmdAction
}

func (n NullAction) Desc() string { return "# Null action does nothing" }

func (n NullAction) Name() string { return "null" }

type ReturnAction struct {
	baseCmdAction
}

func (r ReturnAction) Desc() string { return "# Return action returns from the current function" }

func (r ReturnAction) Name() string { return "return" }

type ContinueAction struct {
	baseCmdAction
}

func (c ContinueAction) Desc() string { return "# Continue action continues to the next iteration" }

func (c ContinueAction) Name() string { return "continue" }

type AbortAction struct {
	baseCmdAction
}

func (a AbortAction) Desc() string { return "# Abort action aborts the current function" }

func (a AbortAction) Name() string { return "abort" }

type SkipAction struct {
	baseCmdAction
}

func (s SkipAction) Desc() string { return "# Skip action skips the current iteration" }

func (s SkipAction) Name() string { return "skip" }

type ExitAction struct {
	baseCmdAction
}

func (e ExitAction) Desc() string { return "# Exit action exits the current function" }

func (e ExitAction) Name() string { return "exit" }
