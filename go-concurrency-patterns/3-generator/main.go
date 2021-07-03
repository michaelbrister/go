package main

import (
	"fmt"
	"math/rand"
	"time"
)

func boring(msg string) <-chan string {
	c := make(chan string)
	go func() { // we launch goroutine inside a function.
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}

	}()
	return c // return a channel to caller.

}

func main() {
	joe := boring("Joe")
	ahn := boring("Ahn")

	for i := 0; i < 5; i++ {
		fmt.Println(<-joe)
		fmt.Println(<-ahn)
	}
	fmt.Println("You're boring. I'm leaving")
}
