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
	//Check the blocking sequence of malicious nodes one epoch ahead of the next epoch；
	//If more than two malicious nodes create blocks in a row, the policy is started；
	//Strategy: Blockdelay to the next slot where the last malicious node exits the block；
	desc_eng := `Exante reorg attack.`
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
