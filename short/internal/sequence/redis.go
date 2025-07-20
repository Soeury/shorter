package sequence

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// 基于redis实现的取号器
var SeqKey = "Shortener:Sequence:SeqKey"

type SRedis struct {
	conn *redis.Redis
}

func NewRedis(host string, ty string, pass string, tls bool) *SRedis {
	conf := redis.RedisConf{
		Host: host,
		Type: ty,
		Pass: pass,
		Tls:  tls,
	}

	return &SRedis{
		conn: redis.MustNewRedis(conf),
	}
}

// Next 取出下一个号
func (rds *SRedis) Next() (ret uint64, err error) {

	conn := rds.conn
	number, err := conn.IncrCtx(context.Background(), SeqKey)
	if err != nil {
		logx.Errorw(
			"sequence: conn.IncrCtx failed",
			logx.LogField{Key: "err", Value: err.Error()},
		)

		return 0, err
	}

	return uint64(number), nil
}
