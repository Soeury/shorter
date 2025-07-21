package svc

import (
	"short/internal/config"
	"short/internal/sequence"
	"short/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config            config.Config
	ShortUrlModel     model.ReflectMapModel // 长短映射表
	Sequence          sequence.SeqInter     // 取号表
	ShortUrlBlackList map[string]struct{}   // 短域名黑名单
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.ShortUrlDB.DSN)

	// 将配置中的短域名黑名单加载到map
	m := make(map[string]struct{}, len(c.ShortUrlBlackList))
	for _, v := range c.ShortUrlBlackList {
		m[v] = struct{}{}
	}

	return &ServiceContext{
		Config:            c,
		ShortUrlModel:     model.NewReflectMapModel(conn),
		Sequence:          sequence.NewSMysql(c.Sequence.DSN),
		ShortUrlBlackList: m,
	}
}
