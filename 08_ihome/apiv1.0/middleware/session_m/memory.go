package session_m

import (
	"errors"
	"sync"
)

type MemorySession struct {
	sessionId string
	data map[string]interface{}
	rwlock sync.RWMutex
}

func NewMemorySession (id string) *MemorySession {
	s := &MemorySession{
		sessionId: id,
		data:      make(map[string]interface{},16),
	}
	return s
}

func (m *MemorySession) Set(key string,value interface{}) (err error) {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	m.data[key] = value
	return
}

func (m *MemorySession) Get(key string) (value interface{},err error) {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	value,ok := m.data[key]
	if !ok{
		err = errors.New("no session")
		return
	}
	return
}

func (m *MemorySession) Del(key string) (err error) {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()
	delete(m.data,key)
	return
}

func  (m *MemorySession) Save(key string) (err error) {
	return
}