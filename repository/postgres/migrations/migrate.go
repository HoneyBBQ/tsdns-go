package migrations

import (
	"embed"
	"fmt"
	"gorm.io/driver/postgres"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

//go:embed *.sql
var migrationFiles embed.FS

// Migration represents the migration history record
type Migration struct {
	ID        uint   `gorm:"primaryKey"`
	Version   string `gorm:"uniqueIndex"`
	AppliedAt int64
}

// AutoMigrate applies all pending migrations
//
// dsn is the PostgreSQL connection string
//
// path is the directory containing migration files
func AutoMigrate(dsn string) error {
	log.Printf("Migrating database schema using DSN: %s", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	// 1. Create migration history table
	if err = db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("create migration table: %w", err)
	}

	// 2. Get all migration files
	files, err := migrationFiles.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read migration files: %w", err)
	}

	// 3. Filter and sort migration files
	var upMigrations []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".up.sql") {
			upMigrations = append(upMigrations, f.Name())
		}
	}

	// 4. Execute migrations
	for _, filename := range upMigrations {
		log.Printf("Applying migration: %s", filename)
		version := strings.Split(filename, "_")[0]

		// Check if already applied
		var count int64
		db.Model(&Migration{}).Where("version = ?", version).Count(&count)
		if count > 0 {
			continue
		}

		// Read and execute migration file
		content, _err := migrationFiles.ReadFile(filename)
		if _err != nil {
			return fmt.Errorf("read migration file %s: %w", filename, err)
		}

		// Execute migration in transaction
		_err = db.Transaction(func(tx *gorm.DB) error {
			// Execute SQL
			if __err := tx.Exec(string(content)).Error; __err != nil {
				return __err
			}

			// Record migration history
			return tx.Create(&Migration{
				Version:   version,
				AppliedAt: time.Now().Unix(),
			}).Error
		})

		if _err != nil {
			return fmt.Errorf("execute migration %s: %w", filename, err)
		}
	}

	return nil
}
