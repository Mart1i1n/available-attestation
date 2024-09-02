package one

import (
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/strategy-gen/types"
	"github.com/tsinghua-cel/strategy-gen/utils"
	"time"
)

type One struct {
}

func (o *One) Description() string {
	//	desc_cn := `
	//提前一个epoch 查看下一个epoch的恶意节点出块顺序；
	//如果有连续两个以上恶意节点出块，开始进行策略；
	//delay策略：blockdelay 到其中最后一个恶意节点出块的下一个slot；
	//恶意节点的投票者 开始做恶，对投票进行delay，执行的策略和blockdelay一样。`
	desc_eng := `	Check the order of malicious nodes in the next epoch one epoch in advance;
	If there are more than three malicious nodes in a row, start the strategy;
	Delay strategy: blockdelay to the next slot after the last malicious node;
	Malicious nodes' voters start to do evil, delay the voting, and execute the
	same strategy as blockdelay.`
	return desc_eng
}

func (o *One) Run(params types.LibraryParams) {
	log.WithField("name", "one").Info("start to run strategy")
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
				strategy.Validators = ValidatorStrategy(hackDuties)
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
