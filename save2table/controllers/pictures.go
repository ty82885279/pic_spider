package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"save2table/logic"
	"save2table/models"
)

// PicsHandler 上传图集接口
// @Summary 上传图集接口
// @Description 上传美女图集接口
// @Tags 上传
// @Accept application/json
// @Produce application/json
// @Param object body models.Pics false "参数列表"
// @Success 200 {object} _ResponseSuccess "请求成功"
// @Failure 1001 {object} _Responsefailure "参数错误"
// @Failure 1002 {object} _Responsefailure "服务繁忙"
// @Router /pics [post]
func PicsHandler(c *gin.Context) {
	//解析参数
	p := new(models.Pics)
	err := c.ShouldBind(p)
	if err != nil {
		zap.L().Error("参数错误", zap.Any("参数列表", p))
		c.JSON(1001, gin.H{
			"code": 1001,
			"msg":  "参数错误",
		})
		return
	}
	//业务逻辑
	err = logic.AddPics(p)
	if err != nil {
		zap.L().Error("添加数据错误", zap.Error(err))
		c.JSON(1002, gin.H{
			"code": 1002,
			"msg":  "服务器繁忙",
		})

		return
	}
	//返回响应
	c.JSON(200, gin.H{
		"code": 2000,
		"msg":  "添加成功",
		"data": p,
	})
}
