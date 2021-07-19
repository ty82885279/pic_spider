package logic

import (
	"save2table/dao"
	"save2table/models"
)

func AddPics(p *models.Pics) error {
	return dao.AddPics(p)
}
