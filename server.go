package tsdns

import (
	"context"
	"fmt"
	"github.com/honeybbq/tsdns-go/types"
	"net"
	"sync"
)

// ServerBuilder represents a builder for TSDNS server
type ServerBuilder struct {
	server *Server
	err    error
}

// Server represents a TSDNS server instance
type Server struct {
	addr       string
	repository types.RecordRepository
	cache      map[string]*types.Record
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	logger     Logger
}

// NewServer creates a new TSDNS server builder
func NewServer(ip string) *ServerBuilder {
	ctx, cancel := context.WithCancel(context.Background())

	builder := &ServerBuilder{
		server: &Server{
			addr:   ip + ":41144",
			cache:  make(map[string]*types.Record),
			ctx:    ctx,
			cancel: cancel,
			logger: newStdLogger(), // Default logger
		},
	}

	// Validate IP address
	if ip != "0.0.0.0" && net.ParseIP(ip) == nil {
		builder.err = fmt.Errorf("invalid IP address")
	}

	return builder
}

// WithRepository sets the repository implementation
func (b *ServerBuilder) WithRepository(repo types.RecordRepository) *ServerBuilder {
	if b.err != nil {
		return b
	}
	b.server.repository = repo
	return b
}

// WithLogger sets the logger implementation
func (b *ServerBuilder) WithLogger(l Logger) *ServerBuilder {
	if b.err != nil {
		return b
	}
	b.server.logger = l
	return b
}

// Build creates and returns the server instance
func (b *ServerBuilder) Build() (*Server, error) {
	if b.err != nil {
		return nil, b.err
	}

	// Validate required fields
	if b.server.repository == nil {
		return nil, fmt.Errorf("repository is required")
	}

	// Start cache updater
	go b.server.cacheUpdater()

	return b.server, nil
}

// MustBuild creates and returns the server instance
func (b *ServerBuilder) MustBuild() *Server {
	server, err := b.Build()
	if err != nil {
		panic(err)
	}
	return server
}

// Start initializes and runs the TSDNS server
// It listens for incoming TCP connections and handles DNS queries
func (s *Server) Start() error {
	addr, err := net.ResolveTCPAddr("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("resolve address error: %v", err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen error: %v", err)
	}
	defer listener.Close()

	// load cache initially
	if err = s.loadCache(); err != nil {
		return fmt.Errorf("load cache error: %v", err)
	}

	s.logger.Info("TSDNS server started at %s\n", s.addr)
	// start handling queries
	for {
		conn, _err := listener.Accept()
		if _err != nil {
			s.logger.Error("accept error: %v\n", _err)
			continue
		}
		go s.handleQuery(conn)
	}
}

// Close shuts down the server and releases resources
func (s *Server) Close() error {
	s.logger.Info("Shutting down tsdns-go server...")
	s.cancel()
	return s.repository.Close()
}
