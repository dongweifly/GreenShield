package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	core "green_shield/core"
	"strings"
)

//Filter 对多个词库进行过滤
func Filter(request TextReq) (sensitiveWords []string, reasons string, desensitization string) {
	types := strings.Split(request.SceneType, ",")

	desensitization = request.Text

	var matchReasons []string
	for _, sceneType := range types {
		wordsName := request.BizType + "_" + sceneType
		s, d := core.GRuleRouterManager.Get(wordsName).Match(wordsName, desensitization)
		if len(s) > 0 {
			desensitization = d
			sensitiveWords = append(sensitiveWords, s...)
			matchReasons = append(matchReasons, sceneType)
		}
	}

	for _, matchReason := range matchReasons {
		if len(reasons) > 0 {
			reasons += ","
		}
		reasons += matchReason
	}

	return
}

func removeDuplicateElement(src []string) []string {
	dst := make([]string, 0, len(src))
	temp := map[string]struct{}{}
	for _, item := range src {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			dst = append(dst, item)
		}
	}
	return dst
}

func TextMatchHandler(c *gin.Context) {
	//httpserver底层将panic捕获处理掉了，这里不需要再处理了;
	body, _ := c.GetRawData()
	var request TextReq

	if err := json.Unmarshal(body, &request); err != nil {
		log.Warningf("httpServerHandler parse json fail : %s\n", err.Error())

		resp := &TextResp{
			Code: 300,
			Msg:  "Parse parameter error!",
		}
		c.JSON(200, resp)
		log.Warnf("body : %s, resp : %v", string(body), resp)
		return
	}

	if len(request.Text) == 0 {
		resp := &TextResp{
			Code: 300,
			Msg:  "request.Text is empty!",
		}
		c.JSON(200, resp)
		log.Warnf("body : %s, resp : %v", string(body), resp)
		return
	}

	var result Result
	//todo: 写死default_abuse, 后续扩展
	result.Requestid = request.RequestID
	result.Suggestion = "pass"
	result.Sensitivewords, result.Reason, result.Desensitization = Filter(request)
	if len(result.Sensitivewords) > 0 {
		result.Suggestion = "block"
	}
	result.Sensitivewords = removeDuplicateElement(result.Sensitivewords)

	resp := &TextResp{
		Code:   200,
		Msg:    "success",
		Result: &result,
	}

	log.Infof("request: %v, result: %v", request, resp.Result)

	c.JSON(200, resp)

}
