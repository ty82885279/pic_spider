package models

type Pics struct {
	Website     string `json:"webside" form:"webside" example:"baidu"`                //必填：网站拼音或者缩写
	Category    string `json:"category" form:"category" example:"清纯" `                //选填：分类，没有就传空
	Title       string `json:"title" form:"title" example:"清纯美少女"`                    // 必填：图集标题
	Description string `json:"description" form:"description" example:"这里是描述"`        //选填：图集描述，没有就传空
	Tags        string `json:"tags" form:"tags" example:"标签1|标签2|标签3"`                // 选填：标签，没有就传空，存在1个以上请用'｜'拼接，
	Pics        string `json:"pics" form:"pics" example:"xxx1.jpg|xxx2.jpg|xxx3.jpg"` //必填：图片链接，存在1个以上请用'｜'拼接
}
