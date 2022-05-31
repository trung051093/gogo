package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 1)
	numbers := []int{1, 2, 3, 4, 5}
	go func() {
		for _, i := range numbers {
			fmt.Println("push :", i)
			ch <- i
			time.Sleep(time.Second)
		}
	}()

	for {
		select {
		case i := <-ch:
			println("get: ", i)
		}
	}

}
