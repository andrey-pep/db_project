package main

import (
	/*"sync"
	"fmt"*/
	"io"
	"crypto/rand"
	"encoding/base64"
	//"github.com/davecgh/go-spew/spew"
	/*"net/http"
	"time"
	//"net/url"*/
)

type Manager struct {
	cookieName	string
	maxlifetime	int64
	Sessions  	map[string]*Session
}

func NewManager(cookieName string, maxlifetime int64) (*Manager, error) {
    return &Manager{cookieName: cookieName, maxlifetime: maxlifetime, Sessions: make(map[string]*Session)}, nil
}

type Session struct {
	Values	map[string]interface{}
    /*Set(key, value interface{}) error //set session value
    Get(key interface{}) interface{}  //get session value
    Delete(key interface{}) error     //delete session value
    SessionID() string                //back current sessionID*/
}

func (mg *Manager) CheckSession(sessionId string) *Session {
	if val, ok := mg.Sessions[sessionId]; ok {
		return val
	}
	return nil
}

func (manager *Manager) DeleteSession (sessionId string) {
	delete(manager.Sessions, sessionId)
}

func (manager *Manager) sessionId() string {
    b := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, b); err != nil {
        return ""
    }
    return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) SessionInit() string {
	id := manager.sessionId()
	manager.Sessions[id] = manager.NewSession(id)
	return id
}

func (manager *Manager) NewSession(id string) *Session {
	return &Session{Values: make(map[string]interface{})}
}

func (s *Session) GetValue(name string) interface{} {
	if val, ok := s.Values[name]; ok {
		return val
	}
	return nil
}

func (s *Session) SetValue(arg string, name string) {
	s.Values[name] = arg
}

/*func (manager *Manager) GC() {
    manager.lock.Lock()
    defer manager.lock.Unlock()
    manager.provider.SessionGC(manager.maxlifetime)
    time.AfterFunc(time.Duration(manager.maxlifetime), func() { manager.GC() })
}*/