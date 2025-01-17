package dbmodel

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type AttestReward struct {
	ID             int64 `orm:"column(id)" db:"id" json:"id" form:"id"`                                                     //  任务类型id
	Epoch          int64 `orm:"column(epoch)" db:"epoch" json:"epoch" form:"epoch"`                                         // epoch
	ValidatorIndex int   `orm:"column(validator_index)" db:"validator_index" json:"validator_index" form:"validator_index"` // 验证者索引
	HeadAmount     int64 `orm:"column(head_amount)" db:"head_amount" json:"head_amount" form:"head_amount"`                 // Head 奖励数量
	TargetAmount   int64 `orm:"column(target_amount)" db:"target_amount" json:"target_amount" form:"target_amount"`         // Target 奖励数量
	SourceAmount   int64 `orm:"column(source_amount)" db:"source_amount" json:"source_amount" form:"source_amount"`         // Source 奖励数量
	//Head	Target	Source	Inclusion Delay	Inactivity
}

func (AttestReward) TableName() string {
	return "t_attest_reward"
}

type AttestRewardRepository interface {
	Create(reward *AttestReward) error
	GetListByFilter(filters ...interface{}) []*AttestReward
}

type attestRewardRepositoryImpl struct {
	o orm.Ormer
}

func NewAttestRewardRepository(o orm.Ormer) AttestRewardRepository {
	return &attestRewardRepositoryImpl{o}
}

func (repo *attestRewardRepositoryImpl) Create(reward *AttestReward) error {
	_, err := repo.o.Insert(reward)
	return err
}

func (repo *attestRewardRepositoryImpl) GetListByFilter(filters ...interface{}) []*AttestReward {
	list := make([]*AttestReward, 0)
	query := repo.o.QueryTable(new(AttestReward).TableName())
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	query.OrderBy("-epoch").All(&list)
	return list
}

func GetRewardListByEpoch(epoch int64) []*AttestReward {
	filters := make([]interface{}, 0)
	filters = append(filters, "epoch", epoch)
	return NewAttestRewardRepository(orm.NewOrm()).GetListByFilter(filters...)
}

func GetRewardListByValidatorIndex(index int) []*AttestReward {
	filters := make([]interface{}, 0)
	filters = append(filters, "validator_index", index)
	return NewAttestRewardRepository(orm.NewOrm()).GetListByFilter(filters...)
}

func GetRewardByValidatorAndEpoch(epoch int64, index int) *AttestReward {
	filters := make([]interface{}, 0)
	filters = append(filters, "epoch", epoch)
	filters = append(filters, "validator_index", index)

	list := NewAttestRewardRepository(orm.NewOrm()).GetListByFilter(filters...)
	if len(list) >= 0 {
		return list[0]
	}
	return nil
}

func GetMaxEpoch() int64 {
	var max int64
	sql := fmt.Sprintf("select max(epoch) from %s", new(AttestReward).TableName())
	if err := orm.NewOrm().Raw(sql).QueryRow(&max); err == orm.ErrNoRows {
		return -1
	}
	return max
}

func GetImpactValidatorCount(maxHackValIdx int, normalTargetAmount int64, epoch int64) int {
	// impact normal validator count
	var countNormal int
	sql := fmt.Sprintf("select count(1) from %s where epoch = ? and target_amount < ? and validator_index > ?", new(AttestReward).TableName())
	orm.NewOrm().Raw(sql, epoch, normalTargetAmount, maxHackValIdx).QueryRow(&countNormal)

	var countHacked int
	sql = fmt.Sprintf("select count(1) from %s where epoch = ? and target_amount >= ? and validator_index <= ?", new(AttestReward).TableName())
	orm.NewOrm().Raw(sql, epoch, normalTargetAmount, maxHackValIdx).QueryRow(&countHacked)
	return countNormal + countHacked

}
