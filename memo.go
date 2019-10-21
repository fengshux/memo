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

	return &Memo{expire, map[string]interface{}{}, map[string]time.Time{}, sync.RWMutex{}}
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
	if ok && current.After(last) {
		m.cache[key] = val
		m.time[key] = current
	}		
	m.rw.Unlock()
}
