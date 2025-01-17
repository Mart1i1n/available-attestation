package server

import (
	"context"
	"errors"
	"fmt"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	ethtype "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/groupcache/lru"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/beaconapi"
	"github.com/tsinghua-cel/attacker-service/common"
	"github.com/tsinghua-cel/attacker-service/config"
	"github.com/tsinghua-cel/attacker-service/dbmodel"
	"github.com/tsinghua-cel/attacker-service/feedback"
	"github.com/tsinghua-cel/attacker-service/openapi"
	"github.com/tsinghua-cel/attacker-service/plugins"
	"github.com/tsinghua-cel/attacker-service/rpc"
	"github.com/tsinghua-cel/attacker-service/server/apis"
	"github.com/tsinghua-cel/attacker-service/strategy"
	"github.com/tsinghua-cel/attacker-service/strategy/slotstrategy"
	"github.com/tsinghua-cel/attacker-service/types"
	"math/big"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	config       *config.Config
	rpcAPIs      []rpc.API   // List of APIs currently provided by the node
	http         *httpServer //
	strategy     *types.Strategy
	internal     []*slotstrategy.InternalSlotStrategy
	execClient   *ethclient.Client
	beaconClient *beaconapi.BeaconGwClient
	honestBeacon *beaconapi.BeaconGwClient

	validatorSetInfo *types.ValidatorDataSet
	mux              sync.Mutex
	attestpool       map[uint64]map[string]*ethpb.Attestation
	openApi          *openapi.OpenAPI
	cache            *lru.Cache
	hotdata          map[string]interface{}

	feedBacker      *feedback.Feedback
	historyStrategy *lru.Cache
	maxMaliciousIdx int
}

func (n *Server) GetBlockBySlot(slot uint64) (interface{}, error) {
	return n.beaconClient.GetDenebBlockBySlot(slot)
}

func (n *Server) GetLatestBeaconHeader() (types.BeaconHeaderInfo, error) {
	return n.beaconClient.GetLatestBeaconHeader()
}

func NewServer(conf *config.Config, plugin plugins.AttackerPlugin, maxMaliciousIdx int) *Server {
	s := &Server{}
	s.maxMaliciousIdx = maxMaliciousIdx
	s.cache = lru.New(10000)
	s.historyStrategy = lru.New(10000)
	s.config = conf
	s.rpcAPIs = apis.GetAPIs(s, plugin)
	client, err := ethclient.Dial(conf.ExecuteRpc)
	if err != nil {
		panic(fmt.Sprintf("dial execute failed with err:%v", err))
	}
	s.execClient = client
	s.beaconClient = beaconapi.NewBeaconGwClient(conf.BeaconRpc)
	s.honestBeacon = beaconapi.NewBeaconGwClient(conf.HonestBeaconRpc)
	s.http = newHTTPServer(log.WithField("module", "server"), rpc.DefaultHTTPTimeouts)
	s.strategy = strategy.ParseStrategy(s, conf.Strategy)
	s.validatorSetInfo = types.NewValidatorSet()
	s.attestpool = make(map[uint64]map[string]*ethpb.Attestation)
	s.openApi = openapi.NewOpenAPI(s, conf)
	s.hotdata = make(map[string]interface{})
	s.feedBacker = feedback.NewFeedback(s)
	return s
}

