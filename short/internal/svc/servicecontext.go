package svc

import (
	"short/internal/config"
	"short/internal/sequence"
	"short/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	ShortUrlModel model.ReflectMapModel // 长短映射表
	Sequence      sequence.SeqInter     // 取号表
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.ShortUrlDB.DSN)
	return &ServiceContext{
		Config:        c,
		ShortUrlModel: model.NewReflectMapModel(conn),
		Sequence:      sequence.NewSMysql(c.Sequence.DSN),
	}
}
