package main

import (
	"fmt"
	"xiaofei/kafka"
	"xiaofei/upload"

	"github.com/spf13/viper"
)

type Conf struct {
	Topic        []string `mapstructure:"topic"`
	KafkaAddress string   `mapstructure:"address"`
}

var cfg = new(Conf)

func main() {

	//从kafka中消费
	viper.SetConfigFile("./cfg/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("read congfig failed:%v\n", err)
		return
	}
	err = viper.Unmarshal(cfg)
	if err != nil {
		fmt.Printf("viper unmarshal failed,err:%v\n", err)
		return
	}

	//

	//根据topic创建消费者
	for i := 0; i < len(cfg.Topic); i++ {
		fmt.Printf("%#v\n", cfg.Topic[i])
		task := new(upload.Task)
		task.Init()
		go kafka.ReadMsgFromKafka(cfg.KafkaAddress, cfg.Topic[i])
	}
	select {}
}
