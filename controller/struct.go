package controller

type CommonResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type PageDataResp struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	PageNo    int         `json:"pageNo,omitempty"`
	PageSize  int         `json:"pageSize,omitempty"`
	TotalSize int         `json:"TotalSize,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type OperateRecord struct {
	Text        string `json:"text"`      // 敏感词，最大128
	BizType     string `json:"bizType"`   // 业务类型
	SceneType   string `json:"sceneType"` // 场景编码，porn, politics, ad
	Operator    string `json:"operator"`
	OperateType string `json:"operateType"`
}

type ModifyWords struct {
	Name        string   `jons:"name"`
	BizType     string   `json:"bizType"`
	SceneType   string   `json:"sceneType"`
	AddTexts    []string `json:"addTexts"`
	RemoveTexts []string `json:"removeTexts"`
	Operator    string   `json:"operator"`
}

type ModifyWord struct {
	Name        string `jons:"name"`
	BizType     string `json:"bizType"`
	SceneType   string `json:"sceneType"`
	Text        string `json:"text"`
	OperateType string `json:"operateType"` //remove or add
}

type ImportWordsReq struct {
	Words    []ModifyWord `json:"words"`
	Operator string           `json:"operator"`
}

type dDeleteOneWordReq struct {
	ID int64 `json:"id"`
}

type AddWordsReq struct {
	WordsName string   `json:"wordsName"` //wordInfo中词库的名字
	Words     []string `json:"words"`
}

type FuzzyWordQueryReq struct {
	Text string `json:"text"`
}

type SearchWords struct {
	Text      string `json:"text"`
	BizType   string `json:"bizType"`
	SceneType string `json:"sceneType"`
}

type TextReq struct {
	RequestID string `json:"requestId"`
	Text      string `json:"text"`
	BizType   string `json:"bizType"`
	SceneType string `json:"sceneType"`
	Extent    string `json:"extent"`
}

type TextResp struct {
	Code   int     `json:"code"`
	Msg    string  `json:"msg"`
	Result *Result `json:"result"`
}
type Result struct {
	Requestid       string   `json:"requestId"`
	Suggestion      string   `json:"suggestion"`
	Reason          string   `json:"reason,omitempty"`
	Sensitivewords  []string `json:"sensitiveWords,omitempty"`
	Desensitization string   `json:"desensitization,omitempty"`
}
