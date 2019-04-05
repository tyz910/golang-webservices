package main

import (
	"math/rand"
	"sync"
)

type Session struct {
	Login     string
	Useragent string
}

type SessionID struct {
	ID string
}

const sessKeyLen = 10

type SessionManager struct {
	mu sync.RWMutex
	sessions map[SessionID]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		mu: sync.RWMutex{}
		sessions: map[SessionID]*Session{},
	}
}

func (sm *SessionManager) Create(in *Session) (*SessionID, error) {
	sm.mu.Lock()
	id := SessionID{RandStringRunes(sessKeyLen)}
	sm.mu.Unlock()
	sm.sessions[id] = in
	return &id, nil
}

func (sm *SessionManager) Check(in *SessionID) *Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if sess, ok := sm.sessions[*in]; ok {
		return sess
	}
	return nil
}

func (sm *SessionManager) Delete(in *SessionID) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, *in)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
