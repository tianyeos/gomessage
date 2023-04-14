package clients

import (
	"gomessage/apps/controllers/clients"
	"gomessage/apps/controllers/send"
	"gomessage/apps/models"
)

type ClientActionDingtalk struct {
	Client *models.Client
}

func (c *ClientActionDingtalk) RendersMessages(client *models.Client, isMerge bool, contentList []string) []any {
	var msgList []any
	if isMerge {
		msg := send.MessageJoint(contentList, "dingtalk")
		data := clients.PackDingtalkMessage(client.ExtendDingtalk.RobotKeyword, msg)
		msgList = append(msgList, data)
	} else {
		for _, msg := range contentList {
			data := clients.PackDingtalkMessage(client.ExtendDingtalk.RobotKeyword, msg)
			msgList = append(msgList, data)
		}
	}
	return msgList
}

func (c *ClientActionDingtalk) PushMessages(messages []any) {
	url := send.RobotRandomUrl(c.Client.ExtendDingtalk.RobotUrlInfoList)
	for _, msg := range messages {
		send.Push(msg, url)
	}
}
