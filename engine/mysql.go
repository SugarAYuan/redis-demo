package engine

import (
	"redisdemo/tools"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var MysqlDB *gorm.DB

func InitMysql() {
	MysqlDB = getConn()
}

func getConn() *gorm.DB {

	dsn := tools.Config.GetString("mysql.dsn")

	db, err := gorm.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}

	// 初始化连接池
	db.DB().SetMaxIdleConns(tools.Config.GetInt("mysql.idle_conns"))
	db.DB().SetMaxOpenConns(tools.Config.GetInt("mysql.max_conns"))
	idleTimeout := tools.Config.GetInt("mysql.idle_timeout")
	if idleTimeout > 0 {
		db.DB().SetConnMaxLifetime(time.Duration(idleTimeout) * time.Second)
	}

	debugFlag := tools.Config.GetBool("mysql.open_orm_stdout")
	if debugFlag {
		db.LogMode(true)
	} else {
		db.LogMode(false)
	}
	return db
}
