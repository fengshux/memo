package memo

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

type item struct {
	val    interface{}
	set    time.Time
	expire time.Time
}

type Memo struct {
	defaultExpire time.Duration
	purgeInterval time.Duration
	cache         map[string]item
	rw            sync.RWMutex
}

func New(defaultExpire, purgeInterval time.Duration) *Memo {
	m := &Memo{defaultExpire, purgeInterval, map[string]item{}, sync.RWMutex{}}
	ticker := time.NewTicker(purgeInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				m.purge()
			}
		}
	}()
	return m
}

func Default() *Memo {
	return New(1*time.Minute, 10*time.Minute)
}

func (m *Memo) Get(key string) interface{} {
	m.rw.RLock()
	// if the key doesn't exist, the itm is zero value
	itm := m.cache[key]
	m.rw.RUnlock()

	if itm.expire.Before(time.Now()) {
		m.rw.Lock()
		if m.cache[key].expire.Before(time.Now()) {
			delete(m.cache, key)
		}
		m.rw.Unlock()
		return nil
	}
	return itm.val
}

func (m *Memo) Set(key string, val interface{}) {
	m.SetEx(key, m.defaultExpire, val)
}

func (m *Memo) SetEx(key string, expire time.Duration, val interface{}) {

	current := time.Now()
	m.rw.Lock()
	itm, ok := m.cache[key]
	if !ok || current.After(itm.set) {
		m.cache[key] = item{
			val,
			current,
			current.Add(expire),
		}
	}
	m.rw.Unlock()
}

func (m *Memo) Del(key string) {
	m.rw.Lock()
	delete(m.cache, key)
	m.rw.Unlock()
}

func getSpan(length int) int {
	sp := 1
	if length <= 1000 {
		return sp
	}

	for length > 10 {
		sp = sp + 1
		length = length / 10
	}
	return sp
}

func (m *Memo) purge() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("purge has an error: ", r)
			debug.PrintStack()
		}
	}()

	m.rw.RLock()
	length := len(m.cache)
	m.rw.RUnlock()

	span := getSpan(length)
	expireCount := length / span
	shard := expireCount / span

	if shard > 5000 {
		shard = 5000
		span = expireCount / shard
	}
	//	keys := make([]string, shard)
	keys := make([]string, 0)

	// 当每次执行的过多的时候，分成span个shard来执行
	for j := 0; j < span; j = j + 1 {
		// get keys
		m.rw.RLock()
		count := 0
		for k, _ := range m.cache {
			if count > shard-1 {
				break
			}
			keys[count] = k
			count = count + 1
		}
		m.rw.RUnlock()

		// range keys and expire
		for i := 0; i < shard; i = i + 1 {
			key := keys[i]
			m.rw.RLock()
			itm, ok := m.cache[key]
			m.rw.RUnlock()
			if ok && itm.expire.Before(time.Now()) {
				m.rw.Lock()
				if m.cache[key].expire.Before(time.Now()) {
					delete(m.cache, key)
				}
				m.rw.Unlock()
			}
		}

	}

}
