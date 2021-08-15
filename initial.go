package WoComApi

import (
	"github.com/thep0y/predator"
	"github.com/tidwall/gjson"
	"time"
)

type WoCom struct {
	corpId   string
	secret   string
	agentId  string
	Messages messages
}

type messages struct {
}

var (
	crawler *predator.Crawler
	token   string
	agentId string
)

func init() {
	//初始化爬虫框架
	crawler = predator.NewCrawler()

}

// NewWoCom 使用一些 WoComOption 创建一个新的 WoCom 实例
func NewWoCom(opts ...WoComOption) *WoCom {
	//创建WoCom
	c := new(WoCom)
	for _, op := range opts {
		op(c)
	}
	agentId = c.agentId
	//Token
	c.getToken()
	go c.updateToken()

	return c
}

/************************* token处理 ****************************/
//获取Token
func (w *WoCom) getToken() {
	var body string
	//爬虫
	crawler.AfterResponse(func(r *predator.Response) {
		body = string(r.Body)
	})
	crawler.Get("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + w.corpId + "&corpsecret=" + w.secret)

	//提取json中的token
	result := gjson.Get(body, "access_token")
	token = result.Str
}

//更新一个小时获取一次
func (w *WoCom) updateToken() {
	//定时刷新token
	for {
		time.Sleep(time.Millisecond * 3600000)
		w.getToken()
	}

}

/************************* 发送消息 ****************************/

// Text 文本消息
func (m *messages) TextAll(content string) error {
	data := `{
   "touser" : "@all",
   "toparty" : "@all",
   "totag" : "@all",
   "msgtype" : "text",
   "agentid" : ` + agentId + `,
   "text" : {
       "content" : ` + `"` + content + `""` + `
   },
   "safe":0,
   "enable_id_trans": 0,
   "enable_duplicate_check": 0,
   "duplicate_check_interval": 1800
}`

	err := crawler.PostRaw("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="+token, []byte(data), nil)
	if err != nil {
		return err
	}
	return nil
}
