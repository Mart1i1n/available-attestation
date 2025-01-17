package dbmodel

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/types"
)

type Strategy struct {
	ID                   int64   `orm:"column(id)" db:"id" json:"id" form:"id"`                                                                                 //  任务类型id
	UUID                 string  `orm:"column(uuid)" db:"uuid" json:"uuid" form:"uuid"`                                                                         //  策略的唯一id
	Content              string  `orm:"column(content);size(3000)" db:"content" json:"content" form:"content"`                                                  //  策略内容
	MinEpoch             int64   `orm:"column(min_epoch)" db:"min_epoch" json:"min_epoch" form:"min_epoch"`                                                     //  最小epoch
	MaxEpoch             int64   `orm:"column(max_epoch)" db:"max_epoch" json:"max_epoch" form:"max_epoch"`                                                     //  最大epoch
	IsEnd                bool    `orm:"column(is_end)" db:"is_end" json:"is_end" form:"is_end"`                                                                 // 是否结束
	ReorgCount           int     `orm:"column(rerog_count)" db:"rerog_count" json:"rerog_count" form:"rerog_count"`                                             // 重组次数
	ImpactValidatorCount int     `orm:"column(impact_validator_count)" db:"impact_validator_count" json:"impact_validator_count" form:"impact_validator_count"` // 影响验证者数量
	HonestLoseRateAvg    float64 `orm:"column(honest_lose_rate_avg)" db:"honest_lose_rate_avg" json:"honest_lose_rate_avg" form:"honest_lose_rate_avg"`         // 诚实验证者平均损失率
	AttackerLoseRateAvg  float64 `orm:"column(attacker_lose_rate_avg)" db:"attacker_lose_rate_avg" json:"attacker_lose_rate_avg" form:"attacker_lose_rate_avg"` // 攻击者平均损失率
}

func (Strategy) TableName() string {
	return "t_strategy"
}

type StrategyRepository interface {
	Create(st *Strategy) error
	Update(st *Strategy) error
	GetByUUID(uuid string) *Strategy
	GetListByFilter(filters ...interface{}) []*Strategy
}

type strategyRepositoryImpl struct {
	o orm.Ormer
}

func NewStrategyRepository(o orm.Ormer) StrategyRepository {
	return &strategyRepositoryImpl{o}
}

func (repo *strategyRepositoryImpl) Create(reward *Strategy) error {
	_, err := repo.o.Insert(reward)
	return err
}

func (repo *strategyRepositoryImpl) Update(st *Strategy) error {
	_, err := repo.o.Update(st)
	return err
}

func (repo *strategyRepositoryImpl) HasByUUID(uuid string) bool {
	filters := make([]interface{}, 0)
	filters = append(filters, "uuid", uuid)
	return len(repo.GetListByFilter(filters...)) > 0
}

func (repo *strategyRepositoryImpl) GetByUUID(uuid string) *Strategy {
	filters := make([]interface{}, 0)
	filters = append(filters, "uuid", uuid)
	list := repo.GetListByFilter(filters...)
	if len(list) > 0 {
		return list[0]
	} else {
		return nil
	}
}

func (repo *strategyRepositoryImpl) GetListByFilter(filters ...interface{}) []*Strategy {
	list := make([]*Strategy, 0)
	query := repo.o.QueryTable(new(Strategy).TableName())
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	query.OrderBy("-id").All(&list)
	return list
}

func InsertNewStrategy(st *types.Strategy) {
	d, _ := json.Marshal(st)
	data := &Strategy{
		UUID:                 st.Uid,
		Content:              string(d),
		IsEnd:                false,
		ReorgCount:           0,
		ImpactValidatorCount: 0,
	}
	if err := NewStrategyRepository(orm.NewOrm()).Create(data); err != nil {
		log.WithError(err).Error("failed to insert new strategy")
	}
}

func GetStrategyByUUID(uuid string) *Strategy {
	return NewStrategyRepository(orm.NewOrm()).GetByUUID(uuid)
}

func StrategyUpdate(st *Strategy) {
	NewStrategyRepository(orm.NewOrm()).Update(st)
}
