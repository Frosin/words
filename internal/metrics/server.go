package metrics

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultMetricsPort = "9091"
	defaultMetricsPath = "/metrics"
)

type Server struct {
	name    string
	addr    string
	server  *http.Server
	handler http.Handler
	mux     *http.ServeMux
	started chan struct{}
	errors  chan error

	listenAddr atomic.Value
}

func NewServer(addr, name string) *Server {
	server := &Server{
		name:    name,
		addr:    addr,
		mux:     http.NewServeMux(),
		started: make(chan struct{}, 1),
		errors:  make(chan error, 1),
	}

	return server
}

func (h *Server) Handle(path string, handler http.Handler) {
	h.mux.Handle(path, handler)
}

// ListenAndServe serves requests to server on the given net.Listener.
func (h *Server) Serve(ln net.Listener) error {
	if ln == nil {
		err := fmt.Errorf("net.Listener is nil")
		h.errors <- err
		return err
	}

	h.listenAddr.Store(ln.Addr())

	handler := h.handler
	if handler == nil {
		handler = h.mux
	}

	h.server = &http.Server{
		Addr:              h.addr,
		Handler:           handler,
		ReadHeaderTimeout: time.Second * 30,
		IdleTimeout:       time.Second * 120, // keep-alive
		ReadTimeout:       time.Minute * 5,
		WriteTimeout:      time.Minute * 5,
		MaxHeaderBytes:    1 << 20,
	}

	h.started <- struct{}{}

	log.Println(h.name, "server started")

	err := h.server.Serve(ln)
	if err == nil || errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	h.errors <- err
	return err
}

func (h *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", h.addr)
	if err != nil {
		h.errors <- err
		return err
	}

	return h.Serve(ln)
}

func RunMetrics() {
	m := newMetricsServer(":"+defaultMetricsPort, defaultMetricsPath)

	go func() {
		if err := m.ListenAndServe(); err != nil {
			log.Fatal("Unable to listen for metrics:", err)
		}
	}()

	log.Println("Metrics init and run")
}

func newMetricsServer(addr, path string) *Server {
	server := NewServer(addr, "Metrics")
	server.Handle(path, promhttp.Handler())

	return server
}
