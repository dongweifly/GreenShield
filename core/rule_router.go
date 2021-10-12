package core

import "sync"

var (
	GRuleRouterManager = NewRuleRouterManager()
)

type RuleRouterManager struct {
	rules map[string]MatchHandler
	mutex sync.RWMutex
}

func NewRuleRouterManager() *RuleRouterManager {
	return &RuleRouterManager{
		rules: make(map[string]MatchHandler),
	}
}

func (r *RuleRouterManager) Register(name string, handler MatchHandler) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.rules[name] = handler
}

func (r *RuleRouterManager) Get(name string) MatchHandler {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if rule, ok := r.rules[name]; ok {
		return rule
	}

	return GWordsMatchManager
}
