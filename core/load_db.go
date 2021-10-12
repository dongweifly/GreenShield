package core

import (
	mather "github.com/dongweifly/sensitive-words-match"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"green_shield/model"
)

type DatabaseWordLoad struct {
	db        *gorm.DB
	wordsInfo map[string]model.WordsInfo
}

func NewDatabaseWordLoad(dbAddress string) *DatabaseWordLoad {

	db, err := gorm.Open(mysql.New(mysql.Config{
		DriverName: "mysql",
		DSN:        dbAddress,
	}))

	if err != nil {
		log.Fatal("Init DB fail ", err.Error())
		return nil
	}

	return &DatabaseWordLoad{
		db: db,
	}
}

func (l *DatabaseWordLoad) Name() string {
	return "load_db"
}

func (l *DatabaseWordLoad) Update() map[string]*mather.MatchService {
	wordInfos := l.diffWordInfo()
	if len(wordInfos) <= 0 {
		return nil
	}

	log.Infof("Load diff words info: %v", wordInfos)

	res := make(map[string]*mather.MatchService, len(wordInfos))
	for _, info := range wordInfos {
		if s := l.loadByType(info.BizType, info.SceneType); s != nil {
			res[info.BizType+"_"+info.SceneType] = s
		}
	}

	return res
}

// loadWordInfo 从数据库中加载green_shield_words_info的信息
func (l *DatabaseWordLoad) loadWordInfo() map[string]model.WordsInfo {
	var wordsInfo []model.WordsInfo
	err := l.db.Find(&wordsInfo).Error
	if err != nil {
		return nil
	}

	m := make(map[string]model.WordsInfo, len(wordsInfo))

	for _, info := range wordsInfo {
		m[info.Name] = info
	}

	return m
}

func (l *DatabaseWordLoad) diffWordInfo() []model.WordsInfo {
	newWordInfo := l.loadWordInfo()
	if newWordInfo == nil || len(newWordInfo) == 0 {
		log.Info("get words info is empty")
		return nil
	}

	var modifies []model.WordsInfo

	for _, n := range newWordInfo {
		if info, ok := l.wordsInfo[n.Name]; ok {
			//更新时间发生变更
			if info.UpdateTime < n.UpdateTime {
				modifies = append(modifies, n)
				l.wordsInfo[n.Name] = n
			}
		} else {
			//如果不存在，那么获取更新
			if modifies == nil {
				modifies = make([]model.WordsInfo, 0)
			}
			modifies = append(modifies, n)
			if l.wordsInfo == nil {
				l.wordsInfo = make(map[string]model.WordsInfo)
			}
			l.wordsInfo[n.Name] = n
		}
	}
	return modifies
}

//loadByType 根据bizType， sceneType构建新的DFA tree
func (l *DatabaseWordLoad) loadByType(bizType, sceneType string) *mather.MatchService {

	var words []model.Words

	err := l.db.Where(&model.Words{
		Valid:     true,
		BizType:   bizType,
		SceneType: sceneType,
	}).Find(&words).Error

	if err != nil {
		log.Error(err.Error())
		return nil
	}

	v := make([]string, 0, 1024)
	for _, word := range words {
		v = append(v, word.Text)
	}

	s := mather.NewMatchService()
	s.Build(v)
	return s
}

//LoadAll 因为只会在启动的调用一次，相当于初始化了一次Update
func (l *DatabaseWordLoad) LoadAll() (map[string]*mather.MatchService, error) {
	return l.Update(), nil
}
