package http

import (
	"crypto/tls"
	"gitlab.petrichor.io/sandstone/v1/concurrency/counter"
	"gitlab.petrichor.io/sandstone/v1/concurrency/trigger"
	"gitlab.petrichor.io/sandstone/v1/log"
	"gitlab.petrichor.io/sandstone/v1/errors/stackerr"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	ln                net.Listener
	lg                log.Logger
	srv               *http.Server
	identifyFn        IdentifyFn
	started           counter.Counter
	stopCh            *trigger.Trigger
	doneCh            *trigger.Trigger
}

type IdentifyFn func(ip string) (identifier string)

func NewServer(lg log.Logger) *Server {
	return &Server{
		srv:    new(http.Server),
		stopCh: trigger.New(),
		doneCh: trigger.New(),
		lg:     lg,
	}
}

func (s *Server) SetIdentifyFn(fn IdentifyFn) {
	s.identifyFn = fn
}

func (s *Server) SetConnStateFn(fn func(net.Conn, http.ConnState)) (error) {
	if !s.started.Zero() {
		return stackerr.New("Server already listening, cannot set conn state")
	}

	s.srv.ConnState = fn

	return nil
}

func (s *Server) Serve(ln net.Listener, handle http.Handler) error {
	return s.listen(ln, handle)
}

func (s *Server) Done() <-chan struct{} {
	return s.doneCh.Triggered()
}

func (s *Server) ServeTLSFromFiles(ln net.Listener, handle http.Handler, certFile, keyFile string) (err error) {
	cfg := new(tls.Config)
	if crt, err := tls.LoadX509KeyPair(certFile, keyFile); err != nil {
		return stackerr.Wrap(err, "Could not load certificate")
	} else {
		cfg.Certificates = []tls.Certificate{crt}
	}

	return s.listen(tls.NewListener(ln, cfg), handle)
}

func (s *Server) ServeTLS(ln net.Listener, handle http.Handler, certPEMBlock []byte, keyPEMBlock []byte) (err error) {
	cfg := new(tls.Config)

	if crt, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock); err != nil {
		return stackerr.Wrap(err, "Could not load certificate")
	} else {
		cfg.Certificates = []tls.Certificate{crt}
	}

	return s.listen(tls.NewListener(ln, cfg), handle)
}

func (s *Server) listen(ln net.Listener, h http.Handler) error {
	if !s.started.Zero() {
		return stackerr.New("Server already listening")
	}

	s.ln = ln
	s.started.Hit()

	var wg sync.WaitGroup

	s.srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-s.stopCh.Triggered():
			return
		default:
			wg.Add(1)
			defer wg.Done()
		}

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		identifier := ip

		if s.identifyFn != nil {
			identifier = s.identifyFn(ip)
		}

		s.lg.Debugf("Serving request %-5s %-35s (%s)", r.Method, r.RequestURI, identifier)

		h.ServeHTTP(w, r)
	})

	s.lg.Successf("HTTP Server listening on %s", ln.Addr())

	go s.serve()

	go func() {
		for {
			select {
			case <-s.stopCh.Triggered():
				wg.Wait()
				s.doneCh.Trigger()
				return
			}
		}
	}()

	return nil
}

func (s *Server) serve() {
listen:
	for {
		select {
		case <-s.stopCh.Triggered():
			// Things are stopping
			break listen
		default:
			// continue
		}

		if err := s.srv.Serve(s.ln); err != nil && s.lg != nil {
			s.lg.Errorf(err.Error())
		}
	}
}

func (s *Server) Stop() (err error) {
	if s.started.Zero() {
		return nil
	}

	s.stopCh.Trigger()

	if err := s.ln.Close(); err != nil {
		return err
	}

	select {
	case <-s.doneCh.Triggered():

	case <-time.After(time.Second * 15):
		return stackerr.New("Timeout while shutting down")
	}

	return nil
}
