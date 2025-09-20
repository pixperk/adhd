package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pixperk/adhd/adhd"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from ADHD server! Time: %s\n", time.Now().Format(time.RFC3339))
	})

	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		fmt.Fprintf(w, "Slow response completed at %s\n", time.Now().Format(time.RFC3339))
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	shutdownCtx, shutdownCancel := adhd.WithCancel(adhd.Background())
	defer shutdownCancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Println("Shutdown signal received")
		shutdownCancel()
	}()

	go func() {
		log.Printf("Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	timeoutCtx, timeoutCancel := adhd.WithTimeout(shutdownCtx, 30*time.Second)
	defer timeoutCancel()

	winner := <-adhd.Race(shutdownCtx, timeoutCtx)

	switch winner.Error {
	case adhd.ErrCanceled:
		log.Println("Graceful shutdown initiated")
	case adhd.ErrDeadlineExceeded:
		log.Println("Shutdown timeout exceeded, forcing exit")
		os.Exit(1)
	}

	shutdownTimeout, cancel := adhd.WithTimeout(adhd.Background(), 10*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- server.Shutdown(adhd.Background())
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Printf("Server shutdown error: %v", err)
		} else {
			log.Println("Server shutdown complete")
		}
	case <-shutdownTimeout.Done():
		log.Println("Shutdown timeout reached, forcing close")
		server.Close()
	}
}
