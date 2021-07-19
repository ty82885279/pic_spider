package dao

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"save2table/models"
)

func AddPics(p *models.Pics) error {
	tableName := p.Website
	err := tableNotExist(tableName)
	if err != nil {
		//fmt.Printf("表（%s）不存在：%v\n", tableName, err)
		err = creatTable(tableName)
		if err != nil {
			return err
		}
		err = addPic(p)
		return err
	} else {
		//fmt.Printf("表（%s）存在：\n", tableName)
		//插入数据
		err = addPic(p)

	}

	return err
}
func tableNotExist(tableName string) (err error) {

	sqlstr := fmt.Sprintf(`select id from %s limit 0,1`, tableName)
	zap.L().Debug("tableNotExist=>sqlstr", zap.String("sql", sqlstr))

	_, err = DB.Exec(sqlstr)
	if err != nil {

	}
	return
}
func creatTable(tableName string) (err error) {
	sqlStr := `
CREATE TABLE ` + "`" + tableName + "`" + ` (` +
		"`id`" + ` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',` +
		"`title`" + ` varchar(255) NOT NULL COMMENT '标题',` +
		"`description`" + ` text(1000) COMMENT '描述',` +
		"`tags`" + ` varchar(255) COMMENT '标签',` +
		"`urls`" + ` text(1000) NOT NULL COMMENT '链接',` +
		"`cate`" + ` varchar(255) COMMENT '分类',PRIMARY KEY (` + "`id`" + `)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	zap.L().Debug("创建表:creatTable=>sqlStr", zap.String("sql", sqlStr))
	_, err = DB.Exec(sqlStr)

	if err != nil {

		zap.L().Error("创建表失败", zap.String("表名", tableName), zap.Error(err))
		return
	} else {
		fmt.Printf("创建表成功\n")
		zap.L().Debug("创建表成功", zap.String("表名", tableName))
	}
	return
}
func addPic(p *models.Pics) (err error) {
	table := p.Website
	var count int
	sqlstr1 := fmt.Sprintf(`select count(id) from %s where title = ?`, table)
	zap.L().Debug("addPic=>sqlstr1", zap.String("sql", sqlstr1), zap.String("title", p.Title))
	err = DB.Get(&count, sqlstr1, p.Title)
	if err != nil && err != sql.ErrNoRows {

		return
	}
	if count > 0 {
		//更新数据

		return UpdatePic(p)

	} else {
		//添加数据

		return InsertPic(p)

	}
}
func UpdatePic(p *models.Pics) error {

	zap.L().Debug("修改数据", zap.String("title", p.Title))
	sqlStr := fmt.Sprintf(`update %s set title=?,cate=?,description=?,tags=?,urls=? 
where title = ?`, p.Website)
	ret, err := DB.Exec(sqlStr, p.Title, p.Category, p.Description, p.Tags, p.Pics, p.Title)
	if err != nil {
		zap.L().Error("UpdatePic=>DB.Exec", zap.Any("sql", sqlStr), zap.Any("object", p))
		return err
	}
	//var n int64
	_, err = ret.RowsAffected() // 操作影响的行数
	if err != nil {
		zap.L().Error("UpdatePic=>ret.RowsAffected()", zap.Error(err))
		return err
	}
	//fmt.Printf("update success, affected rows:%d\n", n)
	return err
}
func InsertPic(p *models.Pics) (err error) {
	zap.L().Debug("添加数据", zap.String("title", p.Title))
	sqlStr := fmt.Sprintf(`insert into %s(title,cate,description,tags,urls) values (?,?,?,?,?)`, p.Website)
	//fmt.Println(sqlStr)
	_, err = DB.Exec(sqlStr, p.Title, p.Category, p.Description, p.Tags, p.Pics)
	return

}