// startRPC is a helper method to configure all the various RPC endpoints during node
// startup. It's not meant to be called at any time afterwards as it makes certain
// assumptions about the state of the node.
func (n *Server) startRPC() error {
	// Filter out personal api
	var (
		servers []*httpServer
	)

	rpcConfig := rpcEndpointConfig{
		batchItemLimit:         config.APIBatchItemLimit,
		batchResponseSizeLimit: config.APIBatchResponseSizeLimit,
	}

	initHttp := func(server *httpServer, port int) error {
		if err := server.setListenAddr(n.config.HttpHost, port); err != nil {
			return err
		}
		if err := server.enableRPC(n.rpcAPIs, httpConfig{
			CorsAllowedOrigins: config.DefaultCors,
			Vhosts:             config.DefaultVhosts,
			Modules:            config.DefaultModules,
			prefix:             config.DefaultPrefix,
			rpcEndpointConfig:  rpcConfig,
		}); err != nil {
			return err
		}
		servers = append(servers, server)
		return nil
	}

	// Set up HTTP.
	// Configure legacy unauthenticated HTTP.
	if err := initHttp(n.http, n.config.HttpPort); err != nil {
		return err
	}

	// Start the servers
	for _, server := range servers {
		if err := server.start(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) monitorEvent() {
	ticker := time.NewTicker(time.Minute * 2)
	defer ticker.Stop()
	totalReorgDepth := uint64(0)

	handler := func(ch chan *apiv1.ChainReorgEvent) {
		for {
			select {
			case reorg := <-ch:
				log.WithFields(log.Fields{
					"slot":            reorg.Slot,
					"depth":           reorg.Depth,
					"totalReorgDepth": totalReorgDepth,
				}).Info("reorg event")
				totalReorgDepth += reorg.Depth
				ev := types.ReorgEvent{
					Epoch:        int64(reorg.Epoch),
					Slot:         int64(reorg.Slot),
					Depth:        int64(reorg.Depth),
					OldHeadState: reorg.OldHeadState.String(),
					NewHeadState: reorg.NewHeadState.String(),
				}
				if oldHeader, err := s.honestBeacon.GetBlockHeaderById(reorg.OldHeadBlock.String()); err == nil {
					ev.OldBlockSlot = int64(oldHeader.Header.Message.Slot)
					ev.OldBlockProposerIndex = int64(oldHeader.Header.Message.ProposerIndex)
				}
				if newHeader, err := s.honestBeacon.GetBlockHeaderById(reorg.NewHeadBlock.String()); err == nil {
					ev.NewBlockSlot = int64(newHeader.Header.Message.Slot)
					ev.NewBlockProposerIndex = int64(newHeader.Header.Message.ProposerIndex)
				}
				dbmodel.InsertNewReorg(ev)
			}
		}
	}

	for {
		select {
		case <-ticker.C:
			eventCh := s.honestBeacon.MonitorReorgEvent()
			if eventCh != nil {
				go handler(eventCh)
				ticker.Reset(time.Hour * 256)
			} else {
				ticker.Reset(time.Minute)
			}
		}
	}

}

func (s *Server) monitorDuties() {
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	dutyTicker := time.NewTicker(time.Minute)
	defer dutyTicker.Stop()

	slotsPerEpoch := int64(32)
	dumped := make(map[int64]bool)

	for {
		select {

		case <-dutyTicker.C:
			header, err := s.honestBeacon.GetLatestBeaconHeader()
			if err != nil {
				log.WithError(err).Debug("duty ticker get latest beacon header failed")
				continue
			}
			curSlot, _ := strconv.ParseInt(header.Header.Message.Slot, 10, 64)
			curEpoch := curSlot / slotsPerEpoch
			nextEpoch := curEpoch + 1
			if curEpoch == 0 && dumped[curEpoch] == false {
				//
				if err := s.dumpDuties(curEpoch); err == nil {
					dumped[curEpoch] = true
				}
			}
			if dumped[nextEpoch] == false {
				if err := s.dumpDuties(nextEpoch); err == nil {
					dumped[nextEpoch] = true
				}
			}

		case <-ticker.C:
			curDuties, err := s.honestBeacon.GetCurrentEpochAttestDuties()
			if err != nil {
				continue
			}
			for _, duty := range curDuties {
				if idx, err := strconv.Atoi(duty.ValidatorIndex); err == nil {
					s.validatorSetInfo.AddValidator(idx, duty.Pubkey)
				}
			}
			nextDuties, _ := s.honestBeacon.GetNextEpochAttestDuties()
			for _, duty := range nextDuties {
				if idx, err := strconv.Atoi(duty.ValidatorIndex); err == nil {
					s.validatorSetInfo.AddValidator(idx, duty.Pubkey)
				}
			}

			ticker.Reset(time.Second * 2)

		}
	}
}

func (s *Server) Start() {
	s.initTools()
	// start RPC endpoints
	err := s.startRPC()
	if err != nil {
		s.stopRPC()
	}
	s.openApi.Start()
	// start collect duties info.
	go s.monitorDuties()
	go s.monitorEvent()
	go s.HandleEndStrategy()
	s.feedBacker.Start()
}

func (s *Server) initTools() {
	init := false
	for !init {
		slotPerEpoch, _ := s.honestBeacon.GetIntConfig(beaconapi.SLOTS_PER_EPOCH)
		interval, _ := s.honestBeacon.GetIntConfig(beaconapi.SECONDS_PER_SLOT)
		genesis, err := s.honestBeacon.GetGenesis()
		if slotPerEpoch == 0 || interval == 0 || err != nil {
			log.WithError(err).Error("initTools get genesis failed, retry")
			time.Sleep(time.Second)
			continue
		} else {
			common.InitSlotTool(interval, int64(slotPerEpoch), genesis.GenesisTime.Unix())
			init = true
			log.WithField("interval", interval).Info("init tool finished")
		}
	}
}

func (s *Server) stopRPC() {
	s.http.stop()
}

// implement backend
func (s *Server) SomeNeedBackend() bool {
	return true
}

func (s *Server) GetBlockHeight() (uint64, error) {
	return s.execClient.BlockNumber(context.Background())
}

func (s *Server) GetBlockByNumber(number *big.Int) (*ethtype.Block, error) {
	return s.execClient.BlockByNumber(context.Background(), number)
}

func (s *Server) GetHeightByNumber(number *big.Int) (*ethtype.Header, error) {
	return s.execClient.HeaderByNumber(context.Background(), number)
}

func (s *Server) GetStrategy() *types.Strategy {
	return s.strategy
}

func (s *Server) GetValidatorRoleByPubkey(slot int, pubkey string) types.RoleType {
	if val := s.validatorSetInfo.GetValidatorByPubkey(pubkey); val != nil {
		return s.GetValidatorRole(slot, int(val.Index))
	} else {
		return types.NormalRole
	}
}

func (s *Server) GetCurrentEpochProposeDuties() ([]types.ProposerDuty, error) {
	return s.honestBeacon.GetCurrentEpochProposerDuties()
}

func (s *Server) GetCurrentEpochAttestDuties() ([]types.AttestDuty, error) {
	return s.honestBeacon.GetCurrentEpochAttestDuties()
}

func (s *Server) GetSlotsPerEpoch() int {
	count, err := s.honestBeacon.GetIntConfig(beaconapi.SLOTS_PER_EPOCH)
	if err != nil {
		return 6
	}
	return count
}

func (s *Server) GetIntervalPerSlot() int {
	interval, _ := s.honestBeacon.GetIntConfig(beaconapi.SECONDS_PER_SLOT)
	return interval
}

func (s *Server) AddSignedAttestation(slot uint64, pubkey string, attestation *ethpb.Attestation) {
	s.validatorSetInfo.AddSignedAttestation(slot, pubkey, attestation)
}

func (s *Server) AddAttestToPool(slot uint64, pubkey string, attestation *ethpb.Attestation) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if _, ok := s.attestpool[slot]; !ok {
		s.attestpool[slot] = make(map[string]*ethpb.Attestation)
	}
	s.attestpool[slot][pubkey] = attestation
}

func (s *Server) GetAttestPool() map[uint64]map[string]*ethpb.Attestation {
	// copy
	s.mux.Lock()
	defer s.mux.Unlock()
	data := make(map[uint64]map[string]*ethpb.Attestation)
	for k, v := range s.attestpool {
		data[k] = make(map[string]*ethpb.Attestation)
		for k1, v1 := range v {
			data[k][k1] = v1
		}
	}
	return data
}

func (s *Server) ResetAttestPool() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.attestpool = make(map[uint64]map[string]*ethpb.Attestation)
}

