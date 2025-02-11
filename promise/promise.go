package promise

import (
	"context"
	"errors"
	"sync"

	"github.com/johannessarpola/gollections/result"
)

type Promise[T any] chan result.Result[T]

func New[T any]() Promise[T] {
	c := make(chan result.Result[T], 1)
	return c
}

func Resolve[T any](value T) Promise[T] {
	p := New[T]()
	p <- result.NewOk(value)
	return p
}

func Reject[T any](err error) Promise[T] {
	p := New[T]()
	p <- result.NewErr[T](err)
	return p
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

// Catch executes a function if the promise results in an error.
func (p Promise[T]) Catch(fn func(error)) Promise[T] {
	go func() {
		res := <-p
		if res.IsErr() {
			fn(res.Err())
		}
	}()
	return p
}

func (p Promise[T]) Wait() result.Result[T] {
	return <-p
}

func All[T any](ctx context.Context, p ...Promise[T]) Promise[[]result.Result[T]] {
	ap := New[[]result.Result[T]]()
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(len(p))
		rss := make([]result.Result[T], len(p))
		for i, promise := range p {
			go func(p Promise[T], i int) {
				defer wg.Done()
				rs := promise.Wait()
				rss[i] = rs
			}(promise, i)
		}
		wg.Wait()
		ap.Resolve(ctx, rss)
	}()

	return ap
}
