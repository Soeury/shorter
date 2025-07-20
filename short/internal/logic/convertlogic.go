package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"short/internal/pkg/connect"
	"short/internal/pkg/md5"
	"short/internal/pkg/urltool"
	"short/internal/svc"
	"short/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Convert 长链转短链
func (l *ConvertLogic) Convert(req *types.ConvertRequest) (resp *types.ConvertResponse, err error) {
	// 1. 取长链，校验数据(validate)
	// - 长链不为空且有效
	if ok := connect.Get(req.LongUrl); !ok {
		return nil, errors.New("invalid longUrl")
	}

	// - 长链未转存过
	md5Value := md5.Sum([]byte(req.LongUrl))
	_, err = l.svcCtx.ShortUrlModel.FindOneByMd5(l.ctx, sql.NullString{String: md5Value, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {
			return nil, fmt.Errorf("longurl has converted yet")
		}

		logx.Errorw(
			"ShortUrlModel.FindOneByMd5 failed",
			logx.LogField{Key: "err", Value: err.Error()},
		)

		return nil, err
	}

	// - 拆分url获取短链
	baseUrl, err := urltool.GetBasePath(req.LongUrl)
	if err != nil {
		logx.Errorw(
			"urltool.GetBasePath failed",
			logx.LogField{Key: "lurl", Value: req.LongUrl}, logx.LogField{Key: "err", Value: err.Error()},
		)

		return nil, err
	}

	// - 不能传入短链
	_, err = l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: baseUrl, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {
			return nil, fmt.Errorf("use shortUrl convert to shortUrl")
		}

		logx.Errorw(
			"ShortUrlModel.FindOneBySurl failed,",
			logx.LogField{Key: "err", Value: err.Error()},
		)

		return nil, err
	}

	return nil, nil
}