func (s *Server) AddSignedBlock(slot uint64, pubkey string, block *ethpb.GenericSignedBeaconBlock) {
	s.validatorSetInfo.AddSignedBlock(slot, pubkey, block)
}

func (s *Server) GetAttestSet(slot uint64) *types.SlotAttestSet {
	return s.validatorSetInfo.GetAttestSet(slot)
}

func (s *Server) GetBlockSet(slot uint64) *types.SlotBlockSet {
	return s.validatorSetInfo.GetBlockSet(slot)
}

func (s *Server) GetValidatorDataSet() *types.ValidatorDataSet {
	return s.validatorSetInfo
}

func (s *Server) GetValidatorByProposeSlot(slot uint64) (int, error) {
	epochPerSlot := uint64(s.GetSlotsPerEpoch())
	epoch := slot / epochPerSlot
	duties, err := s.honestBeacon.GetProposerDuties(int(epoch))
	if err != nil {
		return 0, err
	}
	for _, duty := range duties {
		dutySlot, _ := strconv.ParseInt(duty.Slot, 10, 64)
		if uint64(dutySlot) == slot {
			idx, _ := strconv.Atoi(duty.ValidatorIndex)
			return idx, nil
		}
	}
	return 0, errors.New("not found")
}

func (s *Server) GetProposeDuties(epoch int) ([]types.ProposerDuty, error) {
	return s.honestBeacon.GetProposerDuties(epoch)
}

func (s *Server) SlotsPerEpoch() int {
	return s.GetSlotsPerEpoch()
}

func (s *Server) GetValidatorRole(slot int, valIdx int) types.RoleType {
	if slot < 0 {
		header, err := s.beaconClient.GetLatestBeaconHeader()
		if err != nil {
			return types.NormalRole
		}
		slot, _ = strconv.Atoi(header.Header.Message.Slot)
	}
	return s.strategy.GetValidatorRole(valIdx, int64(slot))
}

