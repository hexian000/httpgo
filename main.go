package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	listen := flag.String("l", ":8080", "local listen port")
	dir := flag.String("d", ".", "serve path")
	help := flag.Bool("h", false, "show usage")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	server := &http.Server{Addr: *listen, Handler: http.FileServer(http.Dir(*dir))}

	go func() {
		c := make(chan os.Signal, 8)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		log.Println("Caught signal:", <-c)
		ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalln(err)
		}
	}()

	log.Println("server starting")
	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Println(err)
		}
	}
	log.Println("server stopping gracefully")
}
