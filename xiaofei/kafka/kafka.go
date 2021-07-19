package kafka

import (
	"encoding/json"
	"fmt"
	"xiaofei/model"
	"xiaofei/upload"

	"github.com/Shopify/sarama"
)

var (
	consumer sarama.Consumer
)

func ReadMsgFromKafka(addr string, topic string) {
	var err error
	consumer, err = sarama.NewConsumer([]string{addr}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println(partitionList)
	fmt.Println("长度--", len(partitionList))

	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		//defer pc.AsyncClose() //不能关闭
		// 异步从每个分区消费信息
		//for i := 0; i <= 10; i++ {
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {

				p := new(model.PicInfo)
				_ = json.Unmarshal(msg.Value, p)

				fmt.Println("取出数据")
				//发送数据，立即返回
				upload.SendToChan(p)
				fmt.Println("数据放置OK")

				//fmt.Printf("%#v\n", p)

				//todo:上传图片
				//get := http.Get("http://127.0.0.1:8088/index", nil, nil)
				//fmt.Println("返回的数据:" + get)
				//paramsMap := make(map[string]string, 10)
				//paramsMap["title"] = "美女图片"
				//paramsMap["des"] = "美女描述美女描述"
				//paramsMap["tags"] = `["青春","阳光","可爱"]`
				//files := make([]http.UploadFile, 0, 10)
				//file1 := http.UploadFile{
				//	Name:     "name",
				//	Filepath: "/Users/mrlee/Desktop/img/车桌面/xinggan/42230/d18831ba8a387ab38619e7d78d5cb835.jpg",
				//}
				//file2 := http.UploadFile{
				//	Name:     "name",
				//	Filepath: "/Users/mrlee/Desktop/img/车桌面/xinggan/42193/82a5c8d17eb813c92834ec281f20691a.jpg",
				//}
				//
				//files = append(files, file1)
				//files = append(files, file2)
				//res := http.PostFile("http://127.0.0.1:8088/pic", paramsMap, files, nil)
				//fmt.Println("正在上传图片....")
				//fmt.Println("图片上传成功....")
				//删除图片
				//for i := 0; i < len(p.Path); i++ {
				//	time.Sleep(time.Second * 1)
				//	_ = os.Remove(p.Path[i])
				//	if i == len(p.Path)-1 {
				//		os.Remove(p.Dir)
				//		fmt.Println("图片目录删除成功")
				//	}
				//
				//}

			}
		}(pc)
		//}

	}
}
