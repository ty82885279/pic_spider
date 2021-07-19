package main

import (
	"fmt"
	"go.uber.org/zap"
	"save2table/logger"

	"save2table/dao"
	"save2table/routers"
	"save2table/settings"
)

// @title 美女图片
// @version 1.0
// @description 上传各种美女图片

// @Host 192.168.3.180:8081
// @BasePath /api/v1
func main() {

	err := settings.Init()
	if err != nil {
		fmt.Printf("加载配置错误:%v\n", err)
		return
	}
	//2. 初始化日志
	err = logger.Init(settings.Cfg.LoggerConf, settings.Cfg.Mode)
	fmt.Printf("日志配置：%v\n", settings.Cfg.LoggerConf)
	if err != nil {
		zap.L().Error("Init log failed", zap.Error(err))
		return
	}
	defer zap.L().Sync() //将缓冲区的日志追加到文件中
	err = dao.Init(settings.Cfg.MysqlConf)
	if err != nil {
		zap.L().Error("Init mysql failed", zap.Error(err))
		return
	}
	r := routers.SetUp()
	r.Run(fmt.Sprintf(":%d", settings.Cfg.Port))
}
