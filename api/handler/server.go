package handler

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	port     string
	handlers Handlers
}

func NewServer(port string, h Handlers) *Server {
	return &Server{
		port:     port,
		handlers: h,
	}
}

func (s *Server) ConfigureAndRun() {
	mux := http.NewServeMux()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	mux.HandleFunc("/hash", s.handlers.HashApi)
	mux.HandleFunc("/hash/", s.handlers.GetHashedPassword)

	httpServer := &http.Server{
		Addr:        s.port,
		Handler:     mux,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen failed: %s\n", err)
		}
	}()

	log.Printf("Server started and listening")

	<-done

	log.Printf("Server stopped. Shutting down...")
	defer func() {
		cancel()
	}()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed %s\n", err)
	}

	log.Printf("Server stopped.")

}
