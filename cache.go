package collapsar

import (
	"collapsar/hash"
	"log"
	"time"
)

const PeriodDuration = 3

func init() {
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
}
// 失效处理函数
type FailHandlerFunc = func(string) (interface{}, error)

// 删除处理函数
type RemoveHandlerFunc = func(string, interface{})

type Cache struct {
	// 偏移量
	offset uint32
	// 存储的单元长度
	length int
	// 存储
	storageList []StorageInterface
	// 计算哈希值
	calculator HashInterface

	// 定期清理下标
	clearIndex int
}

func NewCache(option *Option) *Cache {
	length := option.Length
	i := 0
	adjustLength := 1
	if length > 0 {
		for adjustLength < length {
			i++
			adjustLength = adjustLength << 1
		}
		if i > 8 {
			i = 8
			adjustLength = 256
		}
	} else {
		adjustLength = 16
		i = 4
	}
	var calculator HashInterface
	calculator = option.Calculator
	if calculator == nil {
		calculator = hash.NewAdlerHash()
	}
	cache := &Cache{
		offset:      uint32(i),
		length:      adjustLength,
		storageList: make([]StorageInterface, adjustLength),
		calculator:  calculator,
		clearIndex:  0,
	}
	for i := 0; i < adjustLength; i++ {
		cache.storageList[i] = NewNode(option)
	}
	// 启动周期执行任务
	go cache.periodicWorker()
	return cache
}

func (c *Cache) hash(key string) uint32 {
	h := c.calculator.Hash(key)
	return h
}

// 计算哈希并获取下表
func (c *Cache) getIndex(hash uint32) uint32 {
	// 高位和低位数异或混淆，然后进行取模运算生成下标
	index := ((hash >> c.offset) ^ hash) & (c.offset - 1)
	return index
}

// 添加元素到缓存中
func (c *Cache) Add(key string, val interface{}) (interface{}, error) {
	return c.AddWithTTL(key, val, -1)
}

// 添加元素到缓存中
// 如果添加成功则返回添加的val，如果更新则返回更新前的val
func (c *Cache) AddWithTTL(key string, val interface{}, ttl int64) (interface{}, error) {
	h := c.hash(key)
	index := c.getIndex(h)
	if ttl == -1 {
		return c.storageList[index].Add(key, val, -1)
	}
	return c.storageList[index].Add(key, val, ttl+time.Now().Unix())
}

// 获取元素
func (c *Cache) Get(key string) (interface{}, error) {
	h := c.hash(key)
	index := c.getIndex(h)
	return c.storageList[index].Get(key)
}

// 获取元素剩余的TTL时间
func (c *Cache) TTL(key string) (int64, error) {
	h := c.hash(key)
	index := c.getIndex(h)
	return c.storageList[index].TTL(key)

}

// 清理超时的节点
func (c *Cache) clearExpiredItems() {
	c.clearIndex = (c.clearIndex + 1) % c.length
	if err := c.storageList[c.clearIndex].ClearExpire(); err != nil {
		log.Printf("error when call clearExpiredItems, error:%+v", err)
	}
}

// 定期执行任务
func (c *Cache) periodicWorker() {
	t := time.NewTicker(time.Second * PeriodDuration)
	defer t.Stop()

	for {
		<-t.C
		c.clearExpiredItems()
	}
}
