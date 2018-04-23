package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"reload"
	"syscall"
	"time"
)

// MyHandle 头
type MyHandle struct{}

func main() {
	flag.Parse()

	var server = http.Server{
		Addr:        ":8888",
		Handler:     &MyHandle{},
		ReadTimeout: 6 * time.Second,
	}

	log.Printf("Actual pid is %d\n", syscall.Getpid())

	listener, err := reload.GetListener(server.Addr)
	if err != nil {
		log.Println(err)
	}

	var s = reload.NewService(listener)
	log.Printf("isChild : %v ,listener: %v\n", s.IsChild(), listener)

	go func() {
		err = server.Serve(listener)
		if err != nil {
			log.Println(err)
		}
	}()

	s.Start()
}

func (*MyHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "URL"+r.URL.String())
}
