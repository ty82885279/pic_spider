package kafka

import (
	cfg2 "aso/cfg"
	"aso/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	res       string //返回响应字符串
	cfg       *cfg2.Conf
	EndCh     chan struct{}
	dataSlice = make([]*model.PicUrls, 0, 1000)
)

func VisitWeb(c *cfg2.Conf) {

	cfg = c
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
		//chromedp.Flag("hide-scrollbars", false),
		//chromedp.Flag("mute-audio", false),
		//chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	//创建chrome窗口
	allocCtx, cancel2 := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel2()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	//
	//ctx, cancel = chromedp.NewContext(
	//	context.Background(),
	//)
	//defer cancel()

	//fmt.Sprintf("%s/%s.html", cfg.StartUrl, cfg.Page)

	for {

		EndCh = make(chan struct{})

		err := chromedp.Run(ctx, targetUrl(cfg.StartUrl))
		if err != nil {
			log.Fatal(err)
		}

		_ = chromedp.Run(ctx,
			chromedp.OuterHTML(`body`, &res),
		)

		dom, err := goquery.NewDocumentFromReader(strings.NewReader(res))
		if err != nil {
			log.Fatalln(err)
		} //.

		dataSlice = make([]*model.PicUrls, 0, 1000)
		pp := strconv.Itoa(cfg.Page)
		dom.Find(".egene_spe_pic_li").Each(func(i int, selection *goquery.Selection) {

			//获取图片详情链接
			url, _ := selection.Find("dt > div:nth-child(1) > a").Attr("href")
			//将数据存入数组
			d := i + 1
			index := strconv.Itoa(d)
			//fmt.Println(index)
			Pic := &model.PicUrls{
				Url:   url,
				Page:  pp,
				Index: index,
			}
			dataSlice = append(dataSlice, Pic)
		})
		//fmt.Printf("数组====%#v\n", dataSlice)
		//
		//defer func() {
		//	closeErr := chromedp.Cancel(ctx)
		//	if closeErr != nil {
		//		panic(closeErr)
		//	}
		//}()

		total := len(dataSlice) - cfg.Index
		//total := 5
		go SendToKaka(total)
		//
		tasksCh := make(chan *model.PicUrls, 100)
		for i := cfg.Index; i < len(dataSlice); i++ {
			tasksCh <- dataSlice[i]
		}
		close(tasksCh)

		for i := 1; i <= cfg.MaxWork; i++ {
			go Woker(tasksCh)
		}
		<-EndCh

		fmt.Printf("----------第 %d 页下载成功----------\n", cfg.Page)

		cfg.Page++

	}
}
func targetUrl(url string) chromedp.Tasks {
	res = "" //每次翻页清空响应
	return chromedp.Tasks{
		chromedp.Navigate(url),
		//chromedp.Sleep(10 * time.Second),
		chromedp.WaitVisible(`#auto_width_specialist_0`, chromedp.ByQuery),
		//chromedp.OuterHTML(`body`, &res),
		chromedp.OuterHTML("body", &res, chromedp.ByQuery),
		//chromedp.Evaluate(`document.getElementsByClassName("js-table-content")[0].clientHeight`, &height),

	}
}

