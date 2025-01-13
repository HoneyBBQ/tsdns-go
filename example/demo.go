package main

import (
	"github.com/honeybbq/tsdns-go"
	"github.com/honeybbq/tsdns-go/repository/postgres"
	"github.com/honeybbq/tsdns-go/repository/postgres/migrations"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	dsn := os.Getenv("PG_DSN")
	// use "migrate" as the first argument to run migrations
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		migrate(dsn)
	}

	log.Print("Starting tsdns-go demo...")

	// Create Postgres Repository
	repo := postgres.MustNewRepository(dsn)

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

// Migrate the database schema and exit
func migrate(dsn string) {
	// Migrate the schema
	err := migrations.AutoMigrate(dsn)
	if err != nil {
		log.Fatalf("Failed to migrate schema: %v", err)
	}
	log.Printf("Successfully migrated schema, please remove the 'migrate' argument to start the server, exiting...")
	os.Exit(0)
}
