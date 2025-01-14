package file

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/honeybbq/tsdns-go/types"
	"github.com/vmihailenco/msgpack/v5"
)

type repository struct {
	filePath string
	records  map[string]*types.Record
	mu       sync.RWMutex
}

// NewRepository creates a new file-based repository
//
// filePath: path to the binary file for storage
func NewRepository(filePath string) (types.RecordRepository, error) {
	repo := &repository{
		filePath: filePath,
		records:  make(map[string]*types.Record),
	}

	// Load existing records if file exists
	if err := repo.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load records: %v", err)
	}

	return repo, nil
}

// load reads records from file
func (f *repository) load() error {
	// Check if file exists
	_, err := os.Stat(f.filePath)
	if os.IsNotExist(err) {
		// Create empty file if it doesn't exist
		file, err := os.OpenFile(f.filePath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}
		file.Close()
		return nil // Return as there's nothing to load
	}

	// Read file content
	data, err := os.ReadFile(f.filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Only try to decode if file is not empty
	if len(data) > 0 {
		if err := msgpack.Unmarshal(data, &f.records); err != nil {
			return fmt.Errorf("failed to decode records: %v", err)
		}
	}

	return nil
}

// save writes records to file
func (f *repository) save() error {
	data, err := msgpack.Marshal(f.records)
	if err != nil {
		return fmt.Errorf("failed to encode records: %v", err)
	}

	if err := os.WriteFile(f.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// Find retrieves all records
func (f *repository) Find() ([]*types.Record, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	records := make([]*types.Record, 0, len(f.records))
	for _, record := range f.records {
		if record.DeletedAt == nil {
			records = append(records, record)
		}
	}
	return records, nil
}

// FindByDomain finds a record by domain name
func (f *repository) FindByDomain(domain string) (*types.Record, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	record, exists := f.records[domain]
	if !exists || record.DeletedAt != nil {
		return nil, fmt.Errorf("record not found")
	}
	return record, nil
}

// Create creates a new record
func (f *repository) Create(record *types.Record) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Set timestamps
	now := time.Now()
	record.CreatedAt = now
	record.UpdatedAt = now

	f.records[record.Domain] = record
	return f.save()
}

// Delete removes a record
func (f *repository) Delete(domain string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	record, exists := f.records[domain]
	if !exists {
		return fmt.Errorf("record not found")
	}

	now := time.Now()
	record.DeletedAt = &now
	record.UpdatedAt = now

	return f.save()
}

// DeleteByInstanceID removes all records for a specific instance
func (f *repository) DeleteByInstanceID(instanceID int64) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()
	for _, record := range f.records {
		if record.InstanceID == instanceID {
			record.DeletedAt = &now
			record.UpdatedAt = now
		}
	}

	return f.save()
}

// Close implements repository interface
func (f *repository) Close() error {
	return f.save()
}
