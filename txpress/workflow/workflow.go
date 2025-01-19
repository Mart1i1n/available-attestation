package workflow

import (
	log "github.com/sirupsen/logrus"
	"github.com/xueqianLu/txpress/chains"
	"github.com/xueqianLu/txpress/types"
	"sync"
	"time"
)

type Task struct {
	baseCount int
	batch     int
	interval  time.Duration
}

type Result struct {
	chain    string
	minBlock int
	maxBlock int
}

type Workflow struct {
	chains []types.ChainPlugin
	conf   types.RunConfig
	quit   chan struct{}
}

func NewWorkFlow(chains []types.ChainPlugin, conf types.RunConfig) *Workflow {
	return &Workflow{
		chains: chains,
		conf:   conf,
		quit:   make(chan struct{}),
	}
}

func (w *Workflow) wait() {
	for {
		start := false
		for _, chain := range w.chains {
			block, err := chain.LatestBlockInfo()
			if err != nil || block.Number < int64(w.conf.BeginToStart) {
				continue
			}
			start = true
			break
		}
		if start {
			break
		} else {
			log.Info("wait when latest block.number > ", w.conf.BeginToStart)
			time.Sleep(time.Second * 3)
		}
	}
	log.Info("wait finished")
}

func (w *Workflow) Start() {
	wg := sync.WaitGroup{}
	tss := make([]*txStream, len(w.chains))
	for i, chain := range w.chains {
		tss[i] = newTxStream(chain, w.conf, w.quit)
		wg.Add(1)
		go func(ts *txStream) {
			ts.pause(true)
			ts.start()
			wg.Done()
		}(tss[i])
	}

	lastTps := 0
	baseTxCount := w.conf.BaseCount
	history := make([]types.Record, 0)
	// wait start, if latest block benefit is "", then start the test.
	w.wait()

	for r := 0; r < w.conf.Round; r++ {
		for _, ts := range tss {
			ts.pause(false)
		}
		beginTime := time.Now()
		begin, e := w.chains[0].LatestBlockInfo()
		if e != nil {
			log.Error("get latest block info failed")
			continue
		}
		log.WithFields(log.Fields{
			"begin": begin.Number,
			"wait":  w.conf.Interval * time.Duration(w.conf.Batch) / time.Second,
		}).Info("test one round begin")
		time.Sleep(w.conf.Interval * time.Duration(w.conf.Batch))
		end, e := w.chains[0].LatestBlockInfo()
		if e != nil {
			log.Error("get latest block info failed")
			continue
		}
		// pause all stream
		for _, ts := range tss {
			ts.pause(true)
		}
		endTime := time.Now()

		log.Info("test one round end")

		// calculate tps
		rate := float64(baseTxCount) / w.conf.Interval.Seconds()
		record := w.calculateTps(rate, w.chains[0], int(begin.Number)+1, int(end.Number), endTime.Sub(beginTime))
		if record.Tps > 0 && record.Tps >= lastTps || w.conf.ForceIncrease {
			incs := baseTxCount * w.conf.IncRate / 100
			baseTxCount += incs
			lastTps = record.Tps
		}

		history = append(history, record)
		log.WithFields(log.Fields{
			"begin":     record.Begin,
			"end":       record.End,
			"totaltime": record.TotalTime,
			"totaltx":   record.TotalTx,
			"tps":       record.Tps,
		}).Info("test one round finished")
		newConf := w.conf
		newConf.BaseCount = baseTxCount
		for _, ts := range tss {
			ts.updateSpeed(newConf)
		}
		// wait for next round

	}

	close(w.quit)

	wg.Wait()
	for _, record := range history {
		log.WithFields(log.Fields{
			"begin":     record.Begin,
			"end":       record.End,
			"totaltime": record.TotalTime,
			"totaltx":   record.TotalTx,
			"tps":       record.Tps,
			"tx rate":   record.Rate,
		}).Info("test history")
	}
}

func (w *Workflow) calculateTps(rate float64, chain types.ChainPlugin, minBlock, maxBlock int, duration time.Duration) types.Record {
	return chains.CalcTps(rate, chain, minBlock, maxBlock, duration)
}

func (w *Workflow) makeTx(chain types.ChainPlugin, baseCount int, batch int, checkNonce bool) [][]types.ChainTx {
	txs := make([][]types.ChainTx, batch)
	for i := 0; i < batch; i++ {
		if i > 0 {
			checkNonce = false
		}
		mtxs, err := chain.CreateTxs(baseCount, checkNonce)
		if err != nil {
			return nil
		}
		txs[i] = mtxs
	}
	return txs
}

func (w *Workflow) loop(chain types.ChainPlugin, taskCh chan Task, result chan Result, wg *sync.WaitGroup) {

	defer wg.Done()
	first := true

	for {
		select {
		case task := <-taskCh:
			txs := w.makeTx(chain, task.baseCount, task.batch, first)
			_min, _max := w.runTest(chain, txs, task.interval)
			result <- Result{
				chain:    chain.Id(),
				minBlock: _min,
				maxBlock: _max,
			}
			if first {
				first = false
			}

		case <-w.quit:
			return
		}
	}
}

func (w *Workflow) runTest(chain types.ChainPlugin, txs [][]types.ChainTx, interval time.Duration) (int, int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var waited = make(map[string]int)
	hashes := make([][]string, 0)
	for _, batch := range txs {
		hashs, err := chain.SendTxs(batch)
		if err != nil {
			log.Errorf("send txs error: %s", err)
		} else {
			log.Infof("send txs success, count: %d", len(hashs))
			hashes = append(hashes, hashs)
			for _, hash := range hashs {
				waited[hash] = 0
			}
		}
		<-ticker.C
	}
	var _min, _max int
	for len(waited) > 0 {
		time.Sleep(time.Second)
		for hash, _ := range waited {
			time.Sleep(time.Millisecond * 10)
			block, err := chain.TxBlock(hash)
			if err != nil {
				continue
			}
			if block == 0 {
				continue
			}
			if _min == 0 {
				_min = block
			}

			if block < _min {
				_min = block
			}

			if _max == 0 {
				_max = block
			}

			if block > _max {
				_max = block
			}
			delete(waited, hash)
		}
	}

	return _min, _max
}
