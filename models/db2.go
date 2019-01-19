package models

import (
	"strconv"
	"time"

	"github.com/jetlwx/comm"
)

//Item value type
/*
0 - numeric float;
1 - character;
2 - log;
3 - numeric unsigned;
4 - text.
*/
type Item struct {
	Itemid      int64
	Hostid      int64
	Name        string
	ValueType   int
	Interfaceid int64
	Status      int //是否启用
}
type Host struct {
	Ip   string
	Host string
	Name string
}
type AlertTriggers struct {
	Triggerid   int64
	Description string
}

type Problem struct {
	Eventid     int64
	Objectid    int64 //is Triggerid
	Description string
}

// NoRecoverProblems 找出未恢复的事件
/*
MariaDB [zabbix]> select eventid,objectid from problem where r_eventid is null;
+---------+----------+
| eventid | objectid |
+---------+----------+
| 4693861 |    22789 |
+---------+----------+
status=0 tiggers is enable
*/
func NoRecoverProblems(priority int) (p []Problem) {
	sqlstr := "select problem.eventid,problem.objectid, triggers.description from problem,triggers,events" + " "
	sqlstr += " where problem.r_eventid is null and triggers.status=0 and triggers.state=0 and triggers.priority=" + strconv.Itoa(priority) + " "
	sqlstr += " and problem.objectid=triggers.triggerid  and events.acknowledged=0 and problem.eventid=events.eventid"
	if err := Engine.Sql(sqlstr).Find(&p); err != nil {
		comm.JetLog("E", err)
		return
	}
	comm.JetLog("D", p)
	return
}

//ItemInfo 获取item信息
/*
MariaDB [zabbix]> select itemid,hostid,name,value_type,interfaceid,status from items where itemid=(select itemid from functions  where triggerid=15846 limit 1);
+--------+--------+--------------------------------+------------+-------------+--------+
| itemid | hostid | name                           | value_type | interfaceid | status |
+--------+--------+--------------------------------+------------+-------------+--------+
|  29593 |  10289 | 当前线程等待最大时间           |          3 |          20 |      0 |
+--------+--------+--------------------------------+------------+-------------+--------+
*/
func (p Problem) ItemInfo() (i Item) {
	sql := "select itemid,hostid,name,value_type,interfaceid,status from items where status=0 and itemid=(select itemid from functions  where triggerid=" + strconv.FormatInt(p.Objectid, 10) + " limit 1)"
	_, err := Engine.Sql(sql).Get(&i)
	if err != nil {
		comm.JetLog("E", err)
		return
	}

	return
}

//EventTime　is 事件时间
func (p Problem) EventTime() (t string) {
	var i int64
	sql := "select clock from events where eventid= " + strconv.FormatInt(p.Eventid, 10)
	_, err := Engine.Sql(sql).Get(&i)
	if err != nil {
		comm.JetLog("E", err)
		return
	}

	return time.Unix(i, 0).Format("2006-01-02 15:04:05")
}

//获取主机信息
/*
MariaDB [zabbix]> select interface.ip, hosts.host,hosts.name from interface,hosts where hosts.hostid=10268 and hosts.hostid=interface.hostid;
+---------------+-------------+----------------------------+
| ip            | host        | name                       |
+---------------+-------------+----------------------------+
| 192.168.200.2 | kubenode200 | 容器宿主机（200.2）        |
+---------------+-------------+----------------------------+

*/
func (i Item) HostInfo() (h Host) {
	sql := "select interface.ip, hosts.host,hosts.name from interface,hosts where hosts.status=0 and hosts.hostid=" + strconv.FormatInt(i.Hostid, 10) + " and hosts.hostid=interface.hostid"
	_, err := Engine.Sql(sql).Get(&h)
	if err != nil {
		comm.JetLog("E", err)
		return
	}

	return
}

/*
Possible values:
0 - numeric float;
1 - character;
2 - log;
3 - numeric unsigned;
4 - text.
| history ---0                   |
| history_log  ---2              |
| history_str  ---1              |
| history_text  ---4             |
| history_uint ---3

*/
func (i Item) IetmlatestValue() (res string) {
	if i.Itemid == 0 {
		return
	}
	switch i.ValueType {
	case 0:
		var v float64
		sql := "select  value from history where itemid=" + strconv.FormatInt(i.Itemid, 10) + "  order by clock desc limit 1"
		_, err := Engine.Sql(sql).Get(&v)
		if err != nil {
			comm.JetLog("E", err)
			return
		}
		return strconv.FormatFloat(v, 'f', -1, 64)

	case 1:
		var v string
		sql := "select  value from history_str where itemid=" + strconv.FormatInt(i.Itemid, 10) + "  order by clock desc limit 1"
		_, err := Engine.Sql(sql).Get(&v)
		if err != nil {
			comm.JetLog("E", err)
			return
		}
		return v

	case 2:
		var v string
		sql := "select  value from history_log where itemid=" + strconv.FormatInt(i.Itemid, 10) + "  order by clock desc limit 1"
		_, err := Engine.Sql(sql).Get(&v)
		if err != nil {
			comm.JetLog("E", err)
			return
		}
		return v

	case 3:
		var v int64
		sql := "select  value from history_uint where itemid=" + strconv.FormatInt(i.Itemid, 10) + "  order by clock desc limit 1"
		_, err := Engine.Sql(sql).Get(&v)
		if err != nil {
			comm.JetLog("E", err)
			return strconv.FormatInt(0, 10)
		}
		return strconv.FormatInt(v, 10)

	case 4:
		var v string
		sql := "select  value from history_text where itemid=" + strconv.FormatInt(i.Itemid, 10) + "  order by clock desc limit 1"
		_, err := Engine.Sql(sql).Get(&v)
		if err != nil {
			comm.JetLog("E", err)
			return
		}
		return v

	}

	return

}
