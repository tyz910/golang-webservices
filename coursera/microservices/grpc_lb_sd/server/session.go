package main

import (
	"coursera/microservices/grpc/session"
	"fmt"
	"math/rand"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"golang.org/x/net/context"
)

const sessKeyLen = 10

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[session.SessionID]*session.Session
	host     string
}

func NewSessionManager(port string) *SessionManager {
	return &SessionManager{
		mu:       sync.RWMutex{},
		sessions: map[session.SessionID]*session.Session{},
		host:     port,
	}
}

func (sm *SessionManager) Create(ctx context.Context, in *session.Session) (*session.SessionID, error) {
	fmt.Println("call Create", in)
	id := &session.SessionID{RandStringRunes(sessKeyLen)}
	sm.mu.Lock()
	sm.sessions[*id] = in
	sm.mu.Unlock()
	return id, nil
}

func (sm *SessionManager) Check(ctx context.Context, in *session.SessionID) (*session.Session, error) {
	fmt.Println("call Check", in)
	// между сервисами нет общения, возвращаем заглушку
	fakeLogin := sm.host + " " + in.GetID()
	return &session.Session{Login: fakeLogin}, nil

	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if sess, ok := sm.sessions[*in]; ok {
		return sess, nil
	}
	return nil, grpc.Errorf(codes.NotFound, "session not found")
}

func (sm *SessionManager) Delete(ctx context.Context, in *session.SessionID) (*session.Nothing, error) {
	fmt.Println("call Delete", in)
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, *in)
	return &session.Nothing{Dummy: true}, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