func Woker(tasks <-chan *model.PicUrls) {

	for p := range tasks {

		ctx1, _ := chromedp.NewContext(
			context.Background(),
		)

		htmlStr, dirPath := task(ctx1, p.Url)
		//fmt.Println("返回链接----" + dirPath)

		dom1, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
		if err != nil {
			log.Fatalln(err)
		}

		//获取标签
		var pic = new(model.PicInfo)
		pic.Dir = dirPath
		pic.Page = p.Page
		pic.Index = p.Index

		dom1.Find("div:nth-child(13) > div.fleft.arc_pic > div:nth-child(10) > div.myarc_tag >a").Each(func(i int, selection *goquery.Selection) {
			tag := selection.Text()
			pic.Tags = append(pic.Tags, tag)
		})
		//获取图集
		dom1.Find("div:nth-child(13) > div.fleft.arc_pic > div.arc_pandn > div > div.swiper-wrapper > div").Each(func(i int, selection *goquery.Selection) {
			p, _ := selection.Find("a").Attr("src")
			pic.Pics = append(pic.Pics, p)
			//
		})
		//获取标题
		title := dom1.Find("div:nth-child(13) > div.fleft.arc_pic > div.arc_top > h1").Text()
		//获取简介
		description := dom1.Find("div:nth-child(13) > div.fleft.arc_pic > div.myarc_intro").Text()
		description = strings.Replace(description, "内容简介 ", "", 1)
		description = strings.TrimSpace(description)

		//fmt.Println(title)
		pic.Url = p.Url
		pic.Title = title
		pic.Description = description
		pic.Path = make([]string, 0, len(pic.Pics))

		for i := 0; i < len(pic.Pics); i++ {
			_, err, pathStr := getImg(pic.Pics[i], dirPath)
			if err != nil {
				fmt.Println("保存图片错误" + pathStr)
			}
			pic.Path = append(pic.Path, pathStr)
		}
		pic.Pics = make([]string, 0, 0)
		data, err := json.Marshal(pic)
		if err != nil {
			fmt.Printf("json marshal err:%v\n", err)
			return
		}
		//fmt.Printf("%#v\n", pic)
		//关闭chrome,等待资源被回收
		closeErr := chromedp.Cancel(ctx1)
		if closeErr != nil {
			fmt.Printf("关闭chrome错误:%v\n", closeErr)
		}
		//fmt.Println("发送数据-----" + string(data))
		//发送数据到kafka channel
		fmt.Printf("第 %s 页的第 %s 条数据---图集目录:%s\n", pic.Page, pic.Index, pic.Dir)
		SendChan(cfg.Category, string(data))

	}
}

// 保存图片
func getImg(url, dst string) (n int64, err error, pathStr string) {
	//fmt.Println("图片名称")
	path := strings.Split(url, "/")
	var name string
	if len(path) > 1 {
		name = path[len(path)-1]
	}
	//fmt.Println("目录：" + dst)
	//fmt.Println("文件名：" + name)
	out, err := os.Create(dst + name)

	defer func() {
		err = out.Close()
		if err != nil {
			panic(err)
		}
	}()

	resp, err := http.Get(url)

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			return
		}
	}()
	pix, err := ioutil.ReadAll(resp.Body)
	n, err = io.Copy(out, bytes.NewReader(pix))

	//fmt.Println(out)
	pathStr = fmt.Sprintf("%s%s", dst, name)
	//fmt.Println("path:" + pathStr)
	return n, err, pathStr
}

// 每个图集任务
func task(ctx context.Context, url string) (string, string) {
	//
	path := strings.Split(url, "/")
	var name string
	if len(path) > 1 {
		name = path[len(path)-1]
	}
	nameSlice := []byte(name)
	//fmt.Println(nameSlice)
	nameSlice = nameSlice[0 : len(name)-5]
	//fmt.Println(nameSlice)
	dirName := string(nameSlice)
	//fmt.Println("目录名字：" + dirName)
	//fmt.Println("-------")
	//fmt.Printf("配置:%#v\n", cfg)
	dirPath := fmt.Sprintf(cfg.ImgPath + "/" + cfg.Category + "/" + dirName + "/")
	//fmt.Println("-------：" + dirPath)
	_ = os.Mkdir(cfg.ImgPath+"/"+cfg.Category+"/", os.ModePerm)
	_ = os.Mkdir(dirPath, os.ModePerm)
	//
	var res string
	_ = chromedp.Run(ctx, downloadUrl(url, &res))
	return res, dirPath
}
func downloadUrl(url string, res *string) chromedp.Tasks {

	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(5 * time.Second),
		//chromedp.WaitVisible("body > div:nth-child(12) > div.fleft.arc_pic > div.arc_pandn > div > div.swiper-wrapper", chromedp.ByQuery),
		//chromedp.OuterHTML(`body`, &res),
		chromedp.OuterHTML("body", res, chromedp.ByQuery),
		//chromedp.Evaluate(`document.getElementsByClassName("js-table-content")[0].clientHeight`, &height),
	}
}
