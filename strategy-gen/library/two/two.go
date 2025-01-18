package two

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/strategy-gen/globalinfo"
	"github.com/tsinghua-cel/strategy-gen/types"
	"github.com/tsinghua-cel/strategy-gen/utils"
	"strconv"
	"time"
)

type Two struct {
}

func (o *Two) Name() string {
	return "two"
}

func (o *Two) Description() string {
	desc_eng := `Staircase attack`
	return desc_eng
}

func (o *Two) Run(ctx context.Context, params types.LibraryParams) {
	log.WithField("name", o.Name()).Info("start to run strategy")
	var latestEpoch int64 = -1
	ticker := time.NewTicker(time.Second * 3)
	slotTool := utils.SlotTool{SlotsPerEpoch: globalinfo.ChainBaseInfo().SlotsPerEpoch}
	for {
		select {
		case <-ctx.Done():
			log.WithField("name", o.Name()).Info("stop to run strategy")
			return
		case <-ticker.C:
			slot, err := utils.GetCurSlot(params.Attacker)
			if err != nil {
				log.WithField("error", err).Error("failed to get slot")
				continue
			}
			log.WithFields(log.Fields{
				"slot":      slot,
				"lastEpoch": latestEpoch,
			}).Info("get slot")
			epoch := slotTool.SlotToEpoch(int64(slot))
			// generate new strategy at the end of last epoch.
			if int64(slot) < slotTool.EpochEnd(epoch) {
				continue
			}
			if epoch == latestEpoch {
				continue
			}
			latestEpoch = epoch

			{
				nextEpoch := epoch + 1
				cas := 0

				duties, err := utils.GetEpochDuties(params.Attacker, nextEpoch)
				pduties, err := utils.GetEpochDuties(params.Attacker, epoch-1)
				cduties, err := utils.GetEpochDuties(params.Attacker, epoch)
				if err != nil {
					log.WithFields(log.Fields{
						"error": err,
						"epoch": epoch + 1,
					}).Error("failed to get duties")
					latestEpoch = epoch - 1
					continue
				}
				strategy := types.Strategy{}
				strategy.Validators = ValidatorStrategy(params.MaxValidatorIndex, nextEpoch)
				if checkFirstByzSlot(pduties, params.MaxValidatorIndex) && checkFirstByzSlot(cduties, params.MaxValidatorIndex) && !checkFirstByzSlot(duties, params.MaxValidatorIndex) {
					cas = 1
				}
				strategy.Slots = GenSlotStrategy(getLatestHackerSlot(duties, params.MaxValidatorIndex), nextEpoch, cas)
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

func getLatestHackerSlot(duties []types.ProposerDuty, maxValidatorIndex int) int {
	latest, _ := strconv.Atoi(duties[0].Slot)
	for _, duty := range duties {
		idx, _ := strconv.Atoi(duty.ValidatorIndex)
		slot, _ := strconv.Atoi(duty.Slot)
		if idx > maxValidatorIndex {
			continue
		}
		if slot > latest {
			latest = slot
		}
	}
	return latest

}

func checkFirstByzSlot(duties []types.ProposerDuty, maxValidatorIndex int) bool {
	firstproposerindex, _ := strconv.Atoi(duties[0].ValidatorIndex)
	if firstproposerindex > maxValidatorIndex {
		return false
	}
	return true
}
