package tsdns

import (
	"github.com/honeybbq/tsdns-go/types"
)

// AddRecord adds a new DNS record to the system
// Updates both repository and cache immediately
func (s *Server) AddRecord(domain, target string, port int32) error {
	record := &types.Record{
		Domain: domain,
		Target: target,
		Port:   port,
	}

	err := s.repository.Create(record)
	if err != nil {
		return err
	}

	// update cache
	return s.loadCache()
}

// RemoveRecord deletes a DNS record by domain name
// Updates both repository and cache immediately
func (s *Server) RemoveRecord(domain string) error {
	err := s.repository.Delete(domain)
	if err != nil {
		return err
	}

	// 立即更新缓存
	return s.loadCache()
}

// RemoveInstanceRecords removes all records associated with an instance
// Updates both repository and cache immediately
func (s *Server) RemoveInstanceRecords(instanceID int64) error {
	err := s.repository.DeleteByInstanceID(instanceID)
	if err != nil {
		return err
	}

	// update cache
	return s.loadCache()
}
