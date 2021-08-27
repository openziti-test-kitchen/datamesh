package main

import (
	"fmt"
	"time"
)

var readQueue chan int

func main() {
	readQueue = make(chan int, 3)
	go reader()
	go writer()
	for {
		time.Sleep(30 * time.Minute)
	}
}

func reader() {
	for {
		time.Sleep(200 * time.Millisecond)
		select {
		case in := <-readQueue:
			fmt.Printf("<%d> ", in)
		}
	}
}

func writer() {
	i := 0
	for {
		select {
		case readQueue <- i:
			//

		default:
			fmt.Printf("_%d_ ", i)
		}
		i++
		time.Sleep(50 * time.Millisecond)
	}
}