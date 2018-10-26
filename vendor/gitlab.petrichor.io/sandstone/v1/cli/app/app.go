package app

import (
	"context"
	"gitlab.petrichor.io/sandstone/v1/log"
	"gitlab.petrichor.io/sandstone/v1/errors/stackerr"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
	"gitlab.petrichor.io/sandstone/v1/concurrency/counter"
)

type StartFn func(ctx context.Context) error
type StopFn func() error

type app struct {
	init        sync.Once
	About       About
	stopTimeout time.Duration
	stopFn      StopFn
	started     counter.Counter
	commands    Commands
	Log         *log.Log
	ctx         *Context
	mutex       sync.Mutex
}

type About struct {
	Name        string
	Version     string
	Description string
}

type Opts struct {
	StopTimeout time.Duration
	StopFn      StopFn
}

var (
	a = &app{
		Log:         log.Default,
		stopTimeout: time.Minute,
		ctx:         newContext(),
	}
	ranCh chan struct{}
)

func New(about About, opts *Opts) *app {
	a.init.Do(func() {
		a.About = about

		if opts != nil {
			if opts.StopTimeout > 0 {
				a.stopTimeout = opts.StopTimeout
			}

			if opts.StopFn != nil {
				a.stopFn = opts.StopFn
			}
		}
	})

	return a
}

func Current() *app {
	return a
}

func Ctx() context.Context {
	return a.ctx
}

func Start(startFn StartFn) (err error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if !a.started.Zero() {
		return stackerr.New("App already started")
	}

	defer a.started.Hit()

	// Watch for interrupts and kill signals
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		select {
		case s := <-sig:
			a.Log.Warnf(s.String() + " Signal Caught")
			a.ctx.Cancel(0)
		}
	}()

	go func() {
		select {
		case <-a.ctx.Done():
			if ranCh != nil {
				select {
				case <-ranCh:

				case <-time.After(time.Second * 30):
				}
			}

			if a.stopFn != nil {
				stopCh := make(chan struct{})
				go func() {
					a.Log.Debugf("Stopping...")
					if err := a.stopFn(); err != nil {
						a.Log.Errorf("Error while stopping: %s", err.Error())
					}
					close(stopCh)
				}()

				select {
				case <-stopCh:
				case <-time.After(a.stopTimeout):
					a.Log.Errorf("Timeout while stopping")
					os.Exit(1)
				}
			}
			os.Exit(a.ctx.exitCode)
		}
	}()

	if startFn != nil {
		if e := startFn(a.ctx); e != nil {
			return stackerr.Wrap(e, "Cannot start")
		}
	}

	return nil
}

func Register(cmd Command) {
	a.commands = append(a.commands, cmd)
}

func ExitError(err error) {
	a.Log.Errorf(err.Error())
	Stop(1)
}

func Stop(code int) {
	a.ctx.Cancel(code)
}

func Run() (err error) {
	args := os.Args[1:]

	ranCh = make(chan struct{})
	defer close(ranCh)

	if len(args) == 0 {
		return nil
	}

	for _, v := range a.commands {
		if v.Keyword == strings.ToLower(args[0]) {
			a.ctx.args = a.ctx.args[1:]
			if err := v.Fn(Ctx()); err != nil {
				return err
			}

			return nil
		}
	}

	return stackerr.Newf("Command %s not found", args[0])
}

func SetStopFn(fn StopFn) {
	a.stopFn = fn
}
