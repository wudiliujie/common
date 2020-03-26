package mysql

import (
	"database/sql"
	"github.com/wudiliujie/common/db"
	"github.com/wudiliujie/common/log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	Context    *sql.DB
	PubContext *sql.DB
)

func Init(dbAddr string, dbPoolSize int32) {
	// dbAddr = "root:root@tcp(192.168.22.212:3306)/x_game_s1?charset=utf8&parseTime=True&loc=Local"
	log.Release("初始化Mysql数据库")
	var err error
	Context, err = sql.Open("mysql", dbAddr)

	if err != nil {
		log.Fatal("mysqldb init is error(%v)", err)
	}
	Context.SetMaxOpenConns(int(dbPoolSize))
	Context.SetMaxIdleConns(int(dbPoolSize) / 2)
	Context.SetConnMaxLifetime(time.Hour)
}
func InitPub(dbAddr string, dbPoolSize int32) {
	log.Release("初始化Mysql  Pub数据库")
	var err error
	PubContext, err = sql.Open("mysql", dbAddr)

	if err != nil {
		log.Fatal("mysqldb Pub数据库  init is error(%v)", err)
	}
	PubContext.SetMaxOpenConns(int(dbPoolSize))
	PubContext.SetMaxIdleConns(int(dbPoolSize) / 2)
	PubContext.SetConnMaxLifetime(time.Hour)
}

func Query(strsql string, args ...interface{}) ([]*db.DataRow, error) {
	return db.Query(Context, strsql, args...)
}
func QueryRow(strsql string, args ...interface{}) (*db.DataRow, error) {
	return db.QueryRow(Context, strsql, args...)
}
func PubQueryRow(strsql string, args ...interface{}) (*db.DataRow, error) {
	return db.QueryRow(PubContext, strsql, args...)
}
func PubQuery(strsql string, args ...interface{}) ([]*db.DataRow, error) {
	return db.Query(PubContext, strsql, args...)
}
func Exec(query string, args ...interface{}) (sql.Result, error) {
	//log.Debug("exec:%v", query)
	return Context.Exec(query, args...)
}
func PubExec(query string, args ...interface{}) (sql.Result, error) {
	//log.Debug("exec:%v", query)
	return PubContext.Exec(query, args...)
}
