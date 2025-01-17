package feedback

import (
	"github.com/ethereum/go-ethereum/event"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/dbmodel"
	"github.com/tsinghua-cel/attacker-service/strategy/slotstrategy"
	"github.com/tsinghua-cel/attacker-service/types"
	"sync"
	"time"
)

type Feedback struct {
	historyStrategy map[int64]*pairStrategy
	quit            chan struct{}
	mux             sync.Mutex
	feed            event.Feed
	backend         types.CacheBackend
}

func NewFeedback(backend types.CacheBackend) *Feedback {
	return &Feedback{
		historyStrategy: make(map[int64]*pairStrategy),
		quit:            make(chan struct{}),
		backend:         backend,
	}
}

func (f *Feedback) AddNewStrategy(uid string, origin *types.Strategy, parsed []*slotstrategy.InternalSlotStrategy) {
	f.mux.Lock()
	defer f.mux.Unlock()

	f.historyStrategy[time.Now().UnixMilli()] = &pairStrategy{uid: uid, origin: origin, parsed: parsed}
}

func (f *Feedback) Start() {
	go f.loop()
}

func (f *Feedback) Stop() {
	close(f.quit)
}

type StrategyEndEvent struct {
	Uid      string
	MinEpoch int64
	MaxEpoch int64
}

func (f *Feedback) SubscribeStrategyEndEvent(ch chan StrategyEndEvent) event.Subscription {
	return f.feed.Subscribe(ch)
}

func (f *Feedback) loop() {
	tc := time.NewTicker(time.Second * 10)
	defer tc.Stop()
	for {
		select {
		case <-tc.C:
			f.mux.Lock()
			for timestamp, pair := range f.historyStrategy {
				safeEpoch := dbmodel.GetMaxEpoch()
				ended := pair.IsEnd(safeEpoch)
				log.WithFields(log.Fields{
					"strategy":  pair.origin,
					"safeEpoch": safeEpoch,
				}).Info("check strategy end")
				if ended {
					ev := StrategyEndEvent{
						Uid:      pair.uid,
						MinEpoch: pair.minEpoch.Load().(int64),
						MaxEpoch: pair.maxEpoch.Load().(int64),
					}
					// send event
					f.feed.Send(ev)
					delete(f.historyStrategy, timestamp)
					log.WithFields(log.Fields{
						"uid":       ev.Uid,
						"minEpoch":  ev.MinEpoch,
						"maxEpoch":  ev.MaxEpoch,
						"safeEpoch": safeEpoch,
					}).Info("post strategy end event")
				}
			}
			f.mux.Unlock()
		case <-f.quit:
			return
		}

	}

}
