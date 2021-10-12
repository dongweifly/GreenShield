package core

type MatchHandler interface {
	// Match 文本规则匹配
	Match(wordsName, text string) (sensitive []string, desensitization string)
}
