package controller

import (
	"github.com/gin-gonic/gin"
	"green_shield/model"
)

//fieldsQuery 把表中的某一个字段拉出来
func fieldsQuery(c *gin.Context, fieldName string) {
	var fields []string
	db := GMysqlDB
	err := db.Model(&model.WordsInfo{}).Pluck(
		fieldName,
		&fields).Error

	if err != nil {
		c.JSON(200, &CommonResp{
			Code: 300,
			Msg:  err.Error(),
		})
		return
	} else {
		c.JSON(200, &CommonResp{
			Code: 200,
			Msg:  "success",
			Data: &fields,
		})
	}
}

//WordsInfoNames 查询所有的BizType
func WordsInfoNames(c *gin.Context) {
	db := GMysqlDB

	type Data struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}

	var wordsInfo []model.WordsInfo
	err := db.Model(&model.WordsInfo{}).Find(&wordsInfo).Error

	if err != nil {
		c.JSON(200, &CommonResp{
			Code: 300,
			Msg:  err.Error(),
		})
		return
	} else {
		//todo: 把所有的捞出来再赋值，这个实现方式有点蠢，后面看看gorm有没有更优雅的实现方式；
		data := make([]*Data, 0, len(wordsInfo))
		for _, item := range wordsInfo {
			data = append(data, &Data{
				ID:   item.ID,
				Name: item.Name,
			})
		}
		c.JSON(200, &CommonResp{
			Code: 200,
			Data: data,
		})
	}
}

//WordsInfoQuery 查询词库信息，不分页；
func WordsInfoQuery(c *gin.Context) {
	db := GMysqlDB

	var wordsInfo []model.WordsInfo
	err := db.Model(&model.WordsInfo{}).Find(&wordsInfo).Error

	if err != nil {
		c.JSON(200, &CommonResp{
			Code: 300,
			Msg:  err.Error(),
		})
		return
	} else {
		c.JSON(200, &CommonResp{
			Code: 200,
			Data: wordsInfo,
		})
	}
}

func WordsInfoEdit(c *gin.Context) {

}

func WordsInfoAdd(c *gin.Context) {

}
