package policy


type KeyType uint64
// 淘汰策略
type EliminateInterface interface {
	AddKey(key KeyType, ttl int64)
	RemoveKey(key KeyType)
	AccessKey(key KeyType, ttl int64)
	NeedEliminate() bool
	Eliminate() []KeyType
}
