package format

import "fmt"

func PackDingtalkMessage(keyword string, message string, atAll bool, mobiles []string, dingtalkIds []string) interface{} {
	type MarkdownMessage struct {
		Msgtype  string `json:"msgtype"`
		Markdown struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		} `json:"markdown"`
		At struct {
			AtMobiles     []string `json:"atMobiles"`
			AtDingtalkIds []string `json:"atDingtalkIds"`
			IsAtAll       bool     `json:"isAtAll"`
		} `json:"at"`
	}
	m := MarkdownMessage{}
	m.Msgtype = "markdown"
	m.Markdown.Title = keyword + "-GoMessage"
	m.Markdown.Text = message

	if atAll { //判断是否要@所有人，如果为true，则人员名单设置为空
		m.At.IsAtAll = true //@所有人，本质上就相当于"谁都不@"，考虑到小屏设备的数据对齐效果，这里就真的谁都不@了
		m.At.AtMobiles = nil
		m.At.AtDingtalkIds = nil
	} else {
		//关闭艾特所有人
		m.At.IsAtAll = false

		//设定底层艾特功能
		m.At.AtMobiles = mobiles
		m.At.AtDingtalkIds = dingtalkIds

		//把艾特目标追加到文本中
		if len(mobiles) != 0 || len(dingtalkIds) != 0 {
			//反正钉钉的Markdown策略也不会让被艾特的人"蓝色高亮"显示，因此这里是使用了引用（>），效果上肯能会更美观一些
			m.Markdown.Text += "   \n   &emsp;   \n   >"
			appendAt := func(list []string) {
				for _, item := range list {
					if item != "" {
						m.Markdown.Text += fmt.Sprintf("@%s ", item)
					}
				}
			}
			appendAt(mobiles)
			appendAt(dingtalkIds)
		}

	}

	return m
}
