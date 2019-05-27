package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type RandServer struct {
	srv  *http.Server
	done chan struct{}
}

func NewServer(port uint) *RandServer {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port)}
	return &RandServer{
		srv:  srv,
		done: make(chan struct{}),
	}
}

func (b *RandServer) Start(routes map[string]http.HandlerFunc) <-chan struct{} {
	for route, handler := range routes {
		http.HandleFunc(route, handler)
	}

	go func() {
		if err := b.srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()
	log.Printf("Server started")

	return b.done
}

func (b *RandServer) Stop() {
	defer close(b.done)
	log.Printf("Server shutting down")
	if err := b.srv.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
