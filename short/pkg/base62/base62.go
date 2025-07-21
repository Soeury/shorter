package base62

import (
	"math"
	"strings"
)

var (
	base62Str string
)

func MustInit(bs string) {
	if len(bs) == 0 {
		panic("must need a base string")
	}

	base62Str = bs
}

// ChangeToBase62  uint转换为62进制字符串
func ChangeToBase62(seq uint64) string {

	if seq == 0 {
		return string(base62Str[0])
	}

	b := []byte{}
	for seq > 0 {
		mod := seq % 62
		div := seq / 62
		b = append(b, base62Str[mod])
		seq = div
	}

	return string(Reverse(b))
}

// ChangeToBase10   62进制字符串转换为10进制uint
func ChangeToBase10(str62 string) uint64 {

	seq := uint64(0)
	b := []byte(str62)
	b = Reverse(b)

	for idx, v := range b {
		base := math.Pow(62, float64(idx))
		seq += uint64(strings.Index(base62Str, string(v))) * uint64(base)
	}

	return seq
}

func Reverse(b []byte) []byte {

	for i, j := 0, len(b)-1; i < len(b)/2; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	return b
}
