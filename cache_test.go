package collapsar

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCacheNew(t *testing.T) {
	option := &Option{
		Length: 15,
	}
	cache := NewCache(option)
	if cache.length != 16 {
		t.Errorf("cache length test fail, input:%d, expect:%d, actual:%d ", 15, 16, cache.length)
	}

	if cache.offset != 4 {
		t.Errorf("cache offset test fail, input:%d, expect:%d, actual:%d ", 15, 4, cache.offset)
	}
}

func TestCacheAdd(t *testing.T) {

	option := &Option{
		Length: 514,
	}
	cache := NewCache(option)

	wg1 := sync.WaitGroup{}
	wg1.Add(5)
	ch := make(chan int, 1000)
	addWorker := func(ch chan int) {
		for i := range ch {
			if i == -1 {
				wg1.Done()
				return
			}
			_, err := cache.Add(fmt.Sprintf("test_%d", i), fmt.Sprintf("test_value_%d", i))
			if err != nil {
				t.Errorf("cache add test fail, error:%s", err.Error())
			}
		}
	}

	go addWorker(ch)
	go addWorker(ch)
	go addWorker(ch)
	go addWorker(ch)
	go addWorker(ch)
	for i := 0; i < 1000000; i++ {
		ch <- i
	}

	for i := 0; i < 10; i++ {
		_, err := cache.Add(fmt.Sprintf("test_%d", i), fmt.Sprintf("new_test_value_%d", i))
		if err != nil {
			t.Errorf("cache repeated add test fail, error:%s", err.Error())
		}
	}
	ch <- -1
	ch <- -1
	ch <- -1
	ch <- -1
	ch <- -1

	wg1.Wait()

	wg2 := sync.WaitGroup{}
	wg2.Add(5)
	ch2 := make(chan int, 1000)

	getWorker := func(ch chan int) {
		for i := range ch {
			if i == -1 {
				wg2.Done()
				return
			}
			val, err := cache.Get(fmt.Sprintf("test_%d", i))
			if err == nil {
				value := val.(string)
				if i < 10 {
					expectedValue := fmt.Sprintf("new_test_value_%d", i)
					if value != expectedValue {
						t.Errorf("cache get test fail, expect:%s, acutial:%s", expectedValue, value)
					}
				} else {
					expectedValue := fmt.Sprintf("test_value_%d", i)
					if value != expectedValue {
						t.Errorf("cache get test fail, expect:%s, acutial:%s", expectedValue, value)
					}
				}
			}
		}
	}
	go getWorker(ch2)
	go getWorker(ch2)
	go getWorker(ch2)
	go getWorker(ch2)
	go getWorker(ch2)
	for i := 0; i < 1000000; i++ {
		ch2 <- i
	}
	ch2 <- -1
	ch2 <- -1
	ch2 <- -1
	ch2 <- -1
	ch2 <- -1
	wg2.Wait()
}

func TestCacheTTL(t *testing.T) {
	option := &Option{
		Length: 514,
	}
	cache := NewCache(option)
	_, _ = cache.AddWithTTL("test", 10, 3)
	time.Sleep(10 * time.Second)
	val, err := cache.Get("test")
	if err != ExpiredKeyError || val != nil {
		t.Errorf("check ttl timeout fail, err:%+v, val:%+v", err, val)
	}
	_, _ = cache.AddWithTTL("test", 11, 3)
	val, err = cache.Get("test")
	if err != nil || val == nil || val.(int) != 11 {
		t.Errorf("check ttl timeout fail")
	}
}
