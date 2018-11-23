package models

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type MySqlConfig struct {
	UserName         string
	Password         string
	Port             string
	Host             string
	DBName           string
	MaxIdleConns     int
	MaxOpenConn      int
	ShowSqlOnConsole bool
}

var Engine *xorm.Engine

func (db MySqlConfig) MySqlConn() {
	//var err error
	log.Println("sql--->", db.UserName+":'"+db.Password+"'@tcp("+db.Host+":"+db.Port+")/"+db.DBName+"?charset=utf8")
	E, err := xorm.NewEngine("mysql", db.UserName+":"+db.Password+"@tcp("+db.Host+":"+db.Port+")/"+db.DBName+"?charset=utf8")
	//b, err := sql.Open("mysql", db.UserName+":"+db.Password+"@tcp("+db.Host+":"+db.Port+")/"+db.DBName+"?charset=utf8")
	if err != nil {
		log.Fatal("mySql Conn error:", err)

	}
	//defer Engine.Close()

	E.SetMaxIdleConns(db.MaxIdleConns)
	E.SetMaxOpenConns(db.MaxOpenConn)
	E.ShowSQL(db.ShowSqlOnConsole)

	Engine = E
}

func MysqlPing() {
	if err := Engine.Ping(); err != nil {
		log.Fatal("Mysql连接测试异常：", err)
	}
	log.Println("Mysqll连接测试成功")
}
