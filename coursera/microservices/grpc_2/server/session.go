package main

import (
	"coursera/microservices/grpc/session"
	"fmt"
	"math/rand"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"golang.org/x/net/context"
)

const sessKeyLen = 10

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[session.SessionID]*session.Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		mu:       sync.RWMutex{},
		sessions: map[session.SessionID]*session.Session{},
	}
}

func (sm *SessionManager) Create(ctx context.Context, in *session.Session) (*session.SessionID, error) {
	fmt.Println("call Create", in)

	header := metadata.Pairs("header-key", "42")
	grpc.SendHeader(ctx, header)

	trailer := metadata.Pairs("trailer-key", "3.14")
	grpc.SetTrailer(ctx, trailer)

	id := &session.SessionID{RandStringRunes(sessKeyLen)}
	sm.mu.Lock()
	sm.sessions[*id] = in
	sm.mu.Unlock()
	return id, nil
}

func (sm *SessionManager) Check(ctx context.Context, in *session.SessionID) (*session.Session, error) {
	fmt.Println("call Check", in)
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
