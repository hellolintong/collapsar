package policy

import (
	"sync"
)

const LFUEliminateRate = 0.2
const LFUMaxNumber = 30000

type LFUEliminate struct {
	container *EliminateContainer
	eliminateRate      float32
	maxNumber          int
}

func LFUAddPolicyFunc(container *sync.Map, key string, ttl int64) {
	if counter, ok := container.Load(key); ok {
		container.Store(key, counter.(int) + 1)
	} else {
		container.Store(key, 1)
	}
}

func (l *LFUEliminate) AddKey(key string, ttl int64) {
	l.container.addKey(key, ttl, LFUAddPolicyFunc, LFUAddPolicyFunc)
}

func (l *LFUEliminate) RemoveKey(key string) {
	l.container.removeKey(key)
}

func (l *LFUEliminate) AccessKey(key string, ttl int64) {
	l.container.addKey(key, ttl, LFUAddPolicyFunc, LFUAddPolicyFunc)
}


func (l *LFUEliminate) Eliminate() []string {
	return l.container.eliminate(l.eliminateRate)
}

func (l *LFUEliminate) NeedEliminate() bool {
	return l.container.needEliminate(l.maxNumber)
}

func NewLFUEliminate() *LFUEliminate {
	return &LFUEliminate{
		eliminateRate:      LFUEliminateRate,
		maxNumber:          LFUMaxNumber,
		container: &EliminateContainer{},
	}
}
