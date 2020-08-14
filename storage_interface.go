package collapsar

import "collapsar/policy"

type StorageInterface interface {
	// 添加元素到缓存中
	// 如果添加成功则返回添加的val，如果更新则返回更新前的val
	Add(key policy.KeyType, val interface{}, ttl int64) (interface{}, error)

	// 获取元素
	Get(key policy.KeyType) (interface{}, error)

	// 获取元素剩余的TTL时间
	TTL(key policy.KeyType) (int64, error)

	// 删除元素
	Remove(key policy.KeyType) (interface{}, error)

	// 清理过期的元素集合
	ClearExpire() error
}
