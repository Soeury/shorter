package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"short/internal/pkg/connect"
	"short/internal/pkg/md5"
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
	//    - 长链不为空
	//    - 长链有效
	//    - 长链未转过
	//    - 不能传入短链
	// 2. 取号
	// 3. 号码转短链
	// 4. 存短链

	if ok := connect.Get(req.LongUrl); !ok {
		return nil, errors.New("invalid longUrl")
	}

	md5Value := md5.Sum([]byte(req.LongUrl))
	_, err = l.svcCtx.ShortUrlModel.FindOneByMd5(l.ctx, sql.NullString{String: md5Value, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil { // 长链转存过
			return nil, fmt.Errorf("longurl has converted yet")
		}

		// 其他问题
		logx.Errorw("fineOneByMd5 failed,", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err

	}

	return nil, nil
}
