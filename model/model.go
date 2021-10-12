package model

// Words [...]
type Words struct {
	ID         int64  `json:"id" gorm:"primaryKey;column:id;type:bigint(20);not null"`                                                 // 主键
	Text       string `json:"text" gorm:"uniqueIndex:idx_bizType_text;column:text;type:varchar(128)"`                                  // 敏感词，最大128
	BizType    string `json:"bizType" gorm:"uniqueIndex:idx_bizType_text;index:idx_bizType_sceneType;column:bizType;type:varchar(32)"` // 业务类型
	SceneType  string `json:"sceneType" gorm:"index:idx_bizType_sceneType;column:sceneType;type:varchar(32);not null"`                 // 场景编码，porn, politics, ad
	Extent     string `json:"extent" gorm:"column:extent;type:varchar(255)"`                                                           // 扩展
	Valid      bool   `json:"valid" gorm:"column:valid;type:tinyint(1);default:1"`                                                     // 是否有效，true:有效， false:无效
	UpdateTime int64  `json:"updateTime" gorm:"column:updateTime;type:bigint(20)"`                                                     // 创建时间戳
	CreateTime int64  `json:"createTime" gorm:"column:createTime;type:bigint(20)"`                                                     // 创建时间戳
}

// WordsInfo [...]
type WordsInfo struct {
	ID         int64  `json:"id" gorm:"primaryKey;column:id;type:bigint(20);not null"`                                       // 主键
	Name       string `json:"name" gorm:"unique;column:name;type:varchar(64);not null"`                                      // 词库名称
	BizType    string `json:"bizType" gorm:"uniqueIndex:idx_biztype_scenetype;column:bizType;type:varchar(32);not null"`     // 业务类型
	SceneType  string `json:"sceneType" gorm:"uniqueIndex:idx_biztype_scenetype;column:sceneType;type:varchar(32);not null"` // 场景编码，porn, politics, ad
	ParentName string `json:"parentName" gorm:"column:parentName;type:varchar(64)"`
	UpdateTime int64  `json:"updateTime" gorm:"column:updateTime;type:bigint(20)"` // 更新时间戳
	CreateTime int64  `json:"createTime" gorm:"column:createTime;type:bigint(20)"` // 创建时间戳
}

// OperateRecord [...]
type OperateRecord struct {
	ID          int64  `json:"id" gorm:"primaryKey;column:id;type:bigint(20);not null"`              // 主键
	Text        string `json:"text" gorm:"index:idx_text;column:text;type:varchar(128)"`             // 敏感词，最大128
	BizType     string `json:"bizType" gorm:"column:bizType;type:varchar(32)"`                       // 业务类型
	SceneType   string `json:"sceneType" gorm:"column:sceneType;type:varchar(32);not null"`          // 场景编码，porn, politics, ad
	Operator    string `json:"operator" gorm:"index:idx_operator;column:operator;type:varchar(255)"` // 扩展
	OperateType string `json:"operateType" gorm:"column:operateType;type:varchar(32);not null"`      // 操作类型, remove, add
	CreateTime  int64  `json:"createTime" gorm:"column:createTime;type:bigint(20)"`                  // 创建时间戳
}

func (OperateRecord) TableName() string {
	return "operate_record"
}

func (WordsInfo) TableName() string {
	return "words_info"
}

// WordsInfoColumns get sql column name.获取数据库列名
var WordsInfoColumns = struct {
	ID         string
	Name       string
	BizType    string
	SceneType  string
	ParentName string
	UpdateTime string
	CreateTime string
}{
	ID:         "id",
	Name:       "name",
	BizType:    "bizType",
	SceneType:  "sceneType",
	ParentName: "parentName",
	UpdateTime: "updateTime",
	CreateTime: "createTime",
}

// WordsColumns get sql column name.获取数据库列名
var WordsColumns = struct {
	ID         string
	Text       string
	BizType    string
	SceneType  string
	Extent     string
	Valid      string
	UpdateTime string
	CreateTime string
}{
	ID:         "id",
	Text:       "text",
	BizType:    "bizType",
	SceneType:  "sceneType",
	Extent:     "extent",
	Valid:      "valid",
	UpdateTime: "updateTime",
	CreateTime: "createTime",
}

// OperateRecordColumns get sql column name.获取数据库列名
var OperateRecordColumns = struct {
	ID          string
	Text        string
	BizType     string
	SceneType   string
	Operator    string
	OperateType string
	CreateTime  string
}{
	ID:          "id",
	Text:        "text",
	BizType:     "bizType",
	SceneType:   "sceneType",
	Operator:    "operator",
	OperateType: "operateType",
	CreateTime:  "createTime",
}
