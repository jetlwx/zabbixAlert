package models

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jetlwx/comm"

	"github.com/jetlwx/comm/weixinAPI"
)

var ExcludeKeyWord string
var ExcludeKeyWord2 []string

type Msg struct {
	EventName  string
	ActionTime string
	PushTime   string
	Host       string
	Hostname   string
	IP         string
	Value      string
}

func Action2(priority int) {
	M := []Msg{}
	problem := NoRecoverProblems(priority)
	var msg string
	for _, v := range problem {
		//获取item信息
		if v.Objectid == 0 {
			continue
		}

		i := v.ItemInfo()
		if i.Hostid == 0 {
			continue
		}
		//在维护期
		if IsInMaintenancesPeriods(i.Hostid) {
			comm.JetLog("W", "在维护期间，不推送,hostid=", i.Hostid)
			continue
		}
		//host info
		h := i.HostInfo()
		if h.Host == "" {
			continue
		}
		m := Msg{}
		//leastvalues
		value := i.IetmlatestValue()
		actime := v.EventTime()
		t0 := time.Now().Format("2006-01-02 15:04:05")
		m.EventName = v.Description
		m.ActionTime = actime
		m.PushTime = t0
		m.Host = h.Name
		m.Hostname = h.Host
		m.IP = h.Ip
		m.Value = value
		M = append(M, m)

	}

	N := DuplicateRemoval(M)
	log.Printf("%#v", N)
	for _, v := range N {
		msg += "事件:" + v.EventName + "\n"
		msg += "发生时间:" + v.ActionTime + "\n"
		msg += "推送时间:" + v.PushTime + "\n"
		msg += "主机：" + v.Host + "\n"
		msg += "主机名:" + v.Hostname + "\n"
		msg += "IP:" + v.IP + "\n"
		msg += "当前值：" + v.Value + "\n"
		msg += "－－－－－－－－－－－－－" + "\n"
	}

	if msg == "" {
		return
	}
	sendms(msg)
}

//去重
func DuplicateRemoval(m []Msg) (n []Msg) {
	lenm := len(m)

	for i := 0; i < lenm; i++ {
		log.Printf("%#v", m[i])
		//如果有，则不推送
		if inNotPush(m[i]) {
			continue
		}
		repeat := false
		for j := i + 1; j < lenm; j++ {
			if m[i].EventName == m[j].EventName && m[i].ActionTime == m[j].ActionTime {
				break
			}
		}
		if !repeat {
			n = append(n, m[i])
		}
	}

	return n
}

//不推出去报警关键字
func inNotPush(m Msg) bool {
	for _, v := range ExcludeKeyWord2 {
		if strings.Contains(strings.ToLower(m.Value), strings.ToLower(v)) {
			comm.JetLog("I", "匹配到关键字:", v)
			return true
		}
	}
	fmt.Printf("%#v", m)
	return false
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
