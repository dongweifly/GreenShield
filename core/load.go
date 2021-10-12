package core

import (
	mather "github.com/dongweifly/sensitive-words-match"
)

type WordsLoad interface {
	//Name turn WordsLoader's name
	Name() string

	//LoadAll 把所有的词库load出来p
	LoadAll() (map[string]*mather.MatchService, error)

	//Update 获取更新到的 MatchService
	Update() map[string]*mather.MatchService
}
