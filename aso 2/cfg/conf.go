package cfg

type Conf struct {
	KafkaAddr string `mapstructure:"address"`       // kafakad地址
	StartUrl  string `mapstructure:"starturl"`      // 起始url
	Category  string `mapstructure:"category"`      // 分类
	ImgPath   string `mapstructure:"img_path"`      // 图片绝对路径
	MaxWork   int    `mapstructure:"max_work_size"` //并发数量
	Page      int    `mapstructure:"page"`          //起始页数
	Index     int    `mapstructure:"index"`         //起始下标
}
