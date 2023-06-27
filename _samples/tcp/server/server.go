package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"
	"syscall"

	"github.com/514366607/reload"
)

var (
	port     int
	isAccept int32 = 1
)

func main() {
	flag.IntVar(&port, "p", 8888, `端口`)
	flag.Parse()

	log.Printf("Actual pid is %d\n", syscall.Getpid())

	listener, err := reload.GetListener(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println(err)
	}

	var s = reload.NewServiceWith(listener, reload.WithDefaultHandle(), reload.WithHandleFunc(syscall.SIGUSR1, func(s reload.Service) {
		if err := s.Reload(); err != nil {
			s.Logger().Error(err)
		}
		log.Print("INlasdkjflaksjdflkasjdlkfajsldkfjaslkdfjaskldfjaslkdf\n\n\n\n\n\n\n\n\n\n")
		atomic.StoreInt32(&isAccept, 0)
	}))
	log.Printf("isChild : %v ,listener: %v\n", s.IsChild(), listener)

	go func() {
		defer listener.Close()

		for atomic.LoadInt32(&isAccept) == 1 {
			conn, err := listener.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			s.Add(1)
			log.Println("Accept ", conn.RemoteAddr())
			go recvConnMsg(conn, s)
		}
	}()

	s.Start()
}

func recvConnMsg(conn net.Conn, s reload.Service) {
	//  var buf [4096]byte
	buf := make([]byte, 4096)

	defer conn.Close()
	defer s.Done()

	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			//连接结束
			return
		} else if err != nil {
			log.Println(err)
			return
		}

		var recv = string(buf[0:n])
		recv = fmt.Sprintf(" pid %d Return: %s ", syscall.Getpid(), recv)
		log.Printf("Rev Data : %v", recv)

		conn.Write([]byte(recv))
	}
}
