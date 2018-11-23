package main

import (
	"flag"
	"time"

	"github.com/jetlwx/zabbixAlert/models"
)

var dbuser, dbname, dbpass, dbhost, dbport string
var level0, level1, level2, level3, level4, level5 int64

func init() {
	flag.StringVar(&dbuser, "dbuser", "zabbix", "the user to connect db")
	flag.StringVar(&dbname, "dbname", "zabbix", "the dbname of db")
	flag.StringVar(&dbpass, "dbpass", "zabbix", "the password to connect db")
	flag.StringVar(&dbhost, "dbhost", "192.168.1.229", "the host ip of  db")
	flag.StringVar(&dbport, "dbport", "3306", "the db port to connect db")
	flag.Int64Var(&level0, "level0", 1800, "leve 0 send frequency(seconds)")
	flag.Int64Var(&level1, "level1", 900, "leve 1 send frequency(seconds)")
	flag.Int64Var(&level2, "level2", 480, "leve 2 send frequency(seconds)")
	flag.Int64Var(&level3, "level3", 300, "leve 3 send frequency(seconds)")
	flag.Int64Var(&level4, "level4", 60, "leve 4 send frequency(seconds)")
	flag.Int64Var(&level5, "level5", 30, "leve 5 send frequency(seconds)")
	flag.StringVar(&models.WxCoreID, "wxcoreid", "XXXXXX", "weixin coreID of commpany")
	flag.StringVar(&models.WxCorpSecret, "wxsecret", "xxxxx", "weixin core secret of commpany")
	flag.StringVar(&models.Wxreceive, "wxreceive", "@all", "微信接收者")
	flag.IntVar(&models.WxagentId, "wxagentid", 1, "wenxin company number app agent id")
	flag.Parse()
	db := models.MySqlConfig{}
	db.DBName = dbname
	db.Host = dbhost
	db.MaxIdleConns = 1
	db.MaxOpenConn = 3
	db.Password = dbpass
	db.Port = dbport
	db.ShowSqlOnConsole = true
	db.UserName = dbname
	db.MySqlConn()
	models.MysqlPing()
}
func main() {
	call()
}

func call() {
	count := int64(1)
	for {
		if count == 9223372036854775806 {
			count = 1
		}
		//	fmt.Println("count=", count)
		if count%level0 == 0 {
			models.Action(0)
		}

		if count%level1 == 0 {
			models.Action(1)
		}

		if count%level2 == 0 {
			models.Action(2)
		}

		if count%level3 == 0 {
			models.Action(3)
		}

		if count%level4 == 0 {
			models.Action(4)
		}

		if count%level5 == 0 {
			models.Action(5)
		}

		time.Sleep(time.Duration(1) * time.Second)
		count++
	}
}
