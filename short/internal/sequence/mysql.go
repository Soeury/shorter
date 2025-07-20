package sequence

import (
	"database/sql"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 基于 mysql 实现的取号器
const sqlReplaceIntoStub = `replace into sequence (stub) values ('a')`

type SMysql struct {
	conn sqlx.SqlConn
}

func NewSMysql(dsn string) *SMysql {
	return &SMysql{
		conn: sqlx.NewMysql(dsn),
	}
}

// Next 取出下一个号
func (m *SMysql) Next() (ret uint64, err error) {

	// prepare 预编译
	var stmt sqlx.StmtSession
	stmt, err = m.conn.Prepare(sqlReplaceIntoStub)
	if err != nil {
		logx.Errorw(
			"sequence: m.conn.Prepare failed",
			logx.LogField{Key: "err", Value: err.Error()},
		)

		return 0, err
	}
	defer stmt.Close()

	// exec 执行
	var rest sql.Result
	rest, err = stmt.Exec()
	if err != nil {
		logx.Errorw(
			"sequence: stmt.Exec failed",
			logx.LogField{Key: "err", Value: err.Error()},
		)

		return 0, err
	}

	// get last id 获取上一次插入的 id
	var lid int64
	lid, err = rest.LastInsertId()
	if err != nil {
		logx.Errorw(
			"sequence: rest.LastInsertId failed",
			logx.LogField{Key: "err", Value: err.Error()},
		)

		return 0, err
	}

	return uint64(lid), nil
}
