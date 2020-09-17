package memo

import (
	"time"
	"sync"
)

type Item struct{
	val interface{}
	set time.Time
	expire time.Time
}

type Memo struct {
	defaultExpire time.Duration
	purgeInterval time.Duration
	cache map[string]Item
	rw sync.RWMutex
}

func New( defaultExpire, purgeInterval time.Duration ) *Memo {
	m := &Memo{  defaultExpire, purgeInterval,  map[string]Item{}, sync.RWMutex{}}
	ticker := time.NewTicker(purgeInterval)
	go func () {
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
	return New(1 * time.Minute, 10 * time.Minute)
}

func (m *Memo) Get(key string) interface{} {
	m.rw.RLock()
	item := m.cache[key]
	m.rw.RUnlock()
	if item.expire.Before(time.Now()) {
		return nil
		m.rw.Lock()
		if m.cache[key].expire.Before(time.Now()) {
			delete(m.cache, key)
		}
		m.rw.Unlock()
	}
	return item.val
}

func (m *Memo) Set(key string, val interface{}) {
	m.SetEx(key, m.defaultExpire, val)
}

func (m *Memo) SetEx(key string, expire time.Duration, val interface{}) {

	current := time.Now()
	m.rw.Lock()
	item, ok := m.cache[key]
	if !ok || current.After(item.set) {
		m.cache[key] = Item{
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


func (m *Memo) purge() {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	
	keys := []string{}
	m.rw.RLock()
	for k , _ := range m.cache {
		keys = append(keys, k)
	}
	m.rw.RUnlock()
	length := len(keys)
	span := 1
	if length > 1000 {
		span = 3
	}
	
	for i:=0 ;i< length; i = i+span {
		key := keys[i]
		m.rw.RLock()
		item, ok := m.cache[key]
		m.rw.RUnlock()
		if ok && item.expire.Before(time.Now()) {
			m.rw.Lock()
			if m.cache[key].expire.Before(time.Now()) {
				delete(m.cache, key)
			}
			m.rw.Unlock()		
		}
	}
}
