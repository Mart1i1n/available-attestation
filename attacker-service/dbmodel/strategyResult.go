package dbmodel

import (
	"github.com/astaxie/beego/orm"
	"github.com/tsinghua-cel/attacker-service/types"
)

type StrategyResult struct {
	ID                   int64  `orm:"column(id)" db:"id" json:"id" form:"id"`                                                                                 //  任务类型id
	UUID                 string `orm:"column(uuid)" db:"uuid" json:"uuid" form:"uuid"`                                                                         //  策略的唯一id
	ReorgCount           int    `orm:"column(rerog_count)" db:"rerog_count" json:"rerog_count" form:"rerog_count"`                                             // 重组次数
	ImpactValidatorCount int    `orm:"column(impact_validator_count)" db:"impact_validator_count" json:"impact_validator_count" form:"impact_validator_count"` // 影响验证者数量
}

func (StrategyResult) TableName() string {
	return "t_strategy_result"
}

type StrategyResultRepository interface {
	Create(reward *StrategyResult) error
	GetByUUID(uuid string) *StrategyResult
	GetListByFilter(filters ...interface{}) []*StrategyResult
}

type strategyResultRepositoryImpl struct {
	o orm.Ormer
}

func NewStrategyResultRepository(o orm.Ormer) StrategyResultRepository {
	return &strategyResultRepositoryImpl{o}
}

func (repo *strategyResultRepositoryImpl) Create(sr *StrategyResult) error {
	_, err := repo.o.Insert(sr)
	return err
}

func (repo *strategyResultRepositoryImpl) GetByUUID(uuid string) *StrategyResult {
	filters := make([]interface{}, 0)
	filters = append(filters, "uuid", uuid)
	return repo.GetListByFilter(filters...)[0]
}

func (repo *strategyResultRepositoryImpl) GetListByFilter(filters ...interface{}) []*StrategyResult {
	list := make([]*StrategyResult, 0)
	query := repo.o.QueryTable(new(StrategyResult).TableName())
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	query.OrderBy("-epoch").All(&list)
	return list
}

func InsertNewStrategyResult(st *types.Strategy, reorgCount, impactValidatorCount int) {
	sr := &StrategyResult{
		UUID:                 st.Uid,
		ReorgCount:           reorgCount,
		ImpactValidatorCount: impactValidatorCount,
	}
	NewStrategyResultRepository(orm.NewOrm()).Create(sr)
}

func GetStrategyResultByUUID(uuid string) *StrategyResult {
	return NewStrategyResultRepository(orm.NewOrm()).GetByUUID(uuid)
}
