package vClient

import (
	"github.com/gin-gonic/gin"
	"gomessage/pkg/models"
	"gomessage/pkg/utils"
	"net/http"
)

// ListClient
// @Tags Client
// @Summary 获取所有客户端
// @Router /api/v1/:namespace/client [GET]
func ListClient(g *gin.Context) {
	ns := g.Param("namespace")
	result, _ := models.ListClient(ns)
	g.JSON(http.StatusOK, utils.ResponseSuccessful("获取成功", result))
}
