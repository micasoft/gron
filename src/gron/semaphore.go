package gron

import (
	"sync"
)

type Semaphore struct {
	permits int
	avail   int
	channel chan int
	aMutex  *sync.Mutex
	rMutex  *sync.Mutex
}

func NewSemaphore(permits int) *Semaphore {
	if permits < 1 {
		panic("Invalid number of permits. Less than 1")
	}
	return &Semaphore{
		permits,
		permits,
		make(chan int, permits),
		&sync.Mutex{},
		&sync.Mutex{},
	}
}

//Acquire one permit, if its not available the goroutine will block till its available
func (s *Semaphore) Acquire() {
	s.aMutex.Lock()
	s.channel <- 1
	s.avail--
	s.aMutex.Unlock()
}

//Release one permit
func (s *Semaphore) Release() {
	s.rMutex.Lock()
	<-s.channel
	s.avail++
	s.rMutex.Unlock()
}

func (s *Semaphore) Available() int {
	return s.permits - s.avail
}
