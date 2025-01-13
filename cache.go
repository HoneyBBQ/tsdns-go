package tsdns

import (
	"fmt"
	"github.com/honeybbq/tsdns-go/types"
	"time"
)

// loadCache loads records from repository into memory cache
func (s *Server) loadCache() error {
	records, err := s.repository.Find()
	if err != nil {
		return err
	}

	newCache := make(map[string]*types.Record)
	for _, r := range records {
		newCache[r.Domain] = r
	}

	s.mu.Lock()
	s.cache = newCache
	s.mu.Unlock()

	return nil
}

// cacheUpdater periodically updates the in-memory cache
// Runs every 30 seconds to ensure data consistency
func (s *Server) cacheUpdater() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.loadCache(); err != nil {
				fmt.Printf("Cache update error: %v\n", err)
			}
		case <-s.ctx.Done():
			return
		}
	}
}
