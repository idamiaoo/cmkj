package main

import (
	"fmt"
	"sync"
	"time"
)

type uids struct {
	id    int
	mutex sync.Mutex
}

func (u *uids) nextid() int {
	u.mutex.Lock()
	u.mutex.Unlock()
	u.id++
	return u.id
}

func newgen() chan int {
	var i int
	ch := make(chan int, 10)
	go func() {
		for {
			i++
			ch <- i
		}
	}()
	return ch
}

func main() {
	uids := new(uids)
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println(uids.nextid())
		}()
	}
	<-time.After(time.Duration(2) * time.Second)
	uidc := newgen()
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println(<-uidc)
		}()
	}
	<-time.After(time.Duration(2) * time.Second)
}
