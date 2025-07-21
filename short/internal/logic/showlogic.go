package logic

import (
	"context"
	"database/sql"
	"errors"

	"short/internal/svc"
	"short/internal/types"
	"short/pkg/base62"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Show 根据短链获得长链进行重定向
func (l *ShowLogic) Show(req *types.ShowRequest) (resp *types.ShowResponse, err error) {

	// 1. 过滤器(基于内存or基于redis)
	exists, err := l.svcCtx.Filter.Exists([]byte(req.ShortUrl))
	if err != nil {
		logx.Errorw(
			"l.svcCtx.Filter.Exists failed",
			logx.LogField{Key: "err", Value: err.Error()},
		)

		return nil, err
	}

	if !exists {
		return nil, errors.New("shortUrl not exists in bloomFilter")
	}

	// 2. 根据短链接查询长连接(采用go-zero生成带缓存的Mysql查询, 内嵌singleflight做请求合并)
	// 短链转成 10 进制，判断序号奇偶性，然后去查库
	seq := base62.ChangeToBase10(req.ShortUrl)
	if seq%2 == 1 {
		long, err := l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: req.ShortUrl, Valid: true})

		if err != nil {
			if err == sqlx.ErrNotFound {
				return nil, err
			}

			logx.Errorw(
				"l.svcCtx.ShortUrlModel.FindOneBySurl",
				logx.LogField{Key: "err", Value: err.Error()},
			)
			return nil, err
		}
	}

	// long, err := l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: req.ShortUrl, Valid: true})
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return nil, errors.New("not found longUrl with your shortUrl")
	// 	}

	// 	logx.Errorw(
	// 		"l.svcCtx.ShortUrlModel.FindOneBySurl failed",
	// 		logx.LogField{Key: "err", Value: err.Error()},
	// 	)
	// 	return nil, err
	// }

	return nil, nil
}
