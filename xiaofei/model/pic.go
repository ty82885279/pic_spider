package model

type PicInfo struct {
	Url         string   `json:"url"`
	Dir         string   `json:"dir"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Pics        []string `json:"pics"`
	Path        []string `json:"path"`
	Page        string   `json:"page"`
	Index       string   `json:"index"`
}

type ParamsPicInfo struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}
