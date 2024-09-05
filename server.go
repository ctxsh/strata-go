package strata

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type TLSOpts struct {
	// CertFile is the path to the file containing the SSL certificate or
	// certificate bundle.
	CertFile string
	// Keyfile is the path containing the certificate key.
	KeyFile string
	// InsecureSkipVerify controls whether a client verifies the server's
	// certificate chain and host name.
	InsecureSkipVerify bool
	// MinVersion contains the minimum TLS version that is acceptable.  By
	// default TLS 1.3 is used.
	MinVersion uint16
}

type ServerOpts struct {
	// BindAddr is the address the promethus collector will listen on for
	// connections.
	BindAddr string
	// BaseContext
	// Path is the path used by the HTTP server.
	Path string
	// Port is the path used by the HTTP server.
	Port int
	// TLS
	TLS *TLSOpts
	// TerminationGracePeriod is the amount of time that the server will wait
	// before stopping the HTTP server.  This grace period allows any prometheus
	// scrapers time to scrape.
	TerminationGracePeriod time.Duration
}

type Server struct {
	bindAddr               string
	logger                 Logger
	path                   string
	port                   int
	stopChan               chan struct{}
	stopOnce               sync.Once
	tlsCertFile            string
	tlsKeyFile             string
	tlsInsecureSkipVerify  bool
	tlsMinVersion          uint16
	terminationGracePeriod time.Duration
}

func newServer(opts ServerOpts) *Server {
	opts = defaultedServer(opts)
	return &Server{
		bindAddr:               opts.BindAddr,
		logger:                 logr.New(nil),
		path:                   opts.Path,
		port:                   opts.Port,
		tlsCertFile:            opts.TLS.CertFile,
		tlsKeyFile:             opts.TLS.KeyFile,
		tlsInsecureSkipVerify:  opts.TLS.InsecureSkipVerify,
		tlsMinVersion:          opts.TLS.MinVersion,
		terminationGracePeriod: opts.TerminationGracePeriod,
		stopChan:               make(chan struct{}, 1),
	}
}

// Start creates a new http server which listens on the TCP address addr
// and port.
func (s *Server) Start(ctx context.Context, reg *prometheus.Registry) error {
	mux := http.NewServeMux()
	mux.Handle(s.path, promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		Timeout: DefaultTimeout,
	}))

	server := &http.Server{
		Addr:        fmt.Sprintf("%s:%d", s.bindAddr, s.port),
		ReadTimeout: DefaultTimeout,
		Handler:     mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		<-ctx.Done()
		time.Sleep(s.terminationGracePeriod)
		if err := shutdown(server); err != nil {
			s.logger.Error(err, "shutting down prometheus collector endpoint")
		}
	}()

	if s.tlsCertFile != "" && s.tlsKeyFile != "" {
		server.TLSConfig = &tls.Config{
			MinVersion:         s.tlsMinVersion,
			InsecureSkipVerify: s.tlsInsecureSkipVerify, // nolint:gosec
		}

		s.logger.Info("starting prometheus collector endpoint", "tls", true, "config", s.config())
		return server.ListenAndServeTLS(s.tlsCertFile, s.tlsKeyFile)
	}

	s.logger.Info("starting prometheus collector endpoint", "tls", false, "config", s.config())
	return server.ListenAndServe()
}

// Stop closes the stop channel which initiates the shutdown of the HTTP server.
func (s *Server) Stop() {
	s.logger.Info("shutting down prometheus collector endpoint", "gracePeriod", s.terminationGracePeriod)
	s.stopOnce.Do(func() {
		time.Sleep(s.terminationGracePeriod)
		close(s.stopChan)
	})
}

// WithLogger defines the logger that will be used with the server.
func (s *Server) WithLogger(logger Logger) *Server {
	s.logger = logger
	return s
}

func (s *Server) config() map[string]any {
	return map[string]any{
		"addr":                          s.bindAddr,
		"path":                          s.path,
		"port":                          s.port,
		"terminationGracePeriodSeconds": s.terminationGracePeriod / time.Second,
		"tls": map[string]any{
			"certFile":           s.tlsCertFile,
			"keyFile":            s.tlsKeyFile,
			"insecureSkipVerify": s.tlsInsecureSkipVerify,
		},
	}
}

func shutdown(server *http.Server) error {
	toCtx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return server.Shutdown(toCtx)
}

func defaultedServer(opts ServerOpts) ServerOpts {
	if opts.BindAddr == "" {
		opts.BindAddr = "0.0.0.0"
	}

	if opts.Path == "" {
		opts.Path = "/metrics"
	}

	if opts.Port == 0 {
		opts.Port = 9090
	}

	opts.TLS = defaultedTLS(opts.TLS)

	return opts
}

func defaultedTLS(opts *TLSOpts) *TLSOpts {
	if opts == nil {
		opts = &TLSOpts{}
	}

	if opts.MinVersion == 0 {
		opts.MinVersion = tls.VersionTLS13
	}

	return opts
}
