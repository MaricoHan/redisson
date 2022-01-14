package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
)

var mysqldb *sql.DB

//Get return the db client
func Get() *sql.DB {
	return mysqldb
}

// Connect create a instance of the mysql db
func Connect(user, password, host, port, dbName string, options ...Option) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		user, password, host, port, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	var cfg = config{
		maxIdleConns: 10,
		maxOpenConns: 10,
		maxLifetime:  time.Hour,
		debug:        false,
		writer:       log.Log,
	}
	for _, optionFn := range options {
		if err := optionFn(&cfg); err != nil {
			panic(err)
		}
	}
	//SetMaxIdleConns 设置空闲连接池中连接的最大数量
	db.SetMaxIdleConns(cfg.maxIdleConns)

	//SetMaxOpenConns 设置打开数据库连接的最大数量。
	db.SetMaxOpenConns(cfg.maxOpenConns)

	//SetConnMaxLifetime 设置了连接可复用的最大时间。
	db.SetConnMaxLifetime(cfg.maxLifetime)

	mysqldb = db

	boil.DebugMode = cfg.debug
	boil.SetDB(db)
	boil.DebugWriter = cfg.writer
}

func Close() {
	_ = mysqldb.Close()
}
