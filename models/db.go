package models

import (
	"log"
	"strconv"
	"time"
)

// type Alerts struct {
// 	Alertid       int64
// 	Actionid      int64
// 	Eventid       int64
// 	Userid        int64
// 	Clock         int
// 	Mediatypeid   int64
// 	Sendto        string
// 	Subject       string
// 	Message       string
// 	Status        int
// 	Retries       int
// 	Error         string
// 	Esc_step      int
// 	Alerttype     int
// 	P_eventid     int64
// 	Acknowledgeid int64
// }
type Alerts struct {
	Subject string
	Message string
}

//得取最近３个月未恢复指定级别的触发器ＩＤ
//0-未分类，1-信息，2-警告，3-一般严重，4-严重。5-灾难）
func NoRecoverTriggersID(priority int) (t []int64) {
	t1 := strconv.FormatInt(time.Now().Unix()-int64(90*24*36000), 10)

	sqlstr := "select triggerid from triggers  where priority=" + strconv.Itoa(priority) + " and  status=0 AND value=1 and lastchange>" + t1
	if err := Engine.Sql(sqlstr).Find(&t); err != nil {
		log.Println("error at function NoRecoverTriggersID", err)
		return
	}
	log.Println("trigers=", t)
	return
}

//根据triggerid获取所对应itemid
//根据itemid得到当前监控项是否启用，若启用进行下一步操作，否则放弃当前trigger
// 0-->enable 1--> disable
//获取主机是否有被监控中 0-->监控中，１－－》未被监控中

func GetItemIdByTriggerID(triggerid int64) (itemid int64) {
	str := "select itemid from functions  where triggerid=" + strconv.FormatInt(triggerid, 10) + " limit 1"
	has, err := Engine.Sql(str).Get(&itemid)
	if err != nil {
		log.Println("error at function GetItemIDByTriggerID", err)
		return 0
	}
	if !has {
		return 0
	}
	log.Println("itemid=", itemid)
	return
}

func GetItemStatus(itemid int64) (status int) {
	str := "select status from items where itemid=" + strconv.FormatInt(itemid, 10)
	has, err := Engine.Sql(str).Get(&status)
	if err != nil {
		log.Println("error at function GetItemStatus", err)
		return 128
	}
	if !has {
		return 0
	}
	log.Println("itemid status=", status)
	return
}

func GetHostStatus(itemid int64) (status int) {
	str := "select status from hosts where hostid=(select hostid from items where itemid=" + strconv.FormatInt(itemid, 10) + ")"
	has, err := Engine.Sql(str).Get(&status)
	if err != nil {
		log.Println("error at function GetItemStatus", err)
		return 128
	}
	if !has {
		return 0
	}
	log.Println("host status=", status)
	return
}

//根据triggerid获取所对应itemid
//根据itemid得到当前监控项是否启用，若启用进行下一步操作，否则放弃当前trigger
// 0-->enable 1--> disable
//获取主机是否有被监控中 0-->监控中，１－－》未被监控中
// func GetItemStatus(triggerid int64) (status int) {
// 	var itemstatus int
// 	//sqlstr2 := "select status from hosts where hostid=(select hostid from items where itemid=(select itemid from functions  where triggerid=" + strconv.FormatInt(triggerid, 10) + " limit 1))"
// 	sqlstr := "select status from items where itemid=(select itemid from functions  where triggerid=" + strconv.FormatInt(triggerid, 10) + " limit 1)"
// 	_, err := Engine.Sql(sqlstr2).Get(&itemstatus)
// 	if err != nil {
// 		log.Println("error at function GetItemIDByTriggerID", err)
// 		return 128
// 	}
// 	//如果item为启用，检查host监控是否启用
// 	sqlstr2:=
// 	return
// }

//若 items status 状态为0 ，则找到对应events的EVNETSID，
func GetEventsID(triggerid int64) (eid int64) {
	sqlstr := "select eventid from  events where objectid=" + strconv.FormatInt(triggerid, 10) + " and acknowledged <>1 order by eventid desc limit 1"
	has, err := Engine.Sql(sqlstr).Get(&eid)
	if err != nil {
		log.Println("error at function GetEventsID", err)
		return 0
	}
	if !has {
		return 0
	}
	log.Println("events id =", eid)

	return
}

//根据EVENTSID获取alerts表的报警信息
func GetAlerts(eid int64) (a []Alerts) {
	sqlstr := "select distinct(subject),message from alerts where eventid=" + strconv.FormatInt(eid, 10)
	if err := Engine.Sql(sqlstr).Find(&a); err != nil {
		log.Println("error at function GetAlerts", err)
		return
	}

	return
}
