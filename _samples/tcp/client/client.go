package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

var (
	message string
	port    int
)

func main() {
	flag.StringVar(&message, "s", "123", `要发送的内容`)
	flag.IntVar(&port, "p", 8888, `要发送的内容`)
	flag.Parse()

	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {

		time.Sleep(time.Millisecond * 200)

		conn.Write([]byte(message))

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Rev ERROR : %v", err)
			// return
		}

		buf = buf[0:n]
		log.Printf("Rev Data : %v", string(buf))
	}
}
