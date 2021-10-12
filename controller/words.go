package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
	"green_shield/model"
	"strconv"
	"strings"
	"time"
)

//UpdateWordInfoTime 更新word_info表的最后更新时间，自动加载；
func UpdateWordInfoTime(bizType, sceneType string) error {
	assigns := map[string]interface{}{model.WordsColumns.UpdateTime: time.Now().Unix()}

	return GMysqlDB.Model(&model.WordsInfo{}).
		Where("bizType=? and SceneType=?", bizType, sceneType).
		Updates(assigns).Error
}

//FuzzySearchWordsHandler 模糊查询
func FuzzySearchWordsHandler(c *gin.Context) {
	var request FuzzyWordQueryReq
	if err := parsePostParams(c, &request); err != nil {
		return
	}
	db := GMysqlDB
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 100
	}

	var words []model.Words

	var count = int64(0)
	db.Model(&model.Words{}).
		Where("text like '%" + request.Text + "%'").Count(&count)

	if int(count) > pageSize {
		db = db.Limit(pageSize).Offset((pageNum - 1) * pageSize)
	}
	result := db.Model(&model.Words{}).
		Where("text like '%" + request.Text + "%'").Find(&words)

	if result.Error != nil {
		c.JSON(300, result.Error.Error())
		return
	}

	c.JSON(200, &PageDataResp{
		Code:      200,
		Msg:       "success",
		PageNo:    pageNum,
		PageSize:  pageSize,
		TotalSize: int(count),
		Data:      words,
	})
}

//SearchWordsHandler 精确查找
func SearchWordsHandler(c *gin.Context) {
	var request SearchWords
	if err := parsePostParams(c, &request); err != nil {
		return
	}
	//如果需要分页，设置分页条件；
	db := GMysqlDB

	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 100
	}
	var count = int64(0)

	db.Model(&model.Words{}).Where(&model.Words{
		Text:      request.Text,
		BizType:   request.BizType,
		SceneType: request.SceneType,
	}).Count(&count)

	if int(count) > pageSize {
		db = db.Limit(pageSize).Offset((pageNum - 1) * pageSize)
	}

	var words []model.Words

	result := db.Where(&model.Words{
		Text:      request.Text,
		BizType:   request.BizType,
		SceneType: request.SceneType,
	}).Find(&words)

	if result.Error != nil {
		c.JSON(300, result.Error.Error())
		return
	}

	c.JSON(200, &PageDataResp{
		Code:      200,
		Msg:       "success",
		PageNo:    pageNum,
		PageSize:  pageSize,
		TotalSize: int(count),
		Data:      words,
	})

}

//ExportAllWords 导出所有敏感词
func ExportAllWords(c *gin.Context) {
	db := GMysqlDB
	db = db.Limit(10000)

	var words []model.Words

	result := db.Where(&model.Words{}).Find(&words)

	if result.Error != nil {
		c.JSON(300, result.Error.Error())
		return
	}

	c.JSON(200, &CommonResp{
		Code: 200,
		Msg:  "success",
		Data: words,
	})
}

//ImportWordsHandler 词库导入；
func ImportWordsHandler(c *gin.Context) {
	var request ImportWordsReq
	if err := parsePostParams(c, &request); err != nil {
		return
	}

	var removeWords []*model.Words
	var addWords []*model.Words

	now := time.Now().Unix()
	for _, w := range request.Words {
		word := &model.Words{
			Text:       w.Text,
			BizType:    w.BizType,
			SceneType:  w.SceneType,
			Valid:      true,
			CreateTime: now,
		}
		if w.OperateType == "remove" {
			removeWords = append(removeWords, word)
		} else if w.OperateType == "add" {
			addWords = append(addWords, word)
		} else {
			c.JSON(200, &CommonResp{
				Code: 300,
				Msg:  "operateType not exist or invalid",
			})
			return
		}
	}

	if len(addWords) > 0 {
		if err := insertBatchWords(addWords); err != nil {
			c.JSON(200, &CommonResp{
				Code: 300,
				Msg:  err.Error(),
			})
			return
		}
	}

	if len(removeWords) > 0 {
		if err := removeBatchWords(removeWords); err != nil {
			c.JSON(200, &CommonResp{
				Code: 300,
				Msg:  err.Error(),
			})
			return
		}
	}

	go func() {
		if len(addWords) > 0 {
			for _, w := range addWords {
				_ = addRecord(&OperateRecord{
					Text:        w.Text,
					SceneType:   w.SceneType,
					Operator:    request.Operator,
					OperateType: "add",
				})
			}
		}

		if len(removeWords) > 0 {
			for _, w := range removeWords {
				_ = addRecord(&OperateRecord{
					Text:        w.Text,
					SceneType:   w.SceneType,
					Operator:    request.Operator,
					OperateType: "remove",
				})
			}
		}
	}()

	c.JSON(200, &CommonResp{
		Code: 200,
		Msg:  "Success",
	})
}

