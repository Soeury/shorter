package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	ShortUrlDB struct { // reflect_map
		DSN string
	}

	ShortUrlDB2 struct { // reflect_map2
		DSN string
	}

	Sequence struct { // sequence
		DSN string
	}

	RedisConf struct {
		Host string
		Type string
		Pass string
		Tls  bool
	} // filter

	Base62String string

	ShortUrlBlackList []string

	ShortDomain string

	CacheRedis cache.CacheConf // sequence:redis
}