func (s *Server) GetInternalSlotStrategy() []*slotstrategy.InternalSlotStrategy {
	var err error
	if len(s.internal) == 0 {
		s.internal, err = slotstrategy.ParseToInternalSlotStrategy(s, s.strategy.Slots)
		if err != nil {
			log.WithError(err).Error("parse strategy failed")
			return nil
		}
	}
	return s.internal
}
func (s *Server) GetSlotRoot(slot int64) (string, error) {
	return s.beaconClient.GetSlotRoot(slot)
}

func (s *Server) dumpDuties(epoch int64) error {
	duties, err := s.GetProposeDuties(int(epoch))
	if err != nil {
		return err
	}
	for _, duty := range duties {
		log.WithFields(log.Fields{
			"epoch":     epoch,
			"slot":      duty.Slot,
			"validator": duty.ValidatorIndex,
		}).Info("epoch duty")
	}
	return nil
}

// UpdateStrategy only update strategy slots and actions, not update validators.
func (s *Server) UpdateStrategy(strategy *types.Strategy) error {
	check := false
	if strategy.Uid != "" {
		check = true
	}
	log.WithField("uid", strategy.Uid).Debug("goto parse and update strategy")

	if check {
		if st := dbmodel.GetStrategyByUUID(strategy.Uid); st != nil {
			log.WithError(errors.New("strategy already exist")).Error("strategy already exist")
			return errors.New("strategy already exist")
		}
	}

	parsed, err := slotstrategy.ParseToInternalSlotStrategy(s, strategy.Slots)
	if err != nil {
		log.WithError(err).Error("parse strategy failed")
		return err
	}
	for _, v := range parsed {
		replaced := false
		for _, vi := range s.internal {
			if vi.Slot.StrValue() == v.Slot.StrValue() && vi.Level <= v.Level {
				// replace actions
				vi.Actions = v.Actions
				replaced = true
				break
			}
		}
		if !replaced {
			s.internal = append(s.internal, v)
		}
	}
	// dump internal strategy
	for _, v := range s.internal {
		log.WithFields(log.Fields{
			"slot":  v.Slot.StrValue(),
			"level": v.Level,
		}).Debug("internal strategy slot")
		for k, action := range v.Actions {
			log.WithFields(log.Fields{
				"slot":       v.Slot.StrValue(),
				"level":      v.Level,
				"checkpoint": k,
				"action":     action.Name(),
			}).Debug("internal strategy action")

		}
	}

	for _, v := range strategy.Slots {
		replaced := false
		for _, vi := range s.strategy.Slots {
			if v.Slot == vi.Slot && vi.Level <= v.Level {
				vi.Actions = v.Actions
				replaced = true
				break
			}
		}
		if !replaced {
			s.strategy.Slots = append(s.strategy.Slots, v)
		}
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	// not update validators, you can set it at initial time.
	//s.strategy.Validators = strategy.Validators
	log.WithFields(log.Fields{
		"strategy": s.strategy,
		"check":    check,
	}).Debug("goto check strategy")

	dbmodel.InsertNewStrategy(strategy)

	if check {
		s.historyStrategy.Add(strategy.Uid, HistoryStrategy{
			Strategy: strategy,
		})
		s.feedBacker.AddNewStrategy(strategy.Uid, strategy, parsed)
	}

	return nil
}

func (s *Server) GetSlotStartTime(slot int) (int64, bool) {
	key := fmt.Sprintf("slot_start_time_%d", slot)
	if v, ok := s.cache.Get(key); ok {
		return v.(int64), true
	}
	return 0, false
}

func (s *Server) SetSlotStartTime(slot int, time int64) {
	key := fmt.Sprintf("slot_start_time_%d", slot)
	if _, ok := s.cache.Get(key); !ok {
		s.cache.Add(key, time)
	}
}

func (s *Server) GetCurSlot() int64 {
	key := fmt.Sprintf("cur_slot")
	s.mux.Lock()
	defer s.mux.Unlock()
	if v, ok := s.hotdata[key]; ok {
		return v.(int64)
	}
	return 0
}

func (s *Server) SetCurSlot(slot int64) {
	key := fmt.Sprintf("cur_slot")
	s.mux.Lock()
	defer s.mux.Unlock()
	if v, ok := s.hotdata[key]; !ok {
		s.hotdata[key] = slot
	} else {
		if v.(int64) < slot {
			s.hotdata[key] = slot
		}
	}
}

func (s *Server) HandleEndStrategy() {
	ch := make(chan feedback.StrategyEndEvent, 10)
	sub := s.feedBacker.SubscribeStrategyEndEvent(ch)
	if sub == nil {
		log.Error("subscribe strategy end event failed")
		return
	}
	defer sub.Unsubscribe()
	for {
		select {
		case ev := <-ch:
			log.WithFields(log.Fields{
				"uid": ev.Uid,
			}).Debug("got strategy end event")
			uid := ev.Uid
			if v, exist := s.historyStrategy.Get(uid); exist {
				storeStrategy := dbmodel.GetStrategyByUUID(uid)
				historyInfo := v.(HistoryStrategy)
				// get reorg event count.
				totalReorgCount := 0
				totalImpactCount := 0
				normalTargetAmount := int64(290680)
				//normalHeadAmount := int64(156520)
				finalHonestLoseRate := float64(0.0)
				finalAttackerLoseRate := float64(0.0)
				for i := ev.MinEpoch; i <= ev.MaxEpoch; i++ {
					reorgCount := dbmodel.GetReorgCountByEpoch(i)
					impactCount := dbmodel.GetImpactValidatorCount(s.maxMaliciousIdx, normalTargetAmount, i)
					log.WithFields(log.Fields{
						"epoch":  i,
						"reorg":  reorgCount,
						"impact": impactCount,
					}).Debug("strategy feedback")
					honestLoseRate, attackerLoseRate := calcLoseRate(dbmodel.GetRewardListByEpoch(i), normalTargetAmount, s.maxMaliciousIdx)
					finalHonestLoseRate += honestLoseRate
					finalAttackerLoseRate += attackerLoseRate
					totalReorgCount += reorgCount
					totalImpactCount += impactCount
				}
				historyInfo.FeedBackInfo = &types.FeedBackInfo{
					ReorgCount:           totalReorgCount,
					ImpactValidatorCount: totalImpactCount,
				}
				storeStrategy.MinEpoch = ev.MinEpoch
				storeStrategy.MaxEpoch = ev.MaxEpoch
				storeStrategy.ReorgCount = totalReorgCount
				storeStrategy.ImpactValidatorCount = totalImpactCount
				storeStrategy.IsEnd = true
				storeStrategy.HonestLoseRateAvg = finalHonestLoseRate / float64(ev.MaxEpoch-ev.MinEpoch+1)
				storeStrategy.AttackerLoseRateAvg = finalAttackerLoseRate / float64(ev.MaxEpoch-ev.MinEpoch+1)
				dbmodel.StrategyUpdate(storeStrategy)

				s.historyStrategy.Add(uid, historyInfo)
				log.WithFields(log.Fields{
					"uid":  uid,
					"info": historyInfo,
				}).Debug("update strategy feedback info")
			}
		}

	}
}

func (s *Server) GetFeedBack(uid string) (types.FeedBackInfo, error) {
	v, exist := s.historyStrategy.Get(uid)
	if exist {
		historyInfo := v.(HistoryStrategy)
		if historyInfo.FeedBackInfo != nil {
			return *historyInfo.FeedBackInfo, nil
		} else {
			return types.FeedBackInfo{}, errors.New("strategy feedback not generated")
		}
	}
	if st := dbmodel.GetStrategyByUUID(uid); st != nil && st.IsEnd {
		return types.FeedBackInfo{
			st.ReorgCount,
			st.ImpactValidatorCount,
		}, nil

	} else {
		return types.FeedBackInfo{}, errors.New("strategy not found or not finished")
	}

}

// calcLoseRate return honestLoseRate and attackerLoseRate.
func calcLoseRate(rewards []*dbmodel.AttestReward, normalTargetAmount int64, maxMaliciousIdx int) (float64, float64) {
	honestLoseRate := float64(0.0)
	honestCount := 0
	attackerLoseRate := float64(0.0)
	attackerCount := 0
	// calc every validator loseRate and get average.
	// loseRate := (normalTargetAmount - reward.TargetAmount)/normalTargetAmount
	// and if rewards.ValidatorIndex <= maxMaliciousIdx, it is attacker.
	for _, reward := range rewards {
		loseRate := float64(normalTargetAmount-reward.TargetAmount) / float64(normalTargetAmount)
		if reward.ValidatorIndex <= maxMaliciousIdx {
			attackerLoseRate += loseRate
			attackerCount++
		} else {
			honestLoseRate += loseRate
			honestCount++
		}
	}
	if honestCount > 0 {
		honestLoseRate = honestLoseRate / float64(honestCount)
	}
	if attackerCount > 0 {
		attackerLoseRate = attackerLoseRate / float64(attackerCount)
	}
	return honestLoseRate, attackerLoseRate
}
