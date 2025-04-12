package server

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	conf          *ServerConf
	log           *zap.Logger
	Handler       http.Handler
	Http          *http.Server
	Https         *http.Server
	httpListener  net.Listener
	httpsListener net.Listener
}

type ServerConf struct {
	Addr           string        `yaml:"addr"`
	SSLCert        string        `yaml:"ssl_cert"`
	SSLKey         string        `yaml:"ssl_key"`
	SSLAddr        string        `yaml:"ssl_addr"`
	KeepAlive      bool          `yaml:"keep_alive"`
	ReadTimeout    time.Duration `yaml:"read_timeout"`
	WriteTimeout   time.Duration `yaml:"write_timeout"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
	MaxHeaderBytes int           `yaml:"max_header_bytes"`
}

func NewServer(conf *ServerConf, log *zap.Logger, handler http.Handler) *Server {
	return &Server{conf: conf, log: log, Handler: handler}
}

func (s *Server) Run(errChan chan error) {
	// init 初始化
	if err := s.init(); err != nil {
		errChan <- err
		return
	}

	// start HTTP Server
	s.log.Info("The Server is listening", zap.String("http", s.conf.Addr))
	go func() {
		if err := s.Http.Serve(s.httpListener); err != nil && err != http.ErrServerClosed {
			errChan <- err
			return
		}
	}()

	// start HTTPS Server
	if s.conf.SSLCert != "" && s.conf.SSLKey != "" && s.Https != nil {
		s.log.Info("The Server is listening on", zap.String("https", s.conf.SSLAddr))
		go func() {
			if err := s.Https.ServeTLS(s.httpsListener,
				s.conf.SSLCert, s.conf.SSLKey); err != nil && err != http.ErrServerClosed {
				errChan <- err
				return
			}
		}()
	}
}

func (s *Server) Shutdown() {
	s.shutdownServer(s.Http)
	s.shutdownServer(s.Https)
}

func (s *Server) init() error {
	s.log.Info("Initialize Http API Server")

	// http
	s.Http = &http.Server{
		Addr:           s.conf.Addr,
		Handler:        s.Handler,
		ReadTimeout:    s.conf.ReadTimeout,
		WriteTimeout:   s.conf.WriteTimeout,
		IdleTimeout:    s.conf.IdleTimeout,
		MaxHeaderBytes: s.conf.MaxHeaderBytes,
	}
	s.Http.SetKeepAlivesEnabled(s.conf.KeepAlive)

	// https
	if s.conf.SSLCert != "" && s.conf.SSLKey != "" {
		s.Https = &http.Server{
			Addr:           s.conf.SSLAddr,
			Handler:        s.Handler,
      ReadTimeout:    s.conf.ReadTimeout,
      WriteTimeout:   s.conf.WriteTimeout,
      IdleTimeout:    s.conf.IdleTimeout,
      MaxHeaderBytes: s.conf.MaxHeaderBytes,
			TLSConfig: &tls.Config{
				// 优先使用服务器端的cipherSuite密码套件，确保安全
				PreferServerCipherSuites: true,
			},
		}
		s.Https.SetKeepAlivesEnabled(s.conf.KeepAlive)
	}

	// httpListener
	var err error
	s.httpListener, err = net.Listen("tcp", s.conf.Addr)
	if err != nil {
		s.log.Error("Listen error", zap.Error(err))
		return err
	}
	// httpsListener
	if s.conf.SSLCert != "" && s.conf.SSLKey != "" {
		s.httpsListener, err = net.Listen("tcp", s.conf.SSLAddr)
		if err != nil {
			s.log.Error("Listen error", zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *Server) shutdownServer(server *http.Server) {
	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			s.log.Error("Server Shutdown error", zap.Error(err))
		}
	}
}
