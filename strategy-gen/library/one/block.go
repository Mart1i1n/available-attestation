package one

import (
	"fmt"
	"github.com/tsinghua-cel/strategy-gen/pointset"
	"github.com/tsinghua-cel/strategy-gen/types"
	"github.com/tsinghua-cel/strategy-gen/utils"
	"strconv"
)

/*
提前一个epoch 查看下一个epoch的恶意节点出块顺序；
如果有连续两个以上恶意节点出块，开始进行策略；
delay策略：blockdelay 到最后一个恶意节点出块的下一个slot；
恶意节点的投票者 开始做恶，对投票进行delay，执行的策略和blockdelay一样。
*/
func BlockStrategy(cur, end int, actions map[string]string) {
	point := pointset.GetPointByName("BlockBeforeBroadCast")
	actions[point] = fmt.Sprintf("%s:%d", "delayWithSecond", (end+1-cur)*12)
}

func AttestStrategy(cur, end int, actions map[string]string) {
	point := pointset.GetPointByName("AttestBeforeBroadCast")
	actions[point] = fmt.Sprintf("%s:%d", "delayWithSecond", (end+1-cur)*12)
}

func GenSlotStrategy(allHacks []interface{}) []types.SlotStrategy {
	strategys := make([]types.SlotStrategy, 0)
	for _, subduties := range allHacks {
		duties := subduties.([]utils.ProposerDuty)
		//begin, _ := strconv.Atoi(duties[0].Slot)
		end, _ := strconv.Atoi(duties[len(duties)-1].Slot)

		for i := 0; i < len(duties); i++ {
			slot, _ := strconv.Atoi(duties[i].Slot)
			//idx, _ := strconv.Atoi(duties[i].ValidatorIndex)
			strategy := types.SlotStrategy{
				Slot:    duties[i].Slot,
				Level:   0,
				Actions: make(map[string]string),
			}
			BlockStrategy(slot, end, strategy.Actions)
			AttestStrategy(slot, end, strategy.Actions)
			strategys = append(strategys, strategy)
		}
	}

	return strategys

}
