package policy


// 淘汰策略
type EliminateInterface interface {
	AddKey(key string, ttl int64)
	RemoveKey(key string)
	AccessKey(key string, ttl int64)
	NeedEliminate() bool
	Eliminate() []string
}
