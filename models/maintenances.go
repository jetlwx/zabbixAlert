package models

import (
	"strconv"
	"time"

	"github.com/jetlwx/comm"
)

type HostMaintenances struct {
	Timeperiod_type int
	Every           int
	Month           int
	Dayofweek       int
	Day             int
	Start_time      int64
	Period          int64
	Start_date      int64
	Active_since    int64
	Active_till     int64
}

func IsInMaintenancesPeriods(hostid int64) bool {
	h, b := hostInMaintenancesInfo(hostid)
	g, b1 := hostInMaintenancesGroup(hostid)
	if b {
		if hostInPeriod(h) {
			return true
		}
	}

	if b1 {
		if hostInPeriod(g) {
			return true
		}
	}

	return false
}

func hostInPeriod(h HostMaintenances) bool {
	//只判断每天执行的时间范围
	t0 := time.Now()
	t1 := t0.Unix()
	h1 := t0.Hour()
	m1 := t0.Minute()
	s1 := t0.Second()
	atTime := int64(h1*3600 + m1*60 + s1)
	if t1 < h.Active_since || t1 > h.Active_till {
		return false
	}

	if h.Timeperiod_type == 2 && h.Every == 1 && h.Day == 1 {
		if atTime < h.Start_time || atTime > (h.Start_time+h.Period) {
			return false
		}
		return true
	}
	return false
}

/*
MariaDB [zabbix]> select   timeperiods.timeperiod_type,timeperiods.every,timeperiods.month,timeperiods.dayofweek,timeperiods.day,timeperiods.start_time,timeperiods.period,timeperiods.start_date,maintenances.active_since,maintenances.active_till from maintenances_hosts,timeperiods,maintenances_windows,maintenances where timeperiods.timeperiodid=maintenances_windows.timeperiodid and maintenances_windows.maintenanceid=maintenances.maintenanceid and maintenances_hosts.maintenanceid=maintenances_windows.maintenanceid and maintenances_hosts.hostid=10295;
+-----------------+-------+-------+-----------+-----+------------+--------+------------+--------------+-------------+
| timeperiod_type | every | month | dayofweek | day | start_time | period | start_date | active_since | active_till |
+-----------------+-------+-------+-----------+-----+------------+--------+------------+--------------+-------------+
|               2 |     1 |     0 |         0 |   1 |      82800 |  28800 | 1528731000 |   1528646400 |  2002118400 |
+-----------------+-------+-------+-----------+-----+------------+--------+------------+--------------+-------------+
*/
func hostInMaintenancesInfo(hostid int64) (h HostMaintenances, has bool) {
	sqlstr := "select timeperiods.timeperiod_type,timeperiods.every,timeperiods.month,timeperiods.dayofweek,timeperiods.day,timeperiods.start_time,timeperiods.period,"
	sqlstr += " " + "timeperiods.start_date,maintenances.active_since,maintenances.active_till from maintenances_hosts,timeperiods,maintenances_windows,"
	sqlstr += " " + "maintenances where timeperiods.timeperiodid=maintenances_windows.timeperiodid and maintenances_windows.maintenanceid=maintenances.maintenanceid"
	sqlstr += " " + "  and maintenances_hosts.maintenanceid=maintenances_windows.maintenanceid and maintenances_hosts.hostid=" + strconv.FormatInt(hostid, 10)
	b, err := Engine.Sql(sqlstr).Get(&h)
	if err != nil {
		comm.JetLog("E", err)
	}

	return h, b
}

func hostInMaintenancesGroup(hostid int64) (h HostMaintenances, has bool) {
	sqlstr := "select timeperiods.timeperiod_type,timeperiods.every,timeperiods.month,timeperiods.dayofweek,timeperiods.day,timeperiods.start_time,timeperiods.period,"
	sqlstr += "timeperiods.start_date,maintenances.active_since,maintenances.active_till from timeperiods,maintenances where timeperiodid=(select timeperiodid from "
	sqlstr += " " + " maintenances_windows where maintenanceid=(select maintenanceid from maintenances_groups where groupid in (select groupid from hosts_groups"
	sqlstr += " " + "  where hostid=" + strconv.FormatInt(hostid, 10) + ")) and maintenances.maintenanceid=maintenances_windows.maintenanceid)"
	b, err := Engine.Sql(sqlstr).Get(&h)
	if err != nil {
		comm.JetLog("E", err)

	}
	return h, b
}
