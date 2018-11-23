package weixinAPI

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jetlwx/comm/httplib"
)

/*
how to use
import (
	"fmt"

	"github.com/jetlwx/comm/weixinAPI"
)

func main() {
	wx := weixinAPI.CorpInfo{}
	wx.CorpID = "wxda8ff3f444b720ae"
	wx.CorpSecret = "iv4J641a81eYeXYPKXCsWaFq"
	msg := weixinAPI.SendMsg{}
	msg.Touser = "@all"
	msg.Agentid = 1 //乐有家监控应用ID为1
	msg.Text.Content = "mymsg"

	wx.Send(msg)
}



*/
//认证信息
type CorpInfo struct {
	CorpID     string
	CorpSecret string
}
type tokenInfo struct {
	Errcode      int
	Errmsg       string
	Access_token string
	Expires_in   int
}

//发送信息
//定义一个简单的文本消息格式
/*
参数	是否必须	说明
touser	否	成员ID列表（消息接收者，多个接收者用‘|’分隔，最多支持1000个）。特殊情况：指定为@all，则向该企业应用的全部成员发送
toparty	否	部门ID列表，多个接收者用‘|’分隔，最多支持100个。当touser为@all时忽略本参数
totag	否	标签ID列表，多个接收者用‘|’分隔，最多支持100个。当touser为@all时忽略本参数
msgtype	是	消息类型，此时固定为：text
agentid	是	企业应用的id，整型。可在应用的设置页面查看
content	是	消息内容，最长不超过2048个字节
safe	否	表示是否是保密消息，0表示否，1表示是，默认0

*/
type SendMsg struct {
	Agentid int    `json:"agentid"` //应用id
	Msgtype string `json:"msgtype"`
	Safe    int    `json:"safe"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	Toparty string `json:"toparty"`
	Totag   string `json:"totag"`
	Touser  string `json:"touser"`
}

type send_msg_error struct {
	Errcode      int    `json:"errcode`
	Errmsg       string `json:"errmsg"`
	Invaliduser  string `json:"invaliduser"`
	Invalidparty string `json:"invalidparty"`
	Invalidtag   string `json:"invalidtag"`
}

func token(w *CorpInfo) (string, error) {
	url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + w.CorpID + "&corpsecret=" + w.CorpSecret
	req := httplib.Get(url)
	//fmt.Println(req.String())
	//log.Println("request url:", url)
	res, err := req.Bytes()
	if err != nil {
		log.Println("error at function token", err)
		return "", err
	}
	//log.Println("res=", res)

	t := tokenInfo{}
	if err := json.Unmarshal(res, &t); err != nil {
		log.Println("error at function token2", err)
		return "", err
	}

	if t.Errcode == 0 {
		return t.Access_token, nil
	}
	return "", errors.New("未知错误")
}

func (w *CorpInfo) Send(msg SendMsg) (ee error, info []byte) {
	if msg.Msgtype == "" {
		msg.Msgtype = "text"
	}
	//msg.Agentid = 1
	msg2, err := json.Marshal(msg)
	if err != nil {
		log.Println("error at function Send", err)
		return err, nil
	}
	fmt.Println("msg2:", string(msg2))
	body := bytes.NewBuffer(msg2)
	to, err := token(w)
	if err != nil {
		log.Println("error at function Send2", err)
		return err, nil
	}

	send_url := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + to

	resp, err := http.Post(send_url, "application/json", body)

	if resp.StatusCode != 200 {
		log.Println("error at function Send3", err)
		return errors.New(resp.Status), nil
	}
	fmt.Println(" resp.StatusCode=", resp.StatusCode)
	buf, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var e send_msg_error
	err = json.Unmarshal(buf, &e)
	if err != nil {
		log.Println("error at function Send4", err)
		return err, buf
	}
	if e.Errcode != 0 && e.Errmsg != "ok" {
		return errors.New(string(buf)), buf
	}
	return nil, msg2
}
