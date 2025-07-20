package svc

import (
	"short/internal/config"
	"short/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	ShortUrlModel model.ReflectMapModel // table: reflectMap
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.ShortUrlDB.DSN)
	return &ServiceContext{
		Config:        c,
		ShortUrlModel: model.NewReflectMapModel(conn),
	}
}
