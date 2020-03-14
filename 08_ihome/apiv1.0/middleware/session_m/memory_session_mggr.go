package session_m

import (
	"sync"
	uuid "github.com/satori/go.uuid"
	)

type MemorySessionMgr struct {
	sessionMap map[string]Session
    rwlock sync.RWMutex
}

func NewMemorySessionMgr() *MemorySessionMgr {
	sr := &MemorySessionMgr{
		sessionMap: make(map[string]Session,1024),
	}
	return sr
}

func (s *MemorySessionMgr) Init(addr string, options ...string) (err error) {
	return
}

func (s *MemorySessionMgr) CreateSession()(session Session,err error) {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()

	id := uuid.NewV4()
	sessionId := id.String()
	session = NewMemorySession(sessionId)

}

