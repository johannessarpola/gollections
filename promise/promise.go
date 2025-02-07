package promise

import (
	"context"
	"errors"

	"github.com/johannessarpola/gollections/result"
)

type Promise[T any] chan result.Result[T]

func New[T any]() Promise[T] {
	c := make(chan result.Result[T], 1)
	return c
}

func (p Promise[T]) Resolve(ctx context.Context, value T) Promise[T] {
	select {
	case <-ctx.Done():
		return p
	case p <- result.NewOk(value):
		return p
	}
}

func (p Promise[T]) Reject(ctx context.Context, err error) Promise[T] {
	select {
	case <-ctx.Done():
		return p
	case p <- result.NewErr[T](err):
		return p
	}
}

func (p Promise[T]) getWithinContext(ctx context.Context) result.Result[T] {
	select {
	case <-ctx.Done():
		return result.NewErr[T](errors.New("timeout exceeded"))
	case v := <-p:
		return v
	}
}

func (p Promise[T]) Then(ctx context.Context, f func(context.Context, T) result.Result[T]) Promise[T] {
	r := p.getWithinContext(ctx)
	if r.OK() {
		v := r.Value()
		transform := make(chan result.Result[T], 1)
		go func() {
			select {
			case <-ctx.Done():
				transform <- result.NewErr[T](errors.New("timeout exceeded"))
			default:
				res := f(ctx, v)
				select {
				case <-ctx.Done():
					transform <- result.NewErr[T](errors.New("timeout exceeded"))
				case transform <- res:
				}
			}
		}()
		return transform
	}
	return p.Reject(ctx, r.Err())
}

func (p Promise[T]) Wait() result.Result[T] {
	return <-p
}
