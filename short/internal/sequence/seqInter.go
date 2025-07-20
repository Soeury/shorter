package sequence

// 复用接口
type SeqInter interface {
	Next() (ret uint64, err error)
}
