package tsdns

import (
	"fmt"
	"net"
	"strings"
)

// handleQuery processes incoming DNS queries
// It looks up the domain in the cache and returns the corresponding record
// If no record is found, returns "404"
func (s *Server) handleQuery(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}

	domain := strings.TrimSpace(string(buf[:n]))
	if domain == "" {
		return
	}
	s.logger.Debug("Query received: %s\n", domain)

	// check cache
	s.mu.RLock()
	record, exists := s.cache[domain]
	s.mu.RUnlock()

	// record found
	if exists {
		response := record.Target
		if record.Port != 0 {
			response = fmt.Sprintf("%s:%d", record.Target, record.Port)
		}
		conn.Write([]byte(response))
		s.logger.Debug("Record found: %s -> %s\n", domain, response)
		return
	}

	// record not found
	s.logger.Debug("Record not found: %s\n", domain)
	conn.Write([]byte("404\n"))
}
