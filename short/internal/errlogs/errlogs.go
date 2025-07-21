package errlogs

import "errors"

var (
	ErrInvalidLongUrl             = errors.New("invalid long url in req")
	ErrFindByMd5Failed            = errors.New("l.svcCtx.ShortUrlModel.FindOneByMd5 failed")
	ErrFindBySurlFailed           = errors.New("l.svcCtx.ShortUrlModel.FindOneBySurl failed")
	ErrLongUrlExisted             = errors.New("long url is already existed")
	ErrURLToolGetFailed           = errors.New("url tool get path failed")
	ErrUseShortUrlToConvert       = errors.New("cannot use shortUrl convert to shortUrl")
	ErrGetNextSeqFailed           = errors.New("l.svcCtx.Sequence.Next failed")
	ErrShortExistInBlackList      = errors.New("short existed in balck list")
	ErrInsertDBFailed             = errors.New("l.svcCtx.ShortUrlModel.Insert failed")
	ErrFilterAddFailed            = errors.New("l.svcCtx.Filter.Add failed")
	ErrFilterNotExisted           = errors.New("filter not existed")
	ErrShortUrlNotExistedInFilter = errors.New("short url not existed in filter")
	ErrShortUrlHasBeenDel         = errors.New("shortUrl has been deleted")
	ErrUpdateDBFailed             = errors.New("l.svcCtx.ShortUrlModel.Update failed")
)

var (
	LogInvalidLongUrl        = "invalid long url in req"
	LogFindByMd5Failed       = "l.svcCtx.ShortUrlModel.FindOneByMd5 failed"
	LogFindBySurlFailed      = "l.svcCtx.ShortUrlModel.FindOneBySurl failed"
	LogLongUrlExisted        = "long url is already existed"
	LogURLToolGetFailed      = "url tool get path failed"
	LogGetNextSeqFailed      = "l.svcCtx.Sequence.Next failed"
	LogShortExistInBlackList = "short existed in balck list"
	LogInsertDBFailed        = "l.svcCtx.ShortUrlModel.Insert failed"
	LogFilterAddFailed       = "l.svcCtx.Filter.Add failed"
	LogFilterNotExisted      = "filter not existed"
	LogShortUrlNotExisted    = "short url not existed in filter"
	LogShortUrlHasBeenDel    = "shortUrl has been deleted"
	LogUpdateDBFailed        = "l.svcCtx.ShortUrlModel.Update failed"
)
