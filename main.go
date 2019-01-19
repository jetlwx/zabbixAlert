package main

import (
	"flag"
	"strings"
	"time"

	"github.com/jetlwx/comm"

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
	flag.Int64Var(&level0, "level0", 14400, "leve 0-未分类 send frequency(seconds)")
	flag.Int64Var(&level1, "level1", 1200, "leve 1-信息 send frequency(seconds)")
	flag.Int64Var(&level2, "level2", 240, "leve 2-警告 send frequency(seconds)")
	flag.Int64Var(&level3, "level3", 120, "leve 3-一般严重 send frequency(seconds)")
	flag.Int64Var(&level4, "level4", 90, "leve 4-严重 send frequency(seconds)")
	flag.Int64Var(&level5, "level5", 30, "leve 5-灾难 send frequency(seconds)")

	flag.StringVar(&models.WxCoreID, "wxcoreid", "XXXXXX", "weixin coreID of commpany")
	flag.StringVar(&models.WxCorpSecret, "wxsecret", "xxxxx", "weixin core secret of commpany")
	flag.StringVar(&models.Wxreceive, "wxreceive", "@all", "微信接收者")
	flag.IntVar(&models.WxagentId, "wxagentid", 1, "wenxin company number app agent id")
	flag.StringVar(&models.ExcludeKeyWord, "exclude", "kettle_ro,kettle_rw,bd_ro,cloud_101ro", "不推出的报警关键字，多个用逗号隔开")
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
	models.ExcludeKeyWord2 = strings.Split(models.ExcludeKeyWord, ",")
	comm.JetLog("I", "exclude key word:", models.ExcludeKeyWord2)
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
			go models.Action2(0)
		}

		if count%level1 == 0 {
			go models.Action2(1)
		}

		if count%level2 == 0 {
			go models.Action2(2)
		}

		if count%level3 == 0 {
			go models.Action2(3)
		}

		if count%level4 == 0 {
			go models.Action2(4)
		}

		if count%level5 == 0 {
			go models.Action2(5)
		}

		time.Sleep(time.Duration(1) * time.Second)
		count++
	}
}
