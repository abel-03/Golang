//go:build !solution

package batcher

import (
	"gitlab.com/slon/shad-go/batcher/slow"
	"sync"
)

type Batcher struct {
	slowValue *slow.Value

	firstChanArr []chan chan interface{}

	isLoading bool
	mutex     sync.Mutex
}

func NewBatcher(v *slow.Value) *Batcher {
	return &Batcher{
		slowValue: v,
		firstChanArr:   make([]chan chan interface{}, 0),
	}
}

func (b *Batcher) addWaiter(resChan chan interface{}) {
	b.mutex.Lock()

	firstChan := make(chan chan interface{}, 1)
	b.firstChanArr = append(b.firstChanArr, firstChan)
	firstChan <- resChan

	if !b.isLoading {
		b.isLoading = true
		go func() {
			b.makeLoad()
		}()
	}
	b.mutex.Unlock()

}

func (b *Batcher) makeLoad() {
	b.mutex.Lock()

	secondChanArr := make([]chan interface{}, 0, len(b.firstChanArr))
	for _, firstChan := range b.firstChanArr {
		secondChanArr = append(secondChanArr, <-firstChan)
	}
	b.firstChanArr = b.firstChanArr[:0]

	b.mutex.Unlock()

	res := b.slowValue.Load()

	b.mutex.Lock()
	if len(b.firstChanArr) > 0 {
		go func() {
			b.makeLoad()
		}()
	} else {
		b.isLoading = false
	}
	b.mutex.Unlock()

	for _, secondChan := range secondChanArr {
		secondChan <- res
	}
}

func (b *Batcher) Load() (res interface{}) {
	resChan := make(chan interface{})
	b.addWaiter(resChan)
	res = <-resChan
	return
}