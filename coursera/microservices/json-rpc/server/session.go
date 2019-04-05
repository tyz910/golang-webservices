package main

import (
	"fmt"
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
	mu       sync.RWMutex
	sessions map[SessionID]*Session
}

func NewSessManager() *SessionManager {
	return &SessionManager{
		mu:       sync.RWMutex{},
		sessions: map[SessionID]*Session{},
	}
}

func (sm *SessionManager) Create(in *Session, out *SessionID) error {
	fmt.Println("call Create", in)
	id := &SessionID{RandStringRunes(sessKeyLen)}
	sm.mu.Lock()
	sm.sessions[*id] = in
	sm.mu.Unlock()
	*out = *id
	return nil
}

func (sm *SessionManager) Check(in *SessionID, out *Session) error {
	fmt.Println("call Check", in)
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if sess, ok := sm.sessions[*in]; ok {
		*out = *sess
	}
	return nil
}

func (sm *SessionManager) Delete(in *SessionID, out *int) error {
	fmt.Println("call Delete", in)
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, *in)
	*out = 1
	return nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
