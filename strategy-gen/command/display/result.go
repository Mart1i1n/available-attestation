package display

import "encoding/json"

type ActionPointResult struct {
	BlockPoints []string `json:"block_points"`
	AttPoints   []string `json:"attest_points"`
}

type ActionResult struct {
	BlockActionList []string `json:"block_action_list"`
	AttActionList   []string `json:"attest_action_list"`
}

type DisplayResult struct {
	ActionPointResult ActionPointResult `json:"action_point_list"`
	ActionResult      ActionResult      `json:"action_list"`
}

func (d DisplayResult) GetOutput() string {
	info, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(info)
}
