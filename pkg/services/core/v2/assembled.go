package v2

import (
	"encoding/json"
	"fmt"
	"gomessage/pkg/models"
	"gomessage/pkg/services/core/v1"
	"gomessage/pkg/services/format"
	"gomessage/pkg/utils"
)

// GeneralMessageAssembled 通用消息体封装
type GeneralMessageAssembled struct {
	Assembled
}

func (d *GeneralMessageAssembled) AssembledData(isRenders, isMerge bool, client *models.Client, contentList []string) (string, []any) {
	fmt.Println("消息体原封不动向后传递，不做任何变动...")

	var msgList []any
	for _, msg := range contentList {
		var cc map[string]any
		requestByte := []byte(msg)
		json.Unmarshal(requestByte, &cc)
		msgList = append(msgList, cc)
	}

	return "", msgList
}

// DingtalkMessageAssembled 钉钉消息体封装
type DingtalkMessageAssembled struct {
	Assembled
}

func (d *DingtalkMessageAssembled) AssembledData(isRenders, isMerge bool, client *models.Client, contentList []string) (string, []any) {
	fmt.Println("组装钉钉需要的消息体...")
	var msgList []any
	url := v1.RobotRandomUrl(client.ExtendDingtalk.RobotUrlRandomList)

	if isRenders {
		if isMerge {
			msg := v1.MessageJoint(contentList, utils.VarDingtalk)
			data := format.PackDingtalkMessage(client.ExtendDingtalk.RobotKeyword, msg, false, nil, nil)
			msgList = append(msgList, data)
		} else {
			for _, msg := range contentList {
				data := format.PackDingtalkMessage(client.ExtendDingtalk.RobotKeyword, msg, false, nil, nil)
				msgList = append(msgList, data)
			}
		}
		return url, msgList
	} else {
		for _, msg := range contentList {
			var cc map[string]any
			requestByte := []byte(msg)
			json.Unmarshal(requestByte, &cc)
			msgList = append(msgList, cc)
		}
		return url, msgList
	}
}

// FeishuMessageAssembled 飞书消息体封装
type FeishuMessageAssembled struct {
	Assembled
}

func (d *FeishuMessageAssembled) AssembledData(isRenders, isMerge bool, client *models.Client, contentList []string) (string, []any) {
	fmt.Println("组装飞书需要的消息体...")
	var msgList []any
	url := v1.RobotRandomUrl(client.ExtendFeishu.RobotUrlRandomList)

	if isRenders {
		if isMerge {
			msg := v1.MessageJoint(contentList, utils.VarFeishu)
			data := format.PackFeishuMessage(client, msg)
			msgList = append(msgList, data)
		} else {
			for _, msg := range contentList {
				data := format.PackFeishuMessage(client, msg)
				msgList = append(msgList, data)
			}
		}
		return url, msgList
	} else {
		for _, msg := range contentList {
			var cc map[string]any
			requestByte := []byte(msg)
			json.Unmarshal(requestByte, &cc)
			msgList = append(msgList, cc)
		}
		return url, msgList
	}

}

// WechatMessageAssembled 企业微信消息体封装
type WechatMessageAssembled struct {
	Assembled
}

func (d *WechatMessageAssembled) AssembledData(isRenders, isMerge bool, client *models.Client, contentList []string) (string, []any) {
	fmt.Println("组装企业微信需要的消息体...")
	var msgList []any
	url := ""

	if isRenders {
		if isMerge {
			msg := v1.MessageJoint(contentList, utils.VarWechatApplication)
			msgList = append(msgList, msg)
		} else {
			for _, msg := range contentList {
				msgList = append(msgList, msg)
			}
		}
		return url, msgList

	} else {
		for _, msg := range contentList {
			var cc map[string]any
			requestByte := []byte(msg)
			json.Unmarshal(requestByte, &cc)
			msgList = append(msgList, cc)
		}
		return url, msgList
	}

}

// EmailMessageAssembled 邮箱消息体封装
type EmailMessageAssembled struct {
	Assembled
}

func (d *EmailMessageAssembled) AssembledData(isRenders, isMerge bool, client *models.Client, contentList []string) (string, []any) {
	fmt.Println("组装为Email需要的消息体...")
	return "", nil
}
