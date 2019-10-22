package memo

import (
	"time"
	"sync"
)

type Memo struct {
	expire time.Duration
	cache map[string]interface{}
	time map[string] time.Time
	rw sync.RWMutex
}

func New( expire time.Duration ) *Memo {
	m := &Memo{expire, map[string]interface{}{}, map[string]time.Time{}, sync.RWMutex{}}
	ticker := time.NewTicker(1 * time.Minute)
	done := make(chan bool)
	go func () {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				m.purge()
			}
		}
	}()
	return m
}

func Default() *Memo {
	return New(5 * time.Minute)
}

func (m *Memo) Get(key string) interface{} {
	
	m.rw.RLock()
	val := m.cache[key]
	last := m.time[key]
	m.rw.RUnlock()
	if last.Add(m.expire).Before(time.Now()) {
		val = nil
		m.rw.Lock()
		delete(m.cache, key)
		delete(m.time, key)
		m.rw.Unlock()
	}
	return val
}

func (m *Memo) Set(key string, val interface{})  {	
	current := time.Now()	
	m.rw.Lock()
	last, ok := m.time[key]
	if !ok {
		m.cache[key] = val
		m.time[key] = current		
	} else if current.After(last) {
		m.cache[key] = val
		m.time[key] = current
	}
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
	for _, k := range keys {
		m.rw.RLock()
		createTime := m.time[k]
		m.rw.RUnlock()
		if createTime.Add(m.expire).Before(time.Now()) {
			m.rw.Lock()
			delete(m.cache, k)
			delete(m.time, k)
			m.rw.Unlock()			
		}
	}
}
