package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"short/internal/pkg/base62"
	"short/internal/pkg/connect"
	"short/internal/pkg/md5"
	"short/internal/pkg/urltool"
	"short/internal/svc"
	"short/internal/types"
	"short/model"

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

	// - 查长链(md5): 两张奇偶表
	md5Value := md5.Sum([]byte(req.LongUrl))
	_, err = l.svcCtx.ShortUrlModel.FindOneByMd5(l.ctx, sql.NullString{String: md5Value, Valid: true})
	if err != sql.ErrNoRows && err != nil {
		logx.Errorw(
			"l.svcCtx.ShortUrlModel.FindOneByMd5 failed",
			logx.LogField{Key: "err", Value: err.Error()},
		)
		return nil, err
	}
	if err == nil {
		return nil, errors.New("longUrl already existed in reflect_map")
	}
	_, err = l.svcCtx.ShortUrlModel2.FindOneByMd5(l.ctx, sql.NullString{String: md5Value, Valid: true})
	if err != sql.ErrNoRows && err != nil {
		logx.Errorw(
			"l.svcCtx.ShortUrlModel2.FindOneByMd5 failed",
			logx.LogField{Key: "err", Value: err.Error()},
		)
		return nil, err
	}
	if err == nil {
		return nil, errors.New("longUrl already existed in reflect_map2")
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

	// - 查短链: 两张表
	_, err = l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: baseUrl, Valid: true})
	if err != sqlx.ErrNotFound && err != nil {
		logx.Errorw(
			"ShortUrlModel.FindOneBySurl failed,",
			logx.LogField{Key: "err", Value: err.Error()},
		)
		return nil, err
	}
	if err == nil {
		return nil, fmt.Errorf("cannot use shortUrl convert to shortUrl")
	}
	_, err = l.svcCtx.ShortUrlModel2.FindOneBySurl(l.ctx, sql.NullString{String: baseUrl, Valid: true})
	if err != sqlx.ErrNotFound && err != nil {
		logx.Errorw(
			"ShortUrlModel2.FindOneBySurl failed,",
			logx.LogField{Key: "err", Value: err.Error()},
		)
		return nil, err
	}
	if err == nil {
		return nil, fmt.Errorf("cannot use shortUrl convert to shortUrl")
	}

	var seq uint64
	var short string
	var shortUrl string

	for {
		// 2. 取号
		seq, err = l.svcCtx.Sequence.Next()
		if err != nil {
			logx.Errorw(
				"l.svcCtx.Sequence.Next failed,",
				logx.LogField{Key: "err", Value: err.Error()},
			)

			return nil, err
		}

		// 3. 号码转短链
		// - 安全性: 打乱62进制字符
		// - 避免特殊字符: 短链黑名单
		short = base62.ChangeToBase62(seq)
		if _, ok := l.svcCtx.ShortUrlBlackList[short]; ok {
			logx.Errorw(
				"short existed in balck list",
			)
		} else if !ok {
			break
		}
	}

	// 4. 存储映射关系表, 根据取出号的奇偶性选择表
	if seq%2 == 1 {
		if _, err := l.svcCtx.ShortUrlModel.Insert(
			l.ctx,
			&model.ReflectMap{
				Lurl: sql.NullString{String: req.LongUrl, Valid: true},
				Md5:  sql.NullString{String: md5Value, Valid: true},
				Surl: sql.NullString{String: short, Valid: true},
				ExpireAt: sql.NullTime{
					Time: time.Date(
						time.Now().Year(),
						time.Now().Month()+6,
						time.Now().Day(),
						time.Now().Hour(),
						time.Now().Minute(),
						time.Now().Second(),
						time.Now().Nanosecond(),
						time.UTC,
					),
					Valid: true,
				},
			},
		); err != nil {
			logx.Errorw(
				"l.svcCtx.ShortUrlModel.Insert failed",
				logx.LogField{Key: "err", Value: err.Error()},
			)
			return nil, err
		}
	} else {
		if _, err := l.svcCtx.ShortUrlModel2.Insert(
			l.ctx,
			&model.ReflectMap2{
				Lurl: sql.NullString{String: req.LongUrl, Valid: true},
				Md5:  sql.NullString{String: md5Value, Valid: true},
				Surl: sql.NullString{String: short, Valid: true},
				ExpireAt: sql.NullTime{
					Time: time.Date(
						time.Now().Year(),
						time.Now().Month()+6,
						time.Now().Day(),
						time.Now().Hour(),
						time.Now().Minute(),
						time.Now().Second(),
						time.Now().Nanosecond(),
						time.UTC,
					),
					Valid: true,
				},
			},
		); err != nil {
			logx.Errorw(
				"l.svcCtx.ShortUrlModel2.Insert failed",
				logx.LogField{Key: "key", Value: err.Error()},
			)
			return nil, err
		}
	}

	// 5. 短链存储到 filter
	if err := l.svcCtx.Filter.Add([]byte(short)); err != nil {
		logx.Errorw(
			"l.svcCtx.Filter.Add failed",
			logx.LogField{Key: "err", Value: err.Error()},
		)
		return nil, err
	}

	shortUrl = l.svcCtx.Config.ShortDomain + "/" + short
	return &types.ConvertResponse{ShortUrl: shortUrl}, nil
}
