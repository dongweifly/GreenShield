package core

import (
	"bufio"
	mather "github.com/dongweifly/sensitive-words-match"

	log "github.com/sirupsen/logrus"
	"green_shield/comm"
	"os"
	"path/filepath"
)

type LocalWordLoad struct {
	RootDir string
}

func NewLocalWordLoad(rootDir string) *LocalWordLoad {
	return &LocalWordLoad{
		RootDir: rootDir,
	}
}

func (l *LocalWordLoad) readWordsFromFile(localFileName string) ([]string, error) {
	file, err := os.Open(localFileName)
	if err != nil {
		log.Warnf("LocalWordLoad %s fail : %s", localFileName, err.Error())
		return nil, err
	}

	defer file.Close()

	fi, _ := file.Stat()

	//主要是为性能考虑，append []string
	resp := make([]string, 0, comm.Round8(int(fi.Size())/3))

	scanner := bufio.NewScanner(file)
	//注意: optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		//把不需要的字符过滤掉
		text := comm.TrimString(scanner.Text())
		resp = append(resp, text)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return resp, nil
}

func (l *LocalWordLoad) LoadAll() (map[string]*mather.MatchService, error) {
	var files []string
	//只关注rootDir下面的文件，暂时不支持子文件夹的形式
	err := filepath.Walk(l.RootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Warnf("load %s fail %s", l.RootDir, err.Error())
			return nil
		}

		if !info.IsDir() {
			files = append(files, info.Name())
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	d := make(map[string]*mather.MatchService, len(files))
	for i := 0; i < len(files); i++ {
		if matchService, err := l.loadOne(files[i]); err == nil {
			d[files[i]] = matchService
		}
	}
	return d, nil
}

func (l *LocalWordLoad) Name() string {
	return "load_local"
}

//Update 获取更新到的 matchService
func (l *LocalWordLoad) Update() map[string]*mather.MatchService {
	return nil
}

func (l *LocalWordLoad) loadOne(sensitiveWordsName string) (*mather.MatchService, error) {
	if words, err := l.readWordsFromFile(l.RootDir + "/" + sensitiveWordsName); err == nil {
		s := mather.NewMatchService()
		s.Build(words)
		return s, nil
	} else {
		return nil, err
	}
}
