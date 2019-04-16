package server

import (
	"context"
	"fmt"
	"gitlab.petrichor.io/sandstone/v1/errors/stackerr"
	"gitlab.petrichor.io/sandstone/v1/log"
	"net"
	"net/http"
	"os"

	httpHelper "gitlab.petrichor.io/sandstone/v1/net/http"
)

var (
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

	s := httpHelper.NewServer(log.New("HTTP"))

	tcpAddr, err := net.ResolveTCPAddr("", "::80")
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

	return nil
}
