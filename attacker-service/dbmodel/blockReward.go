package dbmodel

import (
	"github.com/astaxie/beego/orm"
	"github.com/tsinghua-cel/attacker-service/common"
)

type BlockReward struct {
	ID                     int64 `orm:"column(id)" db:"id" json:"id" form:"id"`                                                                                         //  任务类型id
	Slot                   int64 `orm:"column(slot)" db:"slot" json:"slot" form:"slot"`                                                                                 // slot
	ProposerIndex          int   `orm:"column(proposer_index)" db:"proposer_index" json:"proposer_index" form:"proposer_index"`                                         // 验证者索引
	TotalAmount            int64 `orm:"column(total_amount)" db:"total_amount" json:"total_amount" form:"total_amount"`                                                 // Total 奖励数量
	AttestationAmount      int64 `orm:"column(attestation_amount)" db:"attestation_amount" json:"attestation_amount" form:"attestation_amount"`                         // Target 奖励数量
	SyncAggregateAmount    int64 `orm:"column(sync_aggregate_amount)" db:"sync_aggregate_amount" json:"sync_aggregate_amount" form:"sync_aggregate_amount"`             // Sync Aggregate 奖励数量
	ProposerSlashingAmount int64 `orm:"column(proposer_slashing_amount)" db:"proposer_slashing_amount" json:"proposer_slashing_amount" form:"proposer_slashing_amount"` // Proposer Slashing 奖励数量
	AttesterSlashingAmount int64 `orm:"column(attester_slashing_amount)" db:"attester_slashing_amount" json:"attester_slashing_amount" form:"attester_slashing_amount"` // Attester Slashing 奖励数量
}

func (BlockReward) TableName() string {
	return "t_block_reward"
}

type BlockRewardRepository interface {
	Create(reward *BlockReward) error
	GetListByFilter(filters ...interface{}) []*BlockReward
	GetListBySlotRange(start int64, end int64) []*BlockReward
}

type blockRewardRepositoryImpl struct {
	o orm.Ormer
}

func NewBlockRewardRepository(o orm.Ormer) BlockRewardRepository {
	return &blockRewardRepositoryImpl{o}
}

func (repo *blockRewardRepositoryImpl) Create(reward *BlockReward) error {
	_, err := repo.o.Insert(reward)
	return err
}

func (repo *blockRewardRepositoryImpl) GetListByFilter(filters ...interface{}) []*BlockReward {
	list := make([]*BlockReward, 0)
	query := repo.o.QueryTable(new(BlockReward).TableName())
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	query.OrderBy("-slot").All(&list)
	return list
}

func (repo *blockRewardRepositoryImpl) GetListBySlotRange(start int64, end int64) []*BlockReward {
	list := make([]*BlockReward, 0)
	query := repo.o.QueryTable(new(BlockReward).TableName())
	query = query.Filter("slot__gte", start)
	query = query.Filter("slot__lte", end)
	query.OrderBy("-slot").All(&list)

	return list
}

func InsertBlockReward(o orm.Ormer, reward *BlockReward) error {
	return NewBlockRewardRepository(o).Create(reward)
}

func GetBlockRewardListByEpoch(epoch int64) []*BlockReward {
	start := common.EpochStart(epoch)
	end := common.EpochEnd(epoch)
	return NewBlockRewardRepository(orm.NewOrm()).GetListBySlotRange(start, end)
}
