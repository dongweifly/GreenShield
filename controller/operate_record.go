package controller

import (
	"github.com/gin-gonic/gin"
	"green_shield/model"
	"strconv"
	"time"
)

func addRecord(request *OperateRecord) error {
	return GMysqlDB.Model(&model.OperateRecord{}).Create(
		&model.OperateRecord{
			Text:        request.Text,
			Operator:    request.Operator,
			BizType:     request.BizType,
			SceneType:   request.SceneType,
			OperateType: request.OperateType,
			CreateTime:  time.Now().Unix(),
		}).Error
}

// 查询表记录

func SearchRecordHandler(c *gin.Context) {
	//todo: 重复代码合并
	var request OperateRecord
	if err := parsePostParams(c, &request); err != nil {
		return
	}

	//如果需要分页，设置分页条件；
	db := GMysqlDB

	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))

	if pageNum <= 1 {
		pageNum = 1
	}

	if pageSize <= 0 || pageSize > 100 {
		pageSize = 100
	}
	var count = int64(0)
	queryCond := model.OperateRecord{
		Text:        request.Text,
		BizType:     request.BizType,
		SceneType:   request.SceneType,
		Operator:    request.Operator,
		OperateType: request.OperateType,
	}

	db.Model(&model.OperateRecord{}).Where(&queryCond).Count(&count)

	if pageNum > 0 && pageSize > 0 {
		db = db.Limit(pageSize).Offset((pageNum - 1) * pageSize)
	}

	var records []model.OperateRecord

	//这个地方是否需要重复查询有待验证
	result := db.Model(&model.OperateRecord{}).Where(&queryCond).Find(&records)

	if result.Error != nil {
		c.JSON(300, result.Error.Error())
		return
	}

	c.JSON(200, &PageDataResp{
		Code:      200,
		PageNo:    pageNum,
		PageSize:  pageSize,
		TotalSize: int(count),
		Data:      records,
	})
}

func AddRecordHandler(c *gin.Context) {
	var request OperateRecord

	if err := parsePostParams(c, &request); err != nil {
		return
	}

	err := addRecord(&request)

	if err != nil {
		c.JSON(200, CommonResp{
			Code: 300,
			Msg:  err.Error(),
		})
	}

	c.JSON(200, CommonResp{
		Code: 200,
		Msg:  "success",
	})

}
