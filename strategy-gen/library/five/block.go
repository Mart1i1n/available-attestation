package five

import (
	"fmt"
	"github.com/tsinghua-cel/strategy-gen/types"
	"github.com/tsinghua-cel/strategy-gen/utils"
)

/*
提前一个epoch 查看下一个epoch的恶意节点出块顺序；
两个恶意节点之间穿插了一个诚实节点的块，让第二恶意节点的区块的parent，指向上一个恶意节点的slot.
delay 策略：第一个恶意节点的区块，广播delay 1个slot;
*/

func GenSlotStrategy(duties []interface{}) []types.SlotStrategy {
	strategys := make([]types.SlotStrategy, 0)
	for i := 0; i < len(duties); i++ {
		duty := duties[i].([]utils.ProposerDuty)
		if len(duty) != 3 {
			continue
		}
		a := duty[0]
		//b := duty[1]
		c := duty[2]

		slotStrategy := types.SlotStrategy{
			Slot:    fmt.Sprintf("%d", c.Slot),
			Level:   1,
			Actions: make(map[string]string),
		}
		slotStrategy.Actions["BlockGetNewParentRoot"] = fmt.Sprintf("modifyParentRoot:%d", a.Slot)
		strategys = append(strategys, slotStrategy)
	}

	return strategys

}
