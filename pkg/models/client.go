package models

import (
	"encoding/json"
	"errors"
	"gomessage/pkg/utils"
	"gomessage/pkg/utils/database"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Client struct {
	ID                      int                `json:"id" gorm:"primarykey"`
	CreatedAt               time.Time          `json:"-"`
	UpdatedAt               time.Time          `json:"-"`
	DeletedAt               gorm.DeletedAt     `json:"-" gorm:"index"`
	Namespace               string             `json:"namespace"`                 //命名空间（也就是用户使用时看到的通道名称）
	ClientName              string             `json:"client_name"`               //客户端名称
	ClientDescription       string             `json:"client_description"`        //客户端描述
	ClientType              string             `json:"client_type"`               //客户端类型
	IsActive                bool               `json:"is_active"`                 //客户端是否激活
	ClientInfo              json.RawMessage    `json:"client_info" gorm:"-:all" ` //客户端详情（此时拿到的是json格式的数据，需要额外的逻辑把json处理成对应的结构体，存放到"扩展表"中）
	ExtendDingtalk          *Dingtalk          `json:"-" gorm:"-:all"`            //扩展表：存放`钉钉群机器人`客户端的扩展数据
	ExtendFeishu            *Feishu            `json:"-" gorm:"-:all"`            //扩展表：存放`飞书群机器人`客户端的扩展数据
	ExtendWechatApplication *WechatApplication `json:"-" gorm:"-:all"`            //扩展表：存放`企业微信应用号`客户端的扩展数据
	ExtendWechatRobot       *WechatRobot       `json:"-" gorm:"-:all"`            //扩展表：存放`企业微信群机器人`客户端的扩展数据
}

func (*Client) TableName() string {
	return "clients"
}

func AddClient(generalClient *Client) (*Client, error) {

	//TODO: 保存client客户端的"通用数据"
	generalClient.IsActive = false
	createResult := database.DB.Default.Create(&generalClient)
	if createResult.Error != nil {
		return generalClient, createResult.Error
	}

	//TODO: 保存client客户端的"扩展数据"
	switch generalClient.ClientType {

	//钉钉群机器人
	case utils.VarDingtalk:
		generalClient.ExtendDingtalk.ClientId = int(generalClient.ID)                                        //指定扩展表中的client_id，该id就是通用client的通用id
		generalClient.ExtendDingtalk.RobotUrlRandomList = JoinUrl(generalClient.ExtendDingtalk.RobotUrlList) //url随机列表
		generalClient.ExtendDingtalk.RobotUrl = strings.Join(generalClient.ExtendDingtalk.RobotUrlRandomList, "\n")

		//generalClient.ExtendDingtalk.IsAtAll = generalClient.ExtendDingtalk.IsAtAll
		//generalClient.ExtendDingtalk.AtMobiles = generalClient.ExtendDingtalk.AtMobiles
		//generalClient.ExtendDingtalk.AtDingtalkIds = generalClient.ExtendDingtalk.AtDingtalkIds

		dingtalkResult := database.DB.Default.Create(&generalClient.ExtendDingtalk)
		if dingtalkResult.Error != nil {
			return generalClient, dingtalkResult.Error
		}

	//飞书群机器人
	case utils.VarFeishu:
		generalClient.ExtendFeishu.ClientId = int(generalClient.ID)
		generalClient.ExtendFeishu.RobotUrlRandomList = JoinUrl(generalClient.ExtendFeishu.RobotUrlList) //url随机列表
		generalClient.ExtendFeishu.RobotUrl = strings.Join(generalClient.ExtendFeishu.RobotUrlRandomList, "\n")
		feishuResult := database.DB.Default.Create(&generalClient.ExtendFeishu)
		if feishuResult.Error != nil {
			return generalClient, feishuResult.Error
		}

	//企业微信群机器人
	case utils.VarWechatRobot:
		generalClient.ExtendWechatRobot.ClientId = int(generalClient.ID)
		generalClient.ExtendWechatRobot.RobotUrlRandomList = JoinUrl(generalClient.ExtendWechatRobot.RobotUrlList) //url随机列表
		generalClient.ExtendWechatRobot.RobotUrl = strings.Join(generalClient.ExtendWechatRobot.RobotUrlRandomList, "\n")
		result := database.DB.Default.Create(&generalClient.ExtendWechatRobot)
		if result.Error != nil {
			return generalClient, result.Error
		}

	//企业微信应用号
	case utils.VarWechatApplication:
		generalClient.ExtendWechatApplication.ClientId = int(generalClient.ID)
		wechatResult := database.DB.Default.Create(&generalClient.ExtendWechatApplication)
		if wechatResult.Error != nil {
			return generalClient, wechatResult.Error
		}

	//不属于以上任何客户端类型
	default:
		return generalClient, errors.New("未知的ClientType=" + generalClient.ClientType)
	}

	//给出返回值
	return generalClient, nil
}

func UpdateClientInfo(id int, newClient *Client) (*Client, error) {
	oldClient := Client{}
	result := database.DB.Default.Model(&oldClient).Where("id = ? ", id).First(&oldClient)
	if result.Error != nil {
		return nil, result.Error
	}

	database.DB.Default.Model(&oldClient).Updates(Client{
		ClientName:        newClient.ClientName,
		ClientDescription: newClient.ClientDescription,
		IsActive:          newClient.IsActive,
	})

	switch newClient.ClientType {
	case utils.VarDingtalk:
		dingtalk := Dingtalk{}
		database.DB.Default.Model(&dingtalk).Where("client_id = ?", oldClient.ID).First(&dingtalk)
		newClient.ExtendDingtalk.RobotUrlRandomList = JoinUrl(newClient.ExtendDingtalk.RobotUrlList) //url随机列表
		database.DB.Default.Model(&dingtalk).Updates(Dingtalk{
			RobotKeyword: newClient.ExtendDingtalk.RobotKeyword,
			RobotUrl:     strings.Join(newClient.ExtendDingtalk.RobotUrlRandomList, "\n"),
		})

	case utils.VarFeishu:
		feishu := Feishu{}
		database.DB.Default.Model(&feishu).Where("client_id = ?", oldClient.ID).First(&feishu)
		newClient.ExtendFeishu.RobotUrlRandomList = JoinUrl(newClient.ExtendFeishu.RobotUrlList) //url随机列表
		database.DB.Default.Model(&feishu).Updates(Feishu{
			RobotKeyword: newClient.ExtendFeishu.RobotKeyword,
			TitleColor:   newClient.ExtendFeishu.TitleColor,
			RobotUrl:     strings.Join(newClient.ExtendFeishu.RobotUrlRandomList, "\n"),
		})

	case utils.VarWechatRobot:
		wechatRobot := WechatRobot{}
		database.DB.Default.Model(&wechatRobot).Where("client_id = ?", oldClient.ID).First(&wechatRobot)
		newClient.ExtendWechatRobot.RobotUrlRandomList = JoinUrl(newClient.ExtendWechatRobot.RobotUrlList) //url随机列表
		database.DB.Default.Model(&wechatRobot).Updates(WechatRobot{
			RobotKeyword: newClient.ExtendWechatRobot.RobotKeyword,
			RobotUrl:     strings.Join(newClient.ExtendWechatRobot.RobotUrlRandomList, "\n"),
		})

	case utils.VarWechatApplication:
		wechatApp := WechatApplication{}
		database.DB.Default.Model(&wechatApp).Where("client_id = ?", oldClient.ID).First(&wechatApp)
		database.DB.Default.Model(&wechatApp).Updates(WechatApplication{
			CorpId:  newClient.ExtendWechatApplication.CorpId,
			AgentId: newClient.ExtendWechatApplication.AgentId,
			Secret:  newClient.ExtendWechatApplication.Secret,
			Touser:  newClient.ExtendWechatApplication.Touser,
		})

	default:
		return newClient, errors.New("未知的ClientType=" + newClient.ClientType)
	}

	return &oldClient, result.Error

}

func DeleteClient(id int) (int, error) {
	var cli Client
	result := database.DB.Default.Delete(&cli, id)
	return int(result.RowsAffected), result.Error
}

func UpdateClientActive(id int, t *Client) (*Client, error) {
	client := Client{}
	updateResult := database.DB.Default.Model(&client).Where("id = ? ", id).Update("is_active", t.IsActive)
	return &client, updateResult.Error

}

func GetClientById(id int) (*Client, error) {
	var cli Client
	queryResult := database.DB.Default.Where("id = ?", id).First(&cli)
	if queryResult.Error != nil {
		return &cli, queryResult.Error
	}

	switch cli.ClientType {
	case utils.VarDingtalk:
		dingtalk := Dingtalk{}
		dingtalkResult := database.DB.Default.Where("client_id = ?", int(cli.ID)).First(&dingtalk)
		if dingtalkResult.Error != nil {
			return &cli, dingtalkResult.Error
		}
		dingtalk.RobotUrlRandomList = strings.Split(dingtalk.RobotUrl, "\n") //从数据库中取出机器人地址时，展开赋值给RobotUrlInfoList
		cli.ExtendDingtalk = &dingtalk

	case utils.VarFeishu:
		feishu := Feishu{}
		feishuResult := database.DB.Default.Where("client_id = ?", int(cli.ID)).First(&feishu)
		if feishuResult.Error != nil {
			return &cli, feishuResult.Error
		}
		feishu.RobotUrlRandomList = strings.Split(feishu.RobotUrl, "\n")
		cli.ExtendFeishu = &feishu

	case utils.VarWechatRobot:
		wechatRobot := WechatRobot{}
		wechatRobotResult := database.DB.Default.Where("client_id = ?", int(cli.ID)).First(&wechatRobot)
		if wechatRobotResult.Error != nil {
			return &cli, wechatRobotResult.Error
		}
		wechatRobot.RobotUrlRandomList = strings.Split(wechatRobot.RobotUrl, "\n")
		cli.ExtendWechatRobot = &wechatRobot

	case utils.VarWechatApplication:
		wechatApplication := WechatApplication{}
		wechatApplicationResult := database.DB.Default.Where("client_id = ?", int(cli.ID)).First(&wechatApplication)
		if wechatApplicationResult.Error != nil {
			return &cli, wechatApplicationResult.Error
		}
		cli.ExtendWechatApplication = &wechatApplication

	default:
		return nil, errors.New("未知的ClientType=" + cli.ClientType)
	}

	return &cli, nil
}

func ListClient(ns string) (*[]Client, error) {
	var list []Client
	result := database.DB.Default.Where(&Client{Namespace: ns}).Order("id desc").Find(&list)
	if result.Error != nil {
		return &list, result.Error
	}
	return &list, nil
}

// GetActiveClient 获取指定命名空间下的处于激活状态的客户端
func GetActiveClient(ns string) ([]Client, error) {
	var list []Client
	result := database.DB.Default.Where(&Client{IsActive: true, Namespace: ns}).Find(&list)
	if result.Error != nil {
		return list, result.Error
	}
	return list, nil
}
