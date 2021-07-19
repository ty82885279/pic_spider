package main

import (
	cfg2 "aso/cfg"
	"aso/kafka"
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var (
	wg = new(sync.WaitGroup)

	cfg = new(cfg2.Conf)
)

func main() {

	//加载加载配置
	viper.SetConfigFile("./cfg/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("配置文件错误:%v\n", err)

	}
	err = viper.Unmarshal(cfg)
	if err != nil {
		fmt.Println("配置序列化错误:%v\n", err)
		return
	}
	fmt.Printf("配置成功:%#v\n", cfg)
	//kafka.Init
	err = kafka.Init([]string{cfg.KafkaAddr})
	if err != nil {
		fmt.Printf("kafka init err:%v\n", err)
		return
	}
	fmt.Println("kafka")

	kafka.VisitWeb(cfg)

	select {}

}
