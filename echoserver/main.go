package main

import (
	"log"
	"net"
	"sync"
	"time"
)

// An echo server that replays back what was send to it.
// note: not so proud of this code, but it's doing it's job.

var mu sync.Mutex
var count int64

func main() {
	l, err := net.Listen("tcp4", ":8887")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			mu.Lock()
			log.Printf("Current Count: %d", count)
			mu.Unlock()

			time.Sleep(5 * time.Second)
		}
	}()

	for {
		c, err := l.Accept()
		if err != nil {
			continue
		}

		time.Sleep(100 * time.Millisecond)

		go handle(c)

		mu.Lock()
		count++
		mu.Unlock()
	}
}

func handle(c net.Conn) {
	defer c.Close()

	defer func() {
		mu.Lock()
		count--
		mu.Unlock()

	}()

	b := make([]byte, 2048)
	for {
		n, err := c.Read(b)

		if err != nil {
			log.Printf("got error '%v' during reading", err)
			return
		}

		_, err = c.Write(b[0:n])
		if err != nil {
			log.Printf("got error '%v' during writing", err)
			return
		}
	}
}
