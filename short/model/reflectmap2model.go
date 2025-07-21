package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ReflectMap2Model = (*customReflectMap2Model)(nil)

type (
	// ReflectMap2Model is an interface to be customized, add more methods here,
	// and implement the added methods in customReflectMap2Model.
	ReflectMap2Model interface {
		reflectMap2Model
	}

	customReflectMap2Model struct {
		*defaultReflectMap2Model
	}
)

// NewReflectMap2Model returns a model for the database table.
func NewReflectMap2Model(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ReflectMap2Model {
	return &customReflectMap2Model{
		defaultReflectMap2Model: newReflectMap2Model(conn, c, opts...),
	}
}
