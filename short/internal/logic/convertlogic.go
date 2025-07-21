package logic

import (
	"context"
	"database/sql"
	"time"

	erl "short/internal/errlogs"
	"short/internal/svc"
	"short/internal/types"
	"short/model"
	"short/pkg/base62"
	"short/pkg/connect"
	"short/pkg/md5"
	"short/pkg/urltool"

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
		return nil, erl.ErrInvalidLongUrl
	}

	// - 查长链(md5): 两张奇偶表
	md5Value := md5.Sum([]byte(req.LongUrl))
	_, err = l.svcCtx.ShortUrlModel.FindOneByMd5(l.ctx, sql.NullString{String: md5Value, Valid: true})
	if err != sql.ErrNoRows && err != nil {
		logx.Errorw(
			erl.LogFindByMd5Failed,
			logx.LogField{Key: "err", Value: err.Error()},
		)
		return nil, erl.ErrFindByMd5Failed
	}
	if err == nil {
		return nil, erl.ErrLongUrlExisted
	}

	_, err = l.svcCtx.ShortUrlModel2.FindOneByMd5(l.ctx, sql.NullString{String: md5Value, Valid: true})
	if err != sql.ErrNoRows && err != nil {
		logx.Errorw(
			erl.LogFindByMd5Failed,
			logx.LogField{Key: "err", Value: err.Error()},
		)
		return nil, erl.ErrFindByMd5Failed
	}
	if err == nil {
		return nil, erl.ErrLongUrlExisted
	}

	// - 拆分url获取短链
	baseUrl, err := urltool.GetBasePath(req.LongUrl)
	if err != nil {
		logx.Errorw(
			erl.LogURLToolGetFailed,
			logx.LogField{Key: "lurl", Value: req.LongUrl}, logx.LogField{Key: "err", Value: err.Error()},
		)

		return nil, erl.ErrURLToolGetFailed
	}

	// - 查短链: 两张表
	_, err = l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: baseUrl, Valid: true})
	if err != sqlx.ErrNotFound && err != nil {
		logx.Errorw(
			erl.LogFindBySurlFailed,
			logx.LogField{Key: "err", Value: err.Error()},
		)
		return nil, erl.ErrFindBySurlFailed
	}
	if err == nil {
		return nil, erl.ErrUseShortUrlToConvert
	}

	_, err = l.svcCtx.ShortUrlModel2.FindOneBySurl(l.ctx, sql.NullString{String: baseUrl, Valid: true})
	if err != sqlx.ErrNotFound && err != nil {
		logx.Errorw(
			erl.LogFindBySurlFailed,
			logx.LogField{Key: "err", Value: err.Error()},
		)
		return nil, erl.ErrFindBySurlFailed
	}
	if err == nil {
		return nil, erl.ErrUseShortUrlToConvert
	}

	var seq uint64
	var short string
	var shortUrl string

	for {
		// 2. 取号
		seq, err = l.svcCtx.Sequence.Next()
		if err != nil {
			logx.Errorw(
				erl.LogGetNextSeqFailed,
				logx.LogField{Key: "err", Value: err.Error()},
			)

			return nil, erl.ErrGetNextSeqFailed
		}

		// 3. 号码转短链
		// - 安全性: 打乱62进制字符
		// - 避免特殊字符: 短链黑名单
		short = base62.ChangeToBase62(seq)
		if _, ok := l.svcCtx.ShortUrlBlackList[short]; ok {
			logx.Errorw(
				erl.LogShortExistInBlackList,
				logx.LogField{Key: "err", Value: erl.ErrShortExistInBlackList.Error()},
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
				erl.LogInsertDBFailed,
				logx.LogField{Key: "err", Value: err.Error()},
			)
			return nil, erl.ErrInsertDBFailed
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
				erl.LogInsertDBFailed,
				logx.LogField{Key: "key", Value: err.Error()},
			)
			return nil, erl.ErrInsertDBFailed
		}
	}

	// 5. 短链存储到 filter
	if err := l.svcCtx.Filter.Add([]byte(short)); err != nil {
		logx.Errorw(
			erl.LogFilterAddFailed,
			logx.LogField{Key: "err", Value: err.Error()},
		)
		return nil, erl.ErrFilterAddFailed
	}

	shortUrl = l.svcCtx.Config.ShortDomain + "/" + short
	return &types.ConvertResponse{ShortUrl: shortUrl}, nil
}
