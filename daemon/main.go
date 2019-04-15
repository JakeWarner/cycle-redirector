package main

import (
	"context"
	"time"

	"github.com/jakewarner/cycle-redirector/daemon/server"
	"gitlab.petrichor.io/sandstone/v1/builds"
	"gitlab.petrichor.io/sandstone/v1/cli/app"
	"gitlab.petrichor.io/sandstone/v1/concurrency/events"
	"gitlab.petrichor.io/sandstone/v1/errors/stackerr"
	"gitlab.petrichor.io/sandstone/v1/log"
)

var version = "2019.04.14"

func main() {
	builds.SetLog(log.Default)

	a := app.New(app.About{
		Name:        "gohttpd",
		Description: "Simple HTTP Redirection daemon",
		Version:     version,
	}, &app.Opts{
		StopTimeout: time.Minute,
		StopFn:      stop,
	})

	log.DefaultOptions.CaptionPadding = 30

	a.Log.Noticef("%s (%s)", a.About.Name, a.About.Version)

	if err := app.Start(start); err != nil {
		app.ExitError(stackerr.FlattenAll(err))
	}

	if err := app.Run(); err != nil {
		app.ExitError(stackerr.FlattenAll(err))
	}

	select {} // Wait forever
}

func start(ctx context.Context) error {
	if err := events.DefaultDispatcher.Fire("start", nil); err != nil {
		return err
	}

	if err := server.Listen(ctx); err != nil {
		return err
	}

	if err := events.DefaultDispatcher.Fire("started", nil); err != nil {
		return err
	}

	return nil
}

func stop() error {
	if err := events.DefaultDispatcher.Fire("stop", nil); err != nil {
		log.Default.Errorf("Error while stopping: %s", err.Error())
	}

	return nil
}
