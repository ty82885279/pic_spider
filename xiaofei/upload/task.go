package upload

import (
	"fmt"
	"os"
	"xiaofei/model"
)

type Task struct {
}

var (
	dataChan = make(chan *model.PicInfo, 100000)
)

func (t *Task) Init() {

	for i := 0; i <= 20; i++ {
		go uploadInfo(dataChan)
	}
}
func uploadInfo(dataChan <-chan *model.PicInfo) {
	for p := range dataChan {
		//取出数据
		//fmt.Println("从通道出取出数据")
		//上传图片
		//fmt.Println("正在上传图片")
		//time.Sleep(time.Second * 2)
		fmt.Printf("%#v\n", p)
		//删除图片
		fmt.Println("图片上传OK:" + p.Dir + "---" + p.Title)
		for i := 0; i < len(p.Path); i++ {
			_ = os.Remove(p.Path[i])
			if i == len(p.Path)-1 {
				err := os.Remove(p.Dir)
				if err != nil {
					fmt.Printf("删除目录出错：%v\n", err)
					panic(err)
				}
				fmt.Println("目录已经删除:" + p.Dir + "---" + p.Title)
			}
		}
	}
}

func SendToChan(p *model.PicInfo) {
	dataChan <- p
}
