package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Url struct {
	Url string `json:"url"`
}

func JoinUrl(urls []Url) []string {
	var uList []string
	for _, v := range urls {
		uList = append(uList, v.Url)
	}
	return uList
}

type Dingtalk struct {
	ID                 int            `json:"id" gorm:"primarykey"`
	CreatedAt          time.Time      `json:"-"`
	UpdatedAt          time.Time      `json:"-"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`
	ClientId           int            `json:"client_id"`
	RobotKeyword       string         `json:"robot_keyword"`
	RobotUrl           string         `json:"robot_url"`
	IsAtAll            bool           `json:"is_at_all"`
	AtMobiles          datatypes.JSON `json:"at_mobiles"`
	AtDingtalkIds      datatypes.JSON `json:"at_dingtalk_ids"`
	RobotUrlList       []Url          `json:"robot_url_list" gorm:"-:all"` //得到前端提交过来的robot_url_list字段
	RobotUrlRandomList []string       `json:"-" gorm:"-:all"`              //根据robot_url_list，得到一个url随机列表
}

func (*Dingtalk) TableName() string {
	return "client_dingtalk"
}
