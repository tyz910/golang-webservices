/*
	Пример: аукцион
	Аукцион состоит из нескольких участников, которые делают постоянно увеличивающиеся ставки по лоту
	Он может закончиться в следующих случаях:
	1. дана максимальная ставка
	2. сделано определённое количество ставок
	3. прошло определённое ремя
*/
package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Lot struct {
	sync.Mutex
	CurrentBid int
	PlayerID   int
	MaxPrice   int // цена, при которой лот отдаётся назвавшему её игроку
	MaxBids    int // не более Х ставок
	currentCnt int
}

func (b *Lot) GetCurrentBid() int {
	b.Lock()
	defer b.Unlock()
	return b.CurrentBid
}

func (b *Lot) SetNewBid(newBid PlayerBid) bool {
	b.Lock()
	defer b.Unlock()
	if newBid.Bid > b.CurrentBid {
		fmt.Printf("new bid: %+v\n", newBid)

		b.CurrentBid = newBid.Bid
		b.PlayerID = newBid.PlayerID
		b.currentCnt++
		// заканчиваем аукцион если названа максимальная цена
		if b.MaxBids <= b.currentCnt {
			println("finish by count")
			return true
		}
		// заканчиваем аукцион если прошло фиксированное количество ставок
		if b.MaxPrice <= newBid.Bid {
			println("finish by bid")
			return true
		}
	}
	return false
}

type PlayerBid struct {
	Bid      int
	PlayerID int
}

const (
	AvgSleep   = 20000
	AvgBidStep = 100
)

func makeBid(ctx context.Context, playerID int, lot *Lot, bids chan PlayerBid) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(AvgSleep)))
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		myBid := lot.GetCurrentBid() + rand.Intn(AvgBidStep)

		select {
		case <-ctx.Done():
			return
		case bids <- PlayerBid{PlayerID: playerID, Bid: myBid}:
		default:
		}
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(AvgSleep)))
	}
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	lot := &Lot{CurrentBid: 100, MaxBids: 100, MaxPrice: 10000}
	bids := make(chan PlayerBid)
	// создаём контекст, который отсигналит через указанный помежуток
	ctx, finish := context.WithTimeout(context.Background(), time.Second*10)

	// запускаем 5 игроков которые будут делать ставки
	for i := 0; i < 5; i++ {
		go makeBid(ctx, i, lot, bids)
	}

LOOP:
	for {
		select {
		case bid := <-bids:
			if lot.SetNewBid(bid) {
				finish()
				break LOOP
			}
		case <-ctx.Done():
			fmt.Println("Done:", ctx.Err())
			break LOOP

		}
	}

	fmt.Println("Auction finished by player", lot.PlayerID, "with price", lot.CurrentBid)

}
