package promise

import "github.com/johannessarpola/gollections/result"

type Promise[T any] chan result.Result[T]

func NewPromise[T any]() Promise[T] {
	return make(Promise[T])
}

func (p Promise[T]) Resolve(value T) {
	p <- result.NewOk(value)
}

func (p Promise[T]) Reject(err error) {
	p <- result.NewErr[T](err)
}
