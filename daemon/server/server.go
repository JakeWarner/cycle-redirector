package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"gitlab.petrichor.io/sandstone/v1/errors/stackerr"
	"gitlab.petrichor.io/sandstone/v1/log"

	httpHelper "gitlab.petrichor.io/sandstone/v1/net/http"
)

var (
	servers      []*httpHelper.Server
	redirectUrl  string
	redirectHost string
)

var redirectFn = http.HandlerFunc(
	func(w http.ResponseWriter, req *http.Request) {
		if redirectUrl != "" {
			http.Redirect(w, req, redirectUrl, 301)
		} else {
			http.Redirect(w, req, fmt.Sprintf("%s%s", redirectHost, req.URL.EscapedPath()), 301)
		}
	},
)

const TLS_PATH = "/var/run/cycle/tls"

func Listen(ctx context.Context) (err error) {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:

	}

	redirectUrl = os.Getenv("REDIRECT_URL")
	redirectHost = os.Getenv("REDIRECT_HOST")

	if redirectUrl == "" && redirectHost == "" {
		return stackerr.New("Neither REDIRECT_URL nor REDIRECT_HOST environment variables have been configured")
	}

	bundlePath := filepath.Join(TLS_PATH, "current.bundle")
	keyPath := filepath.Join(TLS_PATH, "current.key")

	if _, err := os.Stat(bundlePath); err == nil {
		{
			tcpAddr, err := net.ResolveTCPAddr("", ":443")
			if err != nil {
				return stackerr.Wrap(err, "Could not resolve tcp address")
			}

			ln, err := net.ListenTCP("tcp", tcpAddr)
			if err != nil {
				return stackerr.Wrap(err, "Could not listen on tcp port")
			}

			s := httpHelper.NewServer(log.New("HTTPS"))

			if err := s.ServeTLSFromFiles(ln, redirectFn, bundlePath, keyPath); err != nil {
				return err
			}

			servers = append(servers, s)
		}

		{
			s := httpHelper.NewServer(log.New("HTTP"))

			tcpAddr, err := net.ResolveTCPAddr("", ":80")
			if err != nil {
				return stackerr.Wrap(err, "Could not resolve tcp address")
			}

			ln, err := net.ListenTCP("tcp", tcpAddr)
			if err != nil {
				return stackerr.Wrap(err, "Could not listen on tcp port")
			}

			if err := s.Serve(ln, redirectFn); err != nil {
				return err
			}

			servers = append(servers, s)
		}
	} else {
		s := httpHelper.NewServer(log.New("HTTP"))

		tcpAddr, err := net.ResolveTCPAddr("", ":80")
		if err != nil {
			return stackerr.Wrap(err, "Could not resolve tcp address")
		}

		ln, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			return stackerr.Wrap(err, "Could not listen on tcp port")
		}

		if err := s.Serve(ln, redirectFn); err != nil {
			return err
		}

		servers = append(servers, s)
	}

	return nil
}

func Stop() (err error) {
	if len(servers) == 0 {
		return nil
	}

	for _, v := range servers {
		if err := v.Stop(); err != nil {
			return err
		}
	}

	return nil
}
