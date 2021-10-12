package core

import (
	mather "github.com/dongweifly/sensitive-words-match"
	log "github.com/sirupsen/logrus"
	"green_shield/comm"
	"sync"
	"time"
)

var (
	GWordsMatchManager *WordsMatchManager
)

type WordsMatchManager struct {
	//词库管理
	m            map[string]*mather.MatchService
	replaceDelim rune
	//加载方式
	wordsLoad WordsLoad
	mutex     sync.RWMutex
}

func NewWordsMatchManager(wl WordsLoad) *WordsMatchManager {
	return &WordsMatchManager{
		wordsLoad:    wl,
		replaceDelim: '*',
		m:            make(map[string]*mather.MatchService),
	}
}

func (w *WordsMatchManager) Init() error {
	manger, err := w.wordsLoad.LoadAll()
	if err != nil {
		return err
	}
	w.m = manger

	time.AfterFunc(time.Second*10, w.UpdateLoop)

	return nil
}

func (w *WordsMatchManager) Update() {
	matchService := w.wordsLoad.Update()
	if matchService != nil && len(matchService) > 0 {
		w.mutex.Lock()
		for key, value := range matchService {
			w.m[key] = value
		}
		w.mutex.Unlock()
	}
}

func (w *WordsMatchManager) UpdateLoop() {
	//log.Info("WordsMatchManager UpdateLoop")
	w.Update()
	time.AfterFunc(10*time.Second, w.UpdateLoop)
}

//Match 匹配敏感词 和脱敏后的数据
//wordsName 词库名称, text要检测的词语
func (w *WordsMatchManager) Match(wordsName, text string) (sensitive []string, desensitization string) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	s, ok := w.m[wordsName]
	if !(ok) {
		log.Info("Can not find words " + wordsName)
		return nil, ""
	}
	if wordsName == "default_ad" {
		text = comm.FilterSpecialSymbols(text)
	}

	sensitive, desensitization = s.Match(text, '*')
	return
}
