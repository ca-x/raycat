package tinypool

import "sync"

type Pool[T any] struct {
	p     *sync.Pool
	reset func(*T)
}

func New[T any](resetFunc func(*T)) *Pool[T] {
	return &Pool[T]{
		p: &sync.Pool{New: func() any {
			return new(T)
		}},
		reset: resetFunc,
	}
}

func (p *Pool[T]) Get() *T {
	t := p.p.Get().(*T)
	if p.reset != nil {
		p.reset(t)
	}
	return t
}

func (p *Pool[T]) Free(data *T) {
	p.p.Put(data)
}
