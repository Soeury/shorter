package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ReflectMapModel = (*customReflectMapModel)(nil)

type (
	// ReflectMapModel is an interface to be customized, add more methods here,
	// and implement the added methods in customReflectMapModel.
	ReflectMapModel interface {
		reflectMapModel
		withSession(session sqlx.Session) ReflectMapModel
	}

	customReflectMapModel struct {
		*defaultReflectMapModel
	}
)

// NewReflectMapModel returns a model for the database table.
func NewReflectMapModel(conn sqlx.SqlConn) ReflectMapModel {
	return &customReflectMapModel{
		defaultReflectMapModel: newReflectMapModel(conn),
	}
}

func (m *customReflectMapModel) withSession(session sqlx.Session) ReflectMapModel {
	return NewReflectMapModel(sqlx.NewSqlConnFromSession(session))
}
