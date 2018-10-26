package app

import (
	"context"
	"os"
)

type Context struct {
	context.Context
	cancelFn context.CancelFunc
	exitCode int // code for os.Exit() to return
	args     []string
}

func newContext() *Context {
	ctx, cancelFn := context.WithCancel(context.Background())

	return &Context{
		Context:  ctx,
		cancelFn: cancelFn,
		args:     os.Args[1:],
	}
}

func (ctx *Context) Cancel(exitCode int) {
	ctx.exitCode = exitCode
	ctx.cancelFn()
}

func (ctx *Context) Args() []string {
	return ctx.args
}
