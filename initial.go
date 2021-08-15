package WoComApi

import (
	"github.com/thep0y/predator"
	"github.com/tidwall/gjson"
	"time"
)

type WoCom struct {
	corpId  string
	secret  string
	agentId string
}

type Message interface {
	// MessageText 文本消息
	MessageText(content string) error
	// MessageTextCard 卡片消息 title标题，description内容（支持html），Url链接，btntxt按钮文字
	MessageTextCard(title, description, URL, btntxt string) error
	// MessageMarkdown markdown消息
	MessageMarkdown(content string) error
}

var (
	crawler *predator.Crawler
	token   string
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

// MessageText 文本消息
func (m *WoCom) MessageText(content string) error {
	data := `{
   "touser" : "@all",
   "toparty" : "@all",
   "totag" : "@all",
   "msgtype" : "text",
   "agentid" : ` + m.agentId + `,
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

// MessageTextCard 卡片消息 title标题，description内容（支持html），Url链接，btntxt按钮文字
func (m *WoCom) MessageTextCard(title, description, URL, btntxt string) error {
	data := `{
   "touser" : "@all",
   "toparty" : "@all",
   "totag" : "@all",
   "msgtype" : "textcard",
   "agentid" : ` + m.agentId + `,
   "textcard" : {
            "title" : "` + title + `",
            "description" : "` + description + `",
            "url" : "` + URL + `",
                        "btntxt":"` + btntxt + `"
   },
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

// MessageMarkdown markdown消息
func (m *WoCom) MessageMarkdown(content string) error {
	data := `{
		"touser" : "@all",
			"toparty" : "@all",
			"totag" : "@all",
			"msgtype": "markdown",
			"agentid" : ` + m.agentId + `,
			"markdown": {
			"content": "` + content + `"
	},
	"enable_duplicate_check": 0,
	"duplicate_check_interval": 1800
	}`

	err := crawler.PostRaw("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="+token, []byte(data), nil)
	if err != nil {
		return err
	}

	return nil
}
