package three

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/strategy-gen/globalinfo"
	"github.com/tsinghua-cel/strategy-gen/types"
	"github.com/tsinghua-cel/strategy-gen/utils"
	"time"
)

type Three struct {
}

func (o *Three) Name() string {
	return "three"
}

func (o *Three) Description() string {
	desc_eng := `Unrealized justification reorg attack`
	return desc_eng
}

func (o *Three) Run(ctx context.Context, params types.LibraryParams) {
	log.WithField("name", "three").Info("start to run strategy")
	var latestEpoch int64 = -1
	ticker := time.NewTicker(time.Second * 3)
	slotTool := utils.SlotTool{SlotsPerEpoch: globalinfo.ChainBaseInfo().SlotsPerEpoch}
	for {
		select {
		case <-ctx.Done():
			log.WithField("name", "three").Info("stop to run strategy")
			return
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
				continue
			}
			latestEpoch = epoch
			// get next epoch duties
			nextEpoch := epoch + 1
			duties, err := utils.GetEpochDuties(params.Attacker, nextEpoch)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"epoch": nextEpoch,
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
						"epoch":    nextEpoch,
						"strategy": strategy,
					}).Info("update strategy successfully")
				}
			}
		}
	}
}
