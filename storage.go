package collapsar

import (
	"collapsar/policy"
	"errors"
	"log"
	"sync"
	"time"
)

var ExpiredKeyError = errors.New("expired key")
var NotFoundKeyError = errors.New("can't find key")

// 具体存储的节点
type Item struct {
	ttl   int64
	key   policy.KeyType
	value interface{}
}

func (i *Item) Set(key policy.KeyType, val interface{}, ttl int64) {
	i.key = key
	i.value = val
	i.ttl = ttl
}


// 存储的节点单元
type Node struct {
	locker sync.RWMutex
	// 当前可用的节点
	items map[policy.KeyType]Item
	// 淘汰策略
	eliminate policy.EliminateInterface
	// 失效处理函数
	failHandler FailHandlerFunc
	// 删除处理函数
	removeHandler RemoveHandlerFunc
}

func NewNode(option *Option) *Node {
	node := &Node{
		items:     make(map[policy.KeyType]Item),
		failHandler: option.FailHandler,
		removeHandler: option.RemoveHandler,
		eliminate: option.EliminateHandler,
	}
	return node
}

func (n *Node) clearExpire() error {
	n.locker.Lock()
	defer n.locker.Unlock()

	keys := make([]policy.KeyType, 0)
	for key := range n.items {
		keys = append(keys, key)
	}

	for _, key := range keys {
		if _, err := n.get(key, false); err != nil && err != NotFoundKeyError && err != ExpiredKeyError {
			log.Printf("error when call clearExpire, error:%+v", err)
		}
	}

	// 检查是否需要清理空间
	if n.eliminate != nil && n.eliminate.NeedEliminate() {
		removeKeys := n.eliminate.Eliminate()
		item := &Item{}
		for _, key := range removeKeys {
			item.key = key
			log.Printf("ready to remove key:%d", key)
			_, _ = n.remove(item, false)
		}
	}

	return nil
}

func (n *Node) get(key policy.KeyType, needLock bool) (*Item, error) {
	var item Item
	var ok bool
	if needLock {
		n.locker.RLock()
		item, ok = n.items[key]
		n.locker.RUnlock()
	} else {
		item, ok = n.items[key]
	}

	if ok {
		if n.eliminate != nil {
			n.eliminate.AccessKey(key, item.ttl)
		}
		// 超时key
		if item.ttl == -1 || time.Now().Unix()-item.ttl < 0 {
			return &item, nil
		} else {
			_, err := n.remove(&item, needLock)
			if err != nil {
				return nil, err
			}
			return &item, ExpiredKeyError
		}
	}

	return nil, NotFoundKeyError
}

func (n *Node) remove(item *Item, needLock bool) (interface{}, error) {
	if needLock {
		n.locker.Lock()
	}

	if n.removeHandler != nil {
		n.removeHandler(item.key, item.value)
	}

	delete(n.items, item.key)
	if n.eliminate != nil {
		n.eliminate.RemoveKey(item.key)
	}

	if needLock {
		n.locker.Unlock()
	}

	return item.value, nil
}

// 添加
func (n *Node) Add(key policy.KeyType, val interface{}, ttl int64) (interface{}, error) {
	n.locker.Lock()
	defer n.locker.Unlock()

	if item, ok := n.items[key]; ok {
		item.value = val
		item.ttl = ttl
		n.items[key] = item
		return item.value, nil
	} else {
		availItem := Item{}
		availItem.Set(key, val, ttl)
		n.items[key] = availItem
		if n.eliminate != nil {
			n.eliminate.AddKey(key, ttl)
		}
		return availItem.value, nil
	}
}

// 获取元素
func (n *Node) Get(key policy.KeyType) (interface{}, error) {
	item, err := n.get(key, true)
	if err == nil && item != nil {
		return item.value, nil
	}
	if (err == NotFoundKeyError || err == ExpiredKeyError) && n.failHandler != nil {
		val, err := n.failHandler(key)
		if err == nil {
			if item != nil {
				_, _ = n.Add(key, val, item.ttl)
			} else {
				_, _ = n.Add(key, val, -1)
			}
			return val, nil
	 	} else {
	 		return nil, err
		}
	}
	return nil, err
}

// 获取元素剩余的TTL时间
func (n *Node) TTL(key policy.KeyType) (int64, error) {
	item, err := n.get(key, true)
	if err == nil && item != nil {
		return time.Now().Unix() - item.ttl, nil
	}
	return -1, err
}

// 删除元素
func (n *Node) Remove(key policy.KeyType) (interface{}, error) {
	item, err := n.get(key, true)
	if err == nil && item != nil {
		return n.remove(item, true)
	}
	return nil, nil
}

// 清理过期节点
func (n *Node) ClearExpire() error {
	return n.clearExpire()
}
