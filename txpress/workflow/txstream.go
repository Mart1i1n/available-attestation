package workflow

import (
	log "github.com/sirupsen/logrus"
	"github.com/xueqianLu/txpress/types"
	"time"
)

type txStream struct {
	chain   types.ChainPlugin
	run     types.RunConfig
	paused  bool
	newConf chan types.RunConfig
	quit    chan struct{}
}

func newTxStream(chain types.ChainPlugin, run types.RunConfig, quit chan struct{}) *txStream {
	return &txStream{
		chain:   chain,
		run:     run,
		paused:  false,
		quit:    quit,
		newConf: make(chan types.RunConfig, 1),
	}
}

func (ts *txStream) start() {
	ticker := time.NewTicker(ts.run.Interval)
	defer ticker.Stop()
	for {
		select {
		case conf := <-ts.newConf:
			ticker.Reset(conf.Interval)
		case <-ts.quit:
			return
		case <-ticker.C:
			if ts.paused {
				continue
			}
			// make tx
			txs, err := ts.chain.CreateTxs(ts.run.BaseCount, false)
			if err != nil {
				log.WithError(err).Error("create tx failed")
				continue
			}
			if _, err := ts.chain.SendTxs(txs); err != nil {
				log.WithError(err).Error("send tx failed")
			} else {
				//log.WithFields(log.Fields{
				//	"chain": ts.chain.Id(),
				//	"count": len(hashs),
				//}).Info("send tx success")
			}
			ticker.Reset(ts.run.Interval)
		}
	}
}

func (ts *txStream) pause(value bool) {
	ts.paused = value
}

func (ts *txStream) updateSpeed(conf types.RunConfig) {
	ts.run = conf
	ts.newConf <- conf
}