//AddWordsHandler 前端接口，批量添加
func AddWordsHandler(c *gin.Context) {
	var request AddWordsReq

	if err := parsePostParams(c, &request); err != nil {
		return
	}

	var wordInfo model.WordsInfo
	if err := GMysqlDB.Model(&model.WordsInfo{}).Where(
		"name", request.WordsName).First(&wordInfo).Error; err != nil {
		c.JSON(200, &CommonResp{
			Code: 300,
			Msg:  err.Error(),
		})
		return
	}

	now := time.Now().Unix()
	words := make([]*model.Words, 0, len(request.Words))
	for _, w := range request.Words {
		words = append(words, &model.Words{
			BizType:    wordInfo.BizType,
			SceneType:  wordInfo.SceneType,
			Text:       w,
			UpdateTime: now,
			CreateTime: now,
		})
	}

	if err := insertBatchWords(words); err != nil {
		c.JSON(200, &CommonResp{
			Code: 300,
			Msg:  err.Error(),
		})
	}

	c.JSON(200, &CommonResp{
		Code: 200,
		Msg:  "Success",
	})

	go func() {
		if err := UpdateWordInfoTime(wordInfo.BizType, wordInfo.SceneType); err != nil {
			log.Warningf("Update wordInfo error :%s", err.Error())
		}
		_ = addRecord(&OperateRecord{
			Text:        strings.Join(request.Words, " "),
			BizType:     wordInfo.BizType,
			SceneType:   wordInfo.SceneType,
			Operator:    "HaiTang",
			OperateType: "add",
		})
	}()
}

