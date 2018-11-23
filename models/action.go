package models

import (
	"time"

	"github.com/jetlwx/comm/weixinAPI"
)

var WxCoreID, WxCorpSecret, Wxreceive string
var WxagentId int

func Action(priority int) {
	var msg string
	//get triggers that no recover
	triggers := NoRecoverTriggersID(priority)

	for _, tr := range triggers {
		//get items monitor status,if status is 1(diabled) then  continues

		//if status is 0 ,mean enable;
		itemid := GetItemIdByTriggerID(tr)
		if itemid == 0 {
			continue
		}
		//// itemstatus 0-->enable 1--> disable
		itemStatus := GetItemStatus(itemid)
		if itemStatus != 0 {
			continue
		}

		//check the host or monitor 0-->监控中，１－－》未被监控中
		if GetHostStatus(itemid) != 0 {
			continue
		}

		eid := GetEventsID(tr)
		if eid == 0 {
			continue
		}
		aler := GetAlerts(eid)
		if len(aler) == 0 {
			continue
		}
		for _, a := range aler {
			if a.Subject == "" && a.Message == "" {
				continue
			}

			t0 := time.Now().Format("2006-01-02 15:04:05")
			msg += a.Subject + "\n"
			msg += "推送时间:" + t0 + "\n"
			msg += a.Message + "\n-----------------------" + "\n"
		}

	}

	if msg == "" {
		return
	}
	sendms(msg)
}

func sendms(weixinmsg string) {
	wx := weixinAPI.CorpInfo{}
	wx.CorpID = WxCoreID
	wx.CorpSecret = WxCorpSecret
	msg := weixinAPI.SendMsg{}
	msg.Touser = Wxreceive
	msg.Text.Content = weixinmsg
	msg.Agentid = WxagentId
	//fmt.Println(wx.Send(msg))

	wx.Send(msg)
}
