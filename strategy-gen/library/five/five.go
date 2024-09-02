package five

import (
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/strategy-gen/types"
	"github.com/tsinghua-cel/strategy-gen/utils"
	"time"
)

type Five struct {
}

func (o *Five) Description() string {
	//	desc_cn := `
	//提前一个epoch 查看下一个epoch的恶意节点出块顺序；
	//两个恶意节点之间穿插了一个诚实节点的块，让第二恶意节点的区块的parent，指向上一个恶意节点的slot.
	//delay 策略：第一个恶意节点的区块，广播delay 1个slot;
	desc_eng := `	Check the order of malicious nodes in the next epoch one epoch in advance;
	Two malicious nodes are interspersed with an honest node's block, so that the parent of the second malicious node's block points to the slot of the previous malicious node.
	Delay strategy: the block of the first malicious node broadcasts a delay of 1 slot.`
	return desc_eng
}

func (o *Five) Run(params types.LibraryParams) {
	log.WithField("name", "five").Info("start to run strategy")
	var latestEpoch int64
	ticker := time.NewTicker(time.Second * 3)
	slotTool := utils.SlotTool{SlotsPerEpoch: 32}
	for {
		select {
		case <-ticker.C:
			slot, err := utils.GetCurSlot(params.Attacker)
			if err != nil {
				log.WithField("error", err).Error("failed to get slot")
				continue
			}
			epoch := slotTool.SlotToEpoch(int64(slot))
			if epoch == latestEpoch {
				continue
			}
			if int64(slot) < slotTool.EpochEnd(epoch) {
				//continue
			}
			latestEpoch = epoch
			// get next epoch duties
			duties, err := utils.GetEpochDuties(params.Attacker, epoch+1)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"epoch": epoch + 1,
				}).Error("failed to get duties")
				latestEpoch = epoch - 1
				continue
			}
			if hackDuties, happen := CheckDuties(params.MaxValidatorIndex, duties); happen {
				strategy := types.Strategy{}
				strategy.Slots = GenSlotStrategy(hackDuties)
				if err = utils.UpdateStrategy(params.Attacker, strategy); err != nil {
					log.WithField("error", err).Error("failed to update strategy")
				} else {
					log.WithFields(log.Fields{
						"epoch":    epoch + 1,
						"strategy": strategy,
					}).Info("update strategy successfully")
				}
			}
		}
	}
}
