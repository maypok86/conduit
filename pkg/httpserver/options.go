package httpserver

import "time"

// Option is a functional option for configuring a Server.
type Option func(*Server)

// WithHost sets the host for the server.
func WithHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

// WithPort sets the port for the server.
func WithPort(port string) Option {
	return func(s *Server) {
		s.port = port
	}
}

// WithMaxHeaderBytes sets the max header bytes for the server.
func WithMaxHeaderBytes(maxHeaderBytes int) Option {
	const numberOfBytes = 1 << 20 // 1MB

	return func(s *Server) {
		s.maxHeaderBytes = maxHeaderBytes * numberOfBytes
	}
}

// WithReadTimeout sets the read timeout for the server.
func WithReadTimeout(readTimeout time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = readTimeout
	}
}

// WithReadHeaderTimeout sets the read header timeout for the server.
func WithReadHeaderTimeout(readHeaderTimeout time.Duration) Option {
	return func(s *Server) {
		s.readHeaderTimeout = readHeaderTimeout
	}
}

// WithWriteTimeout sets the write timeout for the server.
func WithWriteTimeout(writeTimeout time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = writeTimeout
	}
}
