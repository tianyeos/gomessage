package vClient

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gomessage/pkg/models"
	"gomessage/pkg/utils"
	"net/http"
)

// PostClient
// @Tags Client
// @Summary 新增一个客户端
// @Router /api/v1/:namespace/client [POST]
func PostClient(g *gin.Context) {

	body := models.Client{}
	if err := g.ShouldBindJSON(&body); err != nil {
		return
	}

	//获取url中的namespace
	body.Namespace = g.Param("namespace")

	//判断客户端类型（绑定客户端信息到client的Extend延伸信息中）
	switch body.ClientType {
	case utils.VarDingtalk:
		json.Unmarshal(body.ClientInfo, &body.ExtendDingtalk)

	case utils.VarWechatRobot:
		json.Unmarshal(body.ClientInfo, &body.ExtendWechatRobot)

	case utils.VarFeishu:
		json.Unmarshal(body.ClientInfo, &body.ExtendFeishu)

	case utils.VarWechatApplication:
		json.Unmarshal(body.ClientInfo, &body.ExtendWechatApplication)

	default:
		return
	}

	result, _ := models.AddClient(&body)

	g.JSON(http.StatusOK, utils.ResponseSuccessful("创建成功", result))
}
