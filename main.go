package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/simonlim94/posts-comments-demo/controller"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "0"
	}

	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/top-posts", controller.GetTopPostsByComments)
	http.HandleFunc("/filtered-comments", controller.GetCommentsByFilter)

	go func() {
		log.Println("Listening at:", "http://"+listener.Addr().String())
		if err := http.Serve(listener, nil); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	// Listen for shutdown with CTRL+C
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	<-signalCh
	log.Println("received shutdown signal, terminating...")
}
