package svc

import (
	"short/internal/config"
	"short/internal/sequence"
	"short/model"

	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config            config.Config
	ShortUrlModel     model.ReflectMapModel // 长短映射表
	Sequence          sequence.SeqInter     // 取号表
	ShortUrlBlackList map[string]struct{}   // 短域名黑名单
	Filter            *bloom.Filter         // bloomfilter做缓存穿透
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.ShortUrlDB.DSN)

	// 将配置中的短域名黑名单加载到map
	m := make(map[string]struct{}, len(c.ShortUrlBlackList))
	for _, v := range c.ShortUrlBlackList {
		m[v] = struct{}{}
	}

	// 初始化 bloomfilter 并加载数据
	store := redis.New(c.RedisConf.Host, func(r *redis.Redis) {
		r.Type = redis.NodeType
	})

	filter := bloom.New(store, "bloom_filter", 20*(1<<20))

	return &ServiceContext{
		Config:            c,
		ShortUrlModel:     model.NewReflectMapModel(conn, c.CacheRedis),
		Sequence:          sequence.NewSMysql(c.Sequence.DSN),
		ShortUrlBlackList: m,
		Filter:            filter,
	}
}
