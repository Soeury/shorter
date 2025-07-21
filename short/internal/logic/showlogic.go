package logic

import (
	"context"
	"database/sql"
	"errors"

	"short/internal/svc"
	"short/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

	// 1. 根据短链接查询长连接(采用go-zero生成带缓存的Mysql查询, 内嵌singleflight做请求合并)
	long, err := l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: req.ShortUrl, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found longUrl with your shortUrl")
		}

		logx.Errorw(
			"l.svcCtx.ShortUrlModel.FindOneBySurl failed",
			logx.LogField{Key: "err", Value: err.Error()},
		)
		return nil, err
	}

	// 2. 返回重定向响应(handler)
	return &types.ShowResponse{LongUrl: long.Lurl.String}, nil
}
