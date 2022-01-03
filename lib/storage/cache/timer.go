package cache

import (
	"sync/atomic"
	"time"
)

type Timer interface {
	Now() uint32
}

type StoppableTimer interface {
	Timer
	Stop()
}

func getUnixTime() uint32 {
	return uint32(time.Now().Unix())
}

type defaultTimer struct {
}

func (timer defaultTimer) Now() uint32 {
	return getUnixTime()
}

type cachedTimer struct {
	now    uint32
	ticker *time.Ticker
	done   chan bool
}

func NewCachedTimer() StoppableTimer {
	timer := &cachedTimer{
		now:    getUnixTime(),
		ticker: time.NewTicker(time.Second),
		done:   make(chan bool, 1),
	}
	go timer.update()
	return timer
}

func (timer *cachedTimer) Now() uint32 {
	return atomic.LoadUint32(&timer.now)
}

func (timer *cachedTimer) Stop() {
	timer.ticker.Stop()
	timer.done <- true
	close(timer.done)
	timer.done = nil
	timer.ticker = nil
}

func (timer *cachedTimer) update() {
	for {
		select {
		case <-timer.done:
			return
		case <-timer.ticker.C:
			atomic.StoreUint32(&timer.now, getUnixTime())
		}
	}
}
