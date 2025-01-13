package main

import (
	"github.com/honeybbq/tsdns-go"
	"github.com/honeybbq/tsdns-go/repository/postgres"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Print("Starting tsdns-go demo...")
	// Create Postgres Repository
	repo := postgres.MustNewRepository(os.Getenv("PG_DSN"))

	// Create TSDNS server
	s := tsdns.NewServer("0.0.0.0").
		WithRepository(repo).
		MustBuild()
	defer s.Close()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server
	go func() {
		if err := s.Start(); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down tsdns-go demo...")
}
