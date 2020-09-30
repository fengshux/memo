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

// Memo is the Type of the instance of memory-cache
type Memo struct {
	defaultExpire time.Duration
	purgeInterval time.Duration
	cache         map[string]item
	rw            sync.RWMutex
}

// New funtion to create the instance
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

// Default function to  new a instance by defaul expire and purge
func Default() *Memo {
	return New(1*time.Minute, 10*time.Minute)
}

// Get value
func (m *Memo) Get(key string) interface{} {

	itm := m.getNotExpire(key)
	if itm == nil {
		return nil
	}
	
	if itm != nil && itm.expire.Before(time.Now()) {
		m.expireKey(key)
		return nil
	}
	return itm.val
}

func (m *Memo) getNotExpire(key string) *item {
	m.rw.RLock()
	defer m.rw.RUnlock()
	itm, ok := m.cache[key]
	if !ok {
		return nil
	}
	return &itm
}

func (m *Memo) expireKey( key string) {
	m.rw.Lock()
	defer m.rw.Unlock()
	if m.cache[key].expire.Before(time.Now()) {
		delete(m.cache, key)
	}
}

// Set a value for a key
func (m *Memo) Set(key string, val interface{}) {
	m.SetEx(key, m.defaultExpire, val)
}

// SetEx set a value for a key with expire
func (m *Memo) SetEx(key string, expire time.Duration, val interface{}) {

	current := time.Now()
	m.rw.Lock()
	defer m.rw.Unlock()
	itm, ok := m.cache[key]
	if !ok || current.After(itm.set) {
		m.cache[key] = item{
			val,
			current,
			current.Add(expire),
		}
	}
}

// Del delete chache
func (m *Memo) Del(key string) {
	m.rw.Lock()
	defer m.rw.Unlock()
	delete(m.cache, key)
}

func (m *Memo) length() int {
	m.rw.RLock()
	defer m.rw.RUnlock()
	return len(m.cache)
}


func (m *Memo) expirekeys(keys []string, num int) []string {
	m.rw.RLock()
	defer m.rw.RUnlock()
	count := 0
	for k := range m.cache {
		if count > num - 1 {
			break
		}
		keys = append(keys, k)
		count = count + 1
	}
	return keys
}

func expireRoundAndShred(length int) (round, shred int) {
	// count span
	span := 1
	if length > 10000 {
		for j := length; j > 10; j = j / 10 {
			span = span + 1
		}
	}

	expireCount := length / span
	round = span + 1;
	shred = expireCount / round
	if shred > 5000 {
		shred = 5000
		round = expireCount / shred
	}
	return
}



func (m *Memo) purge() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("purge has an error: ", r)
			debug.PrintStack()
		}
	}()
	fmt.Println("purge start")
	// get round of expire and count of every round
	length := m.length()	
	round, shred := expireRoundAndShred(length)
	fmt.Println("length,round, shred", length, round, shred)
	keys := make([]string, 0, shred)
	
	for j := 0; j < round; j = j + 1 {
		keys = keys[0:0]
		// get keys
		keys = m.expirekeys(keys, shred)

		// range keys and expire
		for _, key := range keys {
			m.Get(key)
		}

	}
	fmt.Println("purge stop")
}
