package slotstrategy

import (
	"fmt"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	attaggregation "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1/attestation/aggregation/attestations"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/common"
	"github.com/tsinghua-cel/attacker-service/plugins"
	"github.com/tsinghua-cel/attacker-service/types"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type ActionDo func(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse

type FunctionAction struct {
	doFunc ActionDo
	name   string
}

func (f FunctionAction) RunAction(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse {
	if f.doFunc != nil {
		return f.doFunc(backend, slot, pubkey, params...)
	}
	return plugins.PluginResponse{
		Cmd: types.CMD_NULL,
	}
}

func (f FunctionAction) Name() string {
	return f.name
}

func getCmdFromName(name string) types.AttackerCommand {
	switch name {
	case "null":
		return types.CMD_NULL
	case "return":
		return types.CMD_RETURN
	case "continue":
		return types.CMD_CONTINUE
	case "abort":
		return types.CMD_ABORT
	case "skip":
		return types.CMD_SKIP
	case "exit":
		return types.CMD_EXIT
	default:
		return types.CMD_NULL
	}
}

func ParseActionName(action string) (string, []int) {
	strs := strings.Split(action, ":")
	params := make([]int, 0)
	if len(strs) > 1 {
		for _, v := range strs[1:] {
			val, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			params = append(params, val)
		}
	}
	return strs[0], params
}

func GetFunctionAction(backend types.ServiceBackend, action string) (ActionDo, error) {
	name, params := ParseActionName(action)
	switch name {
	case "null", "return", "continue", "abort", "skip", "exit":
		cmd := getCmdFromName(name)
		return func(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse {
			log.WithFields(log.Fields{
				"slot":   slot,
				"action": name,
			}).Debug("do action ")
			r := plugins.PluginResponse{
				Cmd: cmd,
			}
			if len(params) > 0 {
				r.Result = params[0]
			}
			return r
		}, nil
	case "addAttestToPool":
		return func(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse {
			var attestation *ethpb.Attestation
			r := plugins.PluginResponse{
				Cmd: types.CMD_NULL,
			}
			log.WithFields(log.Fields{
				"slot":   slot,
				"action": name,
			}).Debug("do action ")

			if len(params) > 0 {
				attestation = params[0].(*ethpb.Attestation)
				backend.AddAttestToPool(uint64(slot), pubkey, attestation)
				r.Result = attestation
			}

			return r
		}, nil
	case "storeSignedAttest":
		return func(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse {
			var attestation *ethpb.Attestation
			r := plugins.PluginResponse{
				Cmd: types.CMD_NULL,
			}
			log.WithFields(log.Fields{
				"slot":   slot,
				"action": name,
			}).Debug("do action ")

			if len(params) > 0 {
				attestation = params[0].(*ethpb.Attestation)
				backend.AddSignedAttestation(uint64(slot), pubkey, attestation)
				r.Result = attestation
			}

			return r
		}, nil
	case "delayWithSecond", "delayWithDuration":
		var seconds int
		if len(params) == 0 {
			seconds = rand.Intn(10)
		} else {
			seconds = params[0]
		}

		return func(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse {
			r := plugins.PluginResponse{
				Cmd: types.CMD_NULL,
			}
			log.WithFields(log.Fields{
				"slot":   slot,
				"action": name,
			}).Debug("do action ")

			log.WithFields(log.Fields{
				"slot":    slot,
				"seconds": seconds,
			}).Debug("delayWithSecond")
			time.Sleep(time.Second * time.Duration(seconds))
			return r
		}, nil
	case "delayToNextSlot":
		return func(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse {
			r := plugins.PluginResponse{
				Cmd: types.CMD_NULL,
			}
			log.WithFields(log.Fields{
				"slot":   slot,
				"action": name,
			}).Debug("do action ")
			targetTime := common.TimeToSlot(slot + 1)
			total := targetTime - time.Now().Unix()
			log.WithFields(log.Fields{
				"slot":  slot,
				"total": total,
			}).Debug("delayToNextSlot")
			time.Sleep(time.Second * time.Duration(total))
			return r
		}, nil
	case "delayToAfterNextSlot":
		afters := rand.Intn(10)
		if len(params) > 0 {
			afters = params[0]
		}
		return func(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse {
			r := plugins.PluginResponse{
				Cmd: types.CMD_NULL,
			}
			log.WithFields(log.Fields{
				"slot":   slot,
				"action": name,
			}).Info("do action ")
			targetTime := common.TimeToSlot(slot + 1)
			targetTime += int64(afters)
			total := targetTime - time.Now().Unix()
			log.WithFields(log.Fields{
				"slot":  slot,
				"total": total,
			}).Debug("delayToAfterNextSlot")
			time.Sleep(time.Second * time.Duration(total))
			return r
		}, nil
	case "delayToNextNEpochStart":
		n := 1
		if len(params) > 0 {
			n = params[0]
		}
		return func(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse {
			log.WithFields(log.Fields{
				"slot":   slot,
				"action": name,
			}).Debug("do action ")
			epoch := common.SlotToEpoch(slot)
			start := common.EpochStart(epoch + int64(n))
			targetTime := common.TimeToSlot(start)
			total := targetTime - time.Now().Unix()
			log.WithFields(log.Fields{
				"slot":  slot,
				"total": total,
			}).Debug("delayToNextNEpochStart")
			time.Sleep(time.Second * time.Duration(total))
			r := plugins.PluginResponse{
				Cmd: types.CMD_NULL,
			}
			if len(params) > 0 {
				r.Result = params[0]
			}
			return r
		}, nil
	case "delayToNextNEpochEnd":
		n := 0
		if len(params) > 0 {
			n = params[0]
		}
		return func(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse {
			log.WithFields(log.Fields{
				"slot":   slot,
				"action": name,
			}).Debug("do action ")
			epoch := common.SlotToEpoch(slot)
			end := common.EpochEnd(epoch + int64(n))
			targetTime := common.TimeToSlot(end)
			total := targetTime - time.Now().Unix()
			log.WithFields(log.Fields{
				"slot":   slot,
				"target": end,
				"total":  total,
			}).Debug("delayToNextNEpochEnd")
			time.Sleep(time.Second * time.Duration(total))
			r := plugins.PluginResponse{
				Cmd: types.CMD_NULL,
			}
			if len(params) > 0 {
				r.Result = params[0]
			}
			return r
		}, nil
	case "delayToNextNEpochHalf":
		n := 1
		if len(params) > 0 {
			n = params[0]
		}
		slotsPerEpoch := backend.GetSlotsPerEpoch()
		return func(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse {
			log.WithFields(log.Fields{
				"slot":   slot,
				"action": name,
			}).Debug("do action ")
			epoch := common.SlotToEpoch(slot)
			start := common.EpochStart(epoch + int64(n))
			start += int64(slotsPerEpoch) / 2
			targetTime := common.TimeToSlot(start)
			total := targetTime - time.Now().Unix()
			log.WithFields(log.Fields{
				"slot":  slot,
				"total": total,
			}).Debug("delayToNextNEpochHalf")
			time.Sleep(time.Second * time.Duration(total))
			r := plugins.PluginResponse{
				Cmd: types.CMD_NULL,
			}
			if len(params) > 0 {
				r.Result = params[0]
			}
			return r
		}, nil
	case "delayToEpochEnd":
		return func(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse {
			log.WithFields(log.Fields{
				"slot":   slot,
				"action": name,
			}).Debug("do action ")

			epoch := common.SlotToEpoch(slot)
			end := common.EpochEnd(epoch)
			targetTime := common.TimeToSlot(end)
			total := targetTime - time.Now().Unix()
			log.WithFields(log.Fields{
				"slot":  slot,
				"total": total,
			}).Debug("delayToEpochEnd")
			time.Sleep(time.Second * time.Duration(total))
			r := plugins.PluginResponse{
				Cmd: types.CMD_NULL,
			}

			if len(params) > 0 {
				r.Result = params[0]
			}
			return r
		}, nil
	case "delayHalfEpoch":
		slotsPerEpoch := backend.GetSlotsPerEpoch()
		return func(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse {
			log.WithFields(log.Fields{
				"slot":   slot,
				"action": name,
			}).Debug("do action ")

			seconds := backend.GetIntervalPerSlot()
			if seconds == 0 {
				seconds = 12 // default 12 seconds
			}
			total := (seconds) * (slotsPerEpoch / 2)
			log.WithFields(log.Fields{
				"slot":  slot,
				"total": total,
			}).Debug("delayHalfEpoch")
			time.Sleep(time.Second * time.Duration(total))
			r := plugins.PluginResponse{
				Cmd: types.CMD_NULL,
			}
			if len(params) > 0 {
				r.Result = params[0]
			}
			return r
		}, nil
	case "packPooledAttest":
		return func(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse {
			r := plugins.PluginResponse{
				Cmd: types.CMD_NULL,
			}
			epoch := common.SlotToEpoch(slot)
			last := epoch - 1
			if last < 0 {
				last = 0
			}
			minSlot := common.EpochStart(last)
			maxSlot := common.EpochEnd(epoch)
			log.WithFields(log.Fields{
				"slot":   slot,
				"action": name,
			}).Debug("do action ")

			if len(params) == 0 {
				return r
			}
			block, ok := params[0].(*ethpb.SignedBeaconBlockDeneb)
			if !ok {
				log.WithFields(log.Fields{
					"param": fmt.Sprintf("%T", params[0]),
				}).Error("invalid param type, require *ethpb.SignedBeaconBlockDeneb")
				r.Result = params[0]
				return r
			}

			attackerAttestations := make([]*ethpb.Attestation, 0)
			pool := backend.GetAttestPool()
			for ns, atts := range pool {
				if int64(ns) < minSlot || int64(ns) > maxSlot {
					log.WithField("slot", ns).Debug("skip attestation at slot")
					continue
				}
				for _, att := range atts {
					attackerAttestations = append(attackerAttestations, att)
				}
			}
			backend.ResetAttestPool()

			allAtt := append(block.Block.Body.Attestations, attackerAttestations...)
			{
				// Remove duplicates from both aggregated/unaggregated attestations. This
				// prevents inefficient aggregates being created.
				atts, _ := types.ProposerAtts(allAtt).Dedup()
				attsByDataRoot := make(map[[32]byte][]*ethpb.Attestation, len(atts))
				for _, att := range atts {
					attDataRoot, err := att.Data.HashTreeRoot()
					if err != nil {
						continue
					}
					attsByDataRoot[attDataRoot] = append(attsByDataRoot[attDataRoot], att)
				}

				attsForInclusion := types.ProposerAtts(make([]*ethpb.Attestation, 0))
				for _, ass := range attsByDataRoot {
					as, err := attaggregation.Aggregate(ass)
					if err != nil {
						continue
					}
					attsForInclusion = append(attsForInclusion, as...)
				}
				deduped, _ := attsForInclusion.Dedup()
				sorted, err := deduped.SortByProfitability()
				if err != nil {
					log.WithError(err).Error("sort attestation failed")
				} else {
					atts = sorted.LimitToMaxAttestations()
				}
				allAtt = atts
			}

			block.Block.Body.Attestations = allAtt

			r.Result = block
			return r
		}, nil
	case "modifyAttestSource":
		if len(params) < 1 {
			log.WithField("action", action).Error("need at least 1 param.")
			return nil, errors.New("invalid param")
		}
		newSourceSlot := params[0]
		return func(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse {
			var attestation *ethpb.AttestationData
			r := plugins.PluginResponse{
				Cmd: types.CMD_NULL,
			}

			log.WithFields(log.Fields{
				"slot":   slot,
				"action": name,
			}).Debug("do action ")

			if len(params) > 0 {
				attestation = params[0].(*ethpb.AttestationData)
				if root, err := backend.GetSlotRoot(int64(newSourceSlot)); err == nil {
					attestation.Source.Root = common.FromHex(root)
					if r.Result, err = common.AttestationDataToBase64(attestation); err == nil {
						r.Cmd = types.CMD_UPDATE_STATE
					}
				}
			}

			return r
		}, nil
	case "modifyAttestTarget":
		if len(params) < 1 {
			log.WithField("action", action).Error("need at least 1 param.")
			return nil, errors.New("invalid param")
		}
		newSourceSlot := params[0]
		return func(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse {
			var attestation *ethpb.AttestationData
			r := plugins.PluginResponse{
				Cmd: types.CMD_NULL,
			}

			log.WithFields(log.Fields{
				"slot":   slot,
				"action": name,
			}).Debug("do action ")

			if len(params) > 0 {
				attestation = params[0].(*ethpb.AttestationData)
				if root, err := backend.GetSlotRoot(int64(newSourceSlot)); err == nil {
					attestation.Target.Root = common.FromHex(root)
					if r.Result, err = common.AttestationDataToBase64(attestation); err == nil {
						r.Cmd = types.CMD_UPDATE_STATE
					}
				}
			}

			return r
		}, nil
	case "modifyAttestHead":
		if len(params) < 1 {
			log.WithField("action", action).Error("need at least 1 param.")
			return nil, errors.New("invalid param")
		}
		newSourceSlot := params[0]
		return func(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse {
			var attestation *ethpb.AttestationData
			r := plugins.PluginResponse{
				Cmd: types.CMD_NULL,
			}

			log.WithFields(log.Fields{
				"slot":   slot,
				"action": name,
			}).Debug("do action ")

			if len(params) > 0 {
				attestation = params[0].(*ethpb.AttestationData)
				if root, err := backend.GetSlotRoot(int64(newSourceSlot)); err == nil {
					attestation.BeaconBlockRoot = common.FromHex(root)
					if r.Result, err = common.AttestationDataToBase64(attestation); err == nil {
						r.Cmd = types.CMD_UPDATE_STATE
					}
				}
			}
			return r
		}, nil

	case "modifyParentRoot":
		if len(params) < 1 {
			// error.
			log.WithField("action", action).Error("need at least 1 param.")
			return nil, errors.New("invalid param")
		}
		newSlot := params[0]

		return func(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse {
			r := plugins.PluginResponse{
				Cmd: types.CMD_NULL,
			}
			log.WithFields(log.Fields{
				"slot":   slot,
				"action": name,
			}).Debug("do action ")
			// get parent root by newSlot.
			newRoot, err := backend.GetSlotRoot(int64(newSlot))
			if err != nil {
				log.WithError(err).Error("get slot root failed")
				return r
			}

			r.Result = newRoot
			return r
		}, nil
	case "rePackAttestation":
		return func(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse {
			r := plugins.PluginResponse{
				Cmd: types.CMD_NULL,
			}
			log.WithFields(log.Fields{
				"slot":   slot,
				"action": name,
			}).Debug("do action ")

			if len(params) == 0 {
				return r
			}
			block := params[0].(*ethpb.SignedBeaconBlockDeneb)
			epoch := common.SlotToEpoch(slot)
			startEpoch := common.EpochStart(epoch)
			endEpoch := common.EpochEnd(epoch)
			attackerAttestations := make([]*ethpb.Attestation, 0)
			validatorSet := backend.GetValidatorDataSet()
			log.WithFields(log.Fields{
				"slot": slot,
			}).Debug("rePackAttestation")
			for i := startEpoch; i <= endEpoch; i++ {
				allSlotAttest := backend.GetAttestSet(uint64(i))
				if allSlotAttest == nil {
					continue
				}

				for publicKey, att := range allSlotAttest.Attestations {
					val := validatorSet.GetValidatorByPubkey(publicKey)
					if val == nil {
						log.WithField("pubkey", publicKey).Debug("validator not found")
						continue
					}
					valRole := backend.GetValidatorRole(int(i), int(val.Index))
					if val != nil && valRole == types.AttackerRole {
						log.WithField("pubkey", publicKey).Debug("add attacker attestation to block")
						attackerAttestations = append(attackerAttestations, att)
					}
					//log.WithField("pubkey", publicKey).Debug("add attacker attestation to block")
					//attackerAttestations = append(attackerAttestations, att)
				}
			}

			allAtt := append(block.Block.Body.Attestations, attackerAttestations...)
			{
				// Remove duplicates from both aggregated/unaggregated attestations. This
				// prevents inefficient aggregates being created.
				atts, _ := types.ProposerAtts(allAtt).Dedup()
				attsByDataRoot := make(map[[32]byte][]*ethpb.Attestation, len(atts))
				for _, att := range atts {
					attDataRoot, err := att.Data.HashTreeRoot()
					if err != nil {
						continue
					}
					attsByDataRoot[attDataRoot] = append(attsByDataRoot[attDataRoot], att)
				}

				attsForInclusion := types.ProposerAtts(make([]*ethpb.Attestation, 0))
				for _, ass := range attsByDataRoot {
					as, err := attaggregation.Aggregate(ass)
					if err != nil {
						continue
					}
					attsForInclusion = append(attsForInclusion, as...)
				}
				deduped, _ := attsForInclusion.Dedup()
				sorted, err := deduped.SortByProfitability()
				if err != nil {
					log.WithError(err).Error("sort attestation failed")
				} else {
					atts = sorted.LimitToMaxAttestations()
				}
				allAtt = atts
			}

			block.Block.Body.Attestations = allAtt

			r.Result = block
			return r
		}, nil
	default:
		log.WithField("name", name).Error("unknown function action name")
		return nil, errors.New(fmt.Sprintf("unknown function action name:%s", name))
	}
}
