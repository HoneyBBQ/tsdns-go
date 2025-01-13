package types

import "time"

type Record struct {
	ID         int64
	InstanceID int64
	Domain     string
	Target     string
	Port       int32
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

// RecordRepository defines the interface for record storage
type RecordRepository interface {
	// Find retrieves all records
	Find() ([]*Record, error)

	// FindByDomain finds a record by domain name
	FindByDomain(domain string) (*Record, error)

	// Create creates a new record
	Create(record *Record) error

	// Delete removes a record
	Delete(domain string) error

	// DeleteByInstanceID removes all records for a specific instance
	DeleteByInstanceID(instanceID int64) error

	// Close closes the storage connection
	Close() error
}
