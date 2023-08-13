package v3

import (
	"encoding/json"
	"gomessage/pkg/models"
	"gomessage/pkg/services/core/v1"
	"gomessage/pkg/services/format"
	"gomessage/pkg/utils"
)

type ClientActionDingtalk struct {
	Client *models.Client
}

func (c *ClientActionDingtalk) RendersMessages(client *models.Client, isMerge bool, contentList []string) []any {

	//TODO: 这里开发出来@人的功能，感觉没啥意义，因为Markdown格式的信息不支持@人高亮显示
	//参考文档：https://open.dingtalk.com/document/robots/custom-robot-access
	//参考文档2：https://m.dingtalk.com/qidian/help-detail-1060904202

	//var mobiles []string = []string{"18516187075", "13071003191"} //手机号
	//var dingtalkIds []string = []string{"rpn519f"} //钉钉号（注意：不是userID，而是钉钉号，一定要注意区分）

	var mobiles, dingtalkIds []string
	isAtAll := client.ExtendDingtalk.IsAtAll
	json.Unmarshal(client.ExtendDingtalk.AtMobiles, &mobiles)
	json.Unmarshal(client.ExtendDingtalk.AtDingtalkIds, &dingtalkIds)

	var msgList []any
	//判断是否把消息聚合后发送
	if isMerge {
		msg := v1.MessageJoint(contentList, utils.VarDingtalk)
		data := format.PackDingtalkMessage(client.ExtendDingtalk.RobotKeyword, msg, isAtAll, mobiles, dingtalkIds)
		msgList = append(msgList, data)
	} else {
		for _, msg := range contentList {
			data := format.PackDingtalkMessage(client.ExtendDingtalk.RobotKeyword, msg, isAtAll, mobiles, dingtalkIds)
			msgList = append(msgList, data)
		}
	}
	return msgList
}

func (c *ClientActionDingtalk) PushMessages(messages []any) {
	url := v1.RobotRandomUrl(c.Client.ExtendDingtalk.RobotUrlRandomList)
	for _, msg := range messages {
		v1.Push(msg, url)
	}
}
