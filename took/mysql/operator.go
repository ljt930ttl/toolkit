package mysql

import (
	"context"
	"database/sql"
	"sync"

	"took/logger"
)

type Operator interface {
	QueryRows(sqlStr string, args ...interface{}) ([]MapSS, error)
	ExecControl(sqlStr string, args ...interface{}) error
	GetDB() *ConnDB
	Close() error
}

type BaseDB struct {
	sql.DB
	isDebug bool
}

type ConnDB struct {
	sync.Mutex
	Conn    *sql.Conn
	CTX     context.Context
	isDebug bool
}

// QueryRows 查询所有
func (db *BaseDB) QueryRows(sqlStr string, args ...interface{}) ([]MapSS, error) {
	rows, err := db.Query(sqlStr, args...)
	if err != nil {
		logger.Error("err:%s\nsql:%s\n", err.Error(), sqlStr)
		return nil, err
	}
	return ProduceResultForRows(rows)
}

func (db *BaseDB) Close() error {
	return db.Close()
}

func (db *BaseDB) ExecControl(sqlStr string, args ...interface{}) error {
	res, err := db.Exec(sqlStr, args...)
	if err != nil {
		logger.Error("err:%s\nsql:%s,%s \n", err.Error(), sqlStr, args)
		return err
	}
	if db.isDebug == true {
		num, err := res.RowsAffected()
		if err != nil {
			return err
		}
		logger.Debug("Total of %d rows are affected!\n", num)
	}
	return nil
}

func (db *ConnDB) QueryRows(sqlStr string, args ...interface{}) ([]MapSS, error) {
	rows, err := db.Conn.QueryContext(db.CTX, sqlStr, args...)
	if err != nil {
		logger.Error("err:%s\nsql:%s\n", err.Error(), sqlStr)
		return nil, err
	}
	return ProduceResultForRows(rows)
}

func (db *ConnDB) ExecControl(sqlStr string, args ...interface{}) error {
	res, err := db.Conn.ExecContext(db.CTX, sqlStr, args...)
	if err != nil {
		logger.Error("err:%s\nsql:%s,%s \n", err.Error(), sqlStr, args)
		return err
	}
	if db.isDebug == true {
		num, err := res.RowsAffected()
		if err != nil {
			return err
		}
		logger.Debug("Total of %d rows are affected!\n", num)
	}
	return nil
}

func (db *ConnDB) GetDB() *ConnDB {
	return db
}

func (db *ConnDB) Close() error {
	return db.Conn.Close()
}

func ProduceResultForRows(rows *sql.Rows) ([]MapSS, error) {
	//函数结束释放链接
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}(rows)
	//读出查询出的列字段名
	cols, err := rows.Columns()
	if len(cols) == 0 {
		return nil, err
	}
	//values是每个列的值，这里获取到byte里
	values := make([][]byte, len(cols))
	//query.Scan的参数，因为每次查询出来的列是不定长的，用len(cols)定住当次查询的长度
	scans := make([]interface{}, len(cols))
	//让每一行数据都填充到[][]byte里面,狸猫换太子
	for i := range values {
		scans[i] = &values[i]
	}
	results := make([]MapSS, 0, 10)
	for rows.Next() {
		err := rows.Scan(scans...)
		if err != nil {
			return nil, err
		}
		row := make(MapSS, 10)
		for k, v := range values { //每行数据是放在values里面，现在把它挪到row里
			key := cols[k]
			row[key] = string(v)
		}
		results = append(results, row)
	}
	//返回数据
	//fmt.Println(results)
	return results, nil
}

/*
//Example
// A *DB is a pool of connections. Call Conn to reserve a connection for
//exclusive use.

conn, err := db.Conn(ctx)
if err != nil {     log.Fatal(err) }
defer conn.Close()
//Return the connection to the pool.
id := 41
result, err := conn.ExecContext(ctx, `UPDATE balances SET balance = balance + 10 WHERE user_id = ?;`, id)
if err != nil {
	log.Fatal(err)
}
rows, err := result.RowsAffected()
if err != nil {
	log.Fatal(err)
}
if rows != 1 {
	log.Fatalf("expected single row affected, got %d rows affected", rows)
}

*/
