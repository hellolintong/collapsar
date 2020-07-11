package policy

import (
	"sync"
	"time"
)

const LRUEliminateRate = 0.2
const LRUMaxNumber = 30000


type LRUEliminate struct {
	container *EliminateContainer
	eliminateRate      float32
	maxNumber          int
}

func LRUAddPolicyFunc(container *sync.Map, key string, ttl int64) {
	container.Store(key, time.Now().Unix())
}

func (l *LRUEliminate) AddKey(key string, ttl int64) {
	l.container.addKey(key, ttl, LRUAddPolicyFunc, LRUAddPolicyFunc)
}

func (l *LRUEliminate) RemoveKey(key string) {
	l.container.removeKey(key)
}

func (l *LRUEliminate) AccessKey(key string, ttl int64) {
	l.container.addKey(key, ttl, LRUAddPolicyFunc, LRUAddPolicyFunc)
}


func (l *LRUEliminate) Eliminate() []string {
	return l.container.eliminate(l.eliminateRate)
}

func (l *LRUEliminate) NeedEliminate() bool {
	return l.container.needEliminate(l.maxNumber)
}

func NewLRUEliminate() *LRUEliminate {
	return &LRUEliminate{
		eliminateRate:      LRUEliminateRate,
		maxNumber:          LRUMaxNumber,
		container: &EliminateContainer{},
	}
}
