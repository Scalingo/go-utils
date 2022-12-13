package errors

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/pkg/errors"
)

type ErrCtx struct {
	ctx context.Context
	err error
}

func (err ErrCtx) Error() string {
	return err.err.Error()
}

func (err ErrCtx) Ctx() context.Context {
	return err.ctx
}

func New(ctx context.Context, message string) error {
	return ErrCtx{ctx: ctx, err: errgo.New(message)}
}

func Newf(ctx context.Context, format string, args ...interface{}) error {
	return ErrCtx{ctx: ctx, err: errgo.Newf(format, args...)}
}

func NoteMask(ctx context.Context, err error, message string) error {
	return ErrCtx{ctx: ctx, err: errgo.NoteMask(err, message)}
}

func Notef(ctx context.Context, err error, format string, args ...interface{}) error {
	return ErrCtx{ctx: ctx, err: errgo.Notef(err, format, args...)}
}

func Wrap(ctx context.Context, err error, message string) error {
	return ErrCtx{ctx: ctx, err: errors.Wrap(err, message)}
}

func Wrapf(ctx context.Context, err error, format string, args ...interface{}) error {
	return ErrCtx{ctx: ctx, err: errors.Wrapf(err, format, args...)}
}

func Errorf(ctx context.Context, format string, args ...interface{}) error {
	return ErrCtx{ctx: ctx, err: errors.Errorf(format, args...)}
}
