package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"syscall"
	"time"

	"github.com/514366607/reload"
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

	err = server.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}

func (*MyHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "URL"+r.URL.String())
}
