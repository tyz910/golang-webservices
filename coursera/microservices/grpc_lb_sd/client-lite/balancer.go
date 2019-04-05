package main

import (
	"fmt"

	"google.golang.org/grpc/naming"
)

// load balancer
// Copyright 2016 gRPC authors https://github.com/grpc/grpc-go/blob/master/balancer_test.go

type testWatcher struct {
	// the channel to receives name resolution updates
	update chan *naming.Update
	// the side channel to get to know how many updates in a batch
	side chan int
	// the channel to notifiy update injector that the update reading is done
	readDone chan int
}

func (w *testWatcher) Next() (updates []*naming.Update, err error) {
	n := <-w.side
	if n == 0 {
		return nil, fmt.Errorf("w.side is closed")
	}
	for i := 0; i < n; i++ {
		u := <-w.update
		if u != nil {
			updates = append(updates, u)
		}
	}
	w.readDone <- 0
	return
}

func (w *testWatcher) Close() {
	close(w.side)
}

// Inject naming resolution updates to the testWatcher.
func (w *testWatcher) inject(updates []*naming.Update) {
	w.side <- len(updates)
	for _, u := range updates {
		w.update <- u
	}
	<-w.readDone
}

type testNameResolver struct {
	w    *testWatcher
	addr string
}

func (r *testNameResolver) Resolve(target string) (naming.Watcher, error) {
	r.w = &testWatcher{
		update:   make(chan *naming.Update, 1),
		side:     make(chan int, 1),
		readDone: make(chan int),
	}
	r.w.side <- 1
	r.w.update <- &naming.Update{
		Op:   naming.Add,
		Addr: r.addr,
	}
	go func() {
		<-r.w.readDone
	}()
	return r.w, nil
}
