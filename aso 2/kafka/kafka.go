package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

type picData struct {
	topic string
	data  string
}

var (
	producer sarama.SyncProducer //消费者
	dataChan chan *picData       //保存每页图集的Chan
)

// 初始化kafka
func Init(address []string) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	producer, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		return
	}
	dataChan = make(chan *picData, 100000)
	//go SendToKaka()
	return
}

// 向通道内发送数据
func SendChan(topic, data string) {
	msg := &picData{
		topic: topic,
		data:  data,
	}
	dataChan <- msg
}

// 从通道内取数据
func SendToKaka(total int) {
	//fmt.Println(total)
	//loop:
	cnt := 0
	for pic := range dataChan {

		//构建消息
		msg := &sarama.ProducerMessage{}
		msg.Topic = pic.topic
		msg.Value = sarama.StringEncoder(pic.data)
		//发送消息
		_, _, err := producer.SendMessage(msg)
		if err != nil {
			fmt.Printf("kafka send msg err:%v\n", err)
			panic(err)
		}
		//执行结束
		cnt++
		if cnt == total {

			EndCh <- struct{}{}
			break
		}
	}

}