//RemoveWordByIdHandler 根据某一个ID删除敏感词
func RemoveWordByIdHandler(c *gin.Context) {
	var request dDeleteOneWordReq
	if err := parsePostParams(c, &request); err != nil {
		return
	}

	var word model.Words
	res := model.Words{
		ID: request.ID,
	}

	result := GMysqlDB.Model(model.Words{}).Where(&res).First(&word)
	if result.Error != nil {
		c.JSON(200, &CommonResp{
			Code: 300,
			Msg:  result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == int64(0) {
		c.JSON(200, &CommonResp{
			Code: 300,
			Msg:  "Can not find this word ID",
		})
		return
	}

	err := GMysqlDB.Model(model.Words{}).Delete(&res).Error
	log.Infof("###### Delete res %v", res)

	if err != nil {
		c.JSON(200, &CommonResp{
			Code: 300,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, &CommonResp{
		Code: 200,
		Msg:  "Success",
	})

	go func() {
		if err := UpdateWordInfoTime(word.BizType, word.SceneType); err != nil {
			log.Warningf("Update wordInfo error :%s", err.Error())
		}
		_ = addRecord(&OperateRecord{
			Text:        word.Text,
			BizType:     word.BizType,
			SceneType:   word.SceneType,
			Operator:    "HaiTang",
			OperateType: "remove",
		})
	}()

}

//ModifyWordsHandler 批量增减敏感词；前端已经不再使用
func ModifyWordsHandler(c *gin.Context) {
	var request ModifyWords
	if err := parsePostParams(c, &request); err != nil {
		return
	}

	if len(request.BizType) == 0 && len(request.SceneType) == 0 {
		c.JSON(200, &CommonResp{
			Code: 300,
			Msg:  "param error, bizType and  sceneType can not empty",
		})
		return
	}

	//先查一下BizType和SceneType在
	rowAffect := GMysqlDB.Model(&model.WordsInfo{}).Where("bizType=? and SceneType=?",
		request.BizType, request.SceneType).Find(&model.WordsInfo{}).RowsAffected

	if rowAffect <= 0 {
		c.JSON(200, &CommonResp{
			Code: 300,
			Msg:  fmt.Sprintf("bizeType:%s seceneType:%s not exist", request.BizType, request.SceneType),
		})
		return
	}

	for _, text := range request.AddTexts {
		now := time.Now().Unix()
		word := &model.Words{
			Text:       text,
			BizType:    request.BizType,
			SceneType:  request.SceneType,
			Valid:      true,
			UpdateTime: now,
			CreateTime: now,
		}
		err := insertOneWords(word)

		if err != nil {
			c.JSON(200, &CommonResp{
				Code: 300,
				Msg:  err.Error(),
			})
			return
		}
		log.Infof("Add word %v", word)
	}

	for _, text := range request.RemoveTexts {
		err := GMysqlDB.Where("text=? and bizType=? and SceneType=?",
			text, request.BizType, request.SceneType).Delete(&model.Words{}).Error
		if err != nil {
			c.JSON(200, &CommonResp{
				Code: 300,
				Msg:  err.Error(),
			})
			return
		}
		log.Infof("Rmove word %s %s %s", text, request.BizType, request.SceneType)
	}

	//如果操作成功，同时更新wordInfo的最后更新时间；
	if err := UpdateWordInfoTime(request.BizType, request.SceneType); err != nil {
		c.JSON(200, &CommonResp{
			Code: 300,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, &CommonResp{
		Code: 200,
		Msg:  "Success",
	})

	go func() {
		if len(request.AddTexts) > 0 {
			_ = addRecord(&OperateRecord{
				Text:        strings.Join(request.AddTexts, " "),
				SceneType:   request.SceneType,
				Operator:    request.Operator,
				OperateType: "add",
			})
		}

		if len(request.RemoveTexts) > 0 {
			_ = addRecord(&OperateRecord{
				Text:        strings.Join(request.RemoveTexts, " "),
				SceneType:   request.SceneType,
				Operator:    request.Operator,
				OperateType: "remove",
			})
		}
	}()
	return
}

//insertOneWords
func insertOneWords(word *model.Words) error {
	now := time.Now().Unix()
	return GMysqlDB.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: model.WordsColumns.Text},
			{Name: model.WordsColumns.BizType},
			{Name: model.WordsColumns.SceneType},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			model.WordsColumns.Valid:      true,
			model.WordsColumns.UpdateTime: now}),
	}).Create(&word).Error
}

//insertBatchWords
func insertBatchWords(words []*model.Words) error {
	now := time.Now().Unix()
	return GMysqlDB.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: model.WordsColumns.Text},
			{Name: model.WordsColumns.BizType},
			{Name: model.WordsColumns.SceneType},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			model.WordsColumns.Valid:      true,
			model.WordsColumns.UpdateTime: now}),
	}).CreateInBatches(&words, len(words)).Error
}

//removeBatchWords 批量删除词库
func removeBatchWords(words []*model.Words) error {
	for _, word := range words {
		err := GMysqlDB.Where("text=? and bizType=? and SceneType=?",
			word.Text, word.BizType, word.SceneType).Delete(&model.Words{}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

//parsePostParams 解析post Body接口
func parsePostParams(c *gin.Context, postParams interface{}) error {
	body, _ := c.GetRawData()
	if err := json.Unmarshal(body, postParams); err != nil {
		log.Warningf("httpServerHandler parse json fail : %s\n %s", err.Error(), string(body))

		c.JSON(200, &CommonResp{
			Code: 300,
			Msg:  "Parse parameter error!",
		})
		return err
	}
	return nil
}
