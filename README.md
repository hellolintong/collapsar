# collapsar

collapsar是线程安全的进程缓存库

## 使用Demo
```go
func main() {
    option := &Option{
		Length: 514,
	}
	cache := NewCache(option)
	oldValue, err := cache.AddWithTTL("test", 11, 3)
	val, err := cache.Get("test")
	if err != nil || val == nil || val.(int) != 11 {
		t.Errorf("check ttl timeout fail")
	}
}
```