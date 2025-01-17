package feedback

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/common"
	"github.com/tsinghua-cel/attacker-service/config"
	"github.com/tsinghua-cel/attacker-service/strategy/slotstrategy"
	"github.com/tsinghua-cel/attacker-service/types"
	"sync/atomic"
)

type pairStrategy struct {
	uid      string
	origin   *types.Strategy
	parsed   []*slotstrategy.InternalSlotStrategy
	maxEpoch atomic.Value
	minEpoch atomic.Value
}

var (
	FOREVER = int64(1<<63 - 100)
)

// CalcEpochs calculate the min and max epoch of the strategy.
func (p *pairStrategy) CalcEpochs() (int64, int64) {
	var minEpoch, maxEpoch int64 = FOREVER, -1
	for _, s := range p.parsed {
		log.WithField("type", fmt.Sprintf("%T", s.Slot)).Debug("check slot type")
		switch s.Slot.(type) {
		case slotstrategy.NumberSlot:
			slot := s.Slot.(slotstrategy.NumberSlot)
			epoch := common.SlotToEpoch(int64(slot))
			if epoch > maxEpoch {
				maxEpoch = epoch
			}
			if epoch < minEpoch {
				minEpoch = epoch
			}
			log.WithFields(log.Fields{
				"minEpoch": minEpoch,
				"maxEpoch": maxEpoch,
				"epoch":    epoch,
			}).Debug("calc epoch")
		case slotstrategy.FunctionSlot:
			maxEpoch = FOREVER
			log.WithField("maxEpoch", "forever").Debug("set maxEpoch forever")

		default:
			// unknown slot type.
			log.Debug("unknown slot type")
		}
	}
	if minEpoch > maxEpoch {
		return maxEpoch, minEpoch
	}
	return minEpoch, maxEpoch
}

func (p *pairStrategy) IsEnd(epoch int64) bool {
	var maxEpoch int64
	if v := p.maxEpoch.Load(); v != nil {
		maxEpoch = v.(int64)
	} else {
		// calc max epoch
		mi, ma := p.CalcEpochs()
		p.maxEpoch.Store(ma)
		p.minEpoch.Store(mi)
		maxEpoch = ma
	}
	return epoch >= (maxEpoch + config.GetSafeEpochEndInterval())
}
