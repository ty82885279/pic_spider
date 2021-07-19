package controllers

import "save2table/models"

type _ResponseSuccess struct {
	Code    int64        `json:"code"` // 状态码
	Message string       `json:"msg"`  // 提示信息
	Data    *models.Pics `json:"data"` // 数据
}
type _Responsefailure struct {
	Code    int64  `json:"code"` // 状态码
	Message string `json:"msg"`  // 提示信息

}
