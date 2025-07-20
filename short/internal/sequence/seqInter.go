package sequence

type SeqInter interface {
	Next() (ret uint64, err error)
}
