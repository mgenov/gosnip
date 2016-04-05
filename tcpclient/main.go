// A tcpclient that generates ECHO traffic to the echo server.
// App is connecting n clients to the target server and sends ECHO messages and waits for replay from the server.

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

var server = flag.String("server", "127.0.0.1:2090", "the target server that will listen")
var clients = flag.Int("client", 100, "number of concurrent clients")

func main() {
	flag.Parse()

	for i := 0; i < *clients; i++ {
		go Client(i, *server)
		time.Sleep(5 * time.Millisecond)
	}
	select {}

}

func Client(id int, host string) {

	c, err := net.DialTimeout("tcp4", host, 3*time.Second)
	if err != nil {
		log.Println(err)
		return
	}
	buf := make([]byte, 2048)
	for {
		_, err := c.Write([]byte(fmt.Sprintf("ECHO: %d\n", id)))
		if err != nil {
			log.Printf("Write error: %v", err)
			break
		}

		_, err = c.Read(buf)
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		time.Sleep(30 * time.Second)
	}
	log.Println("Client terminted.")
	c.Close()
}
