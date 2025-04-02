package stream

import (
	it "github.com/dzherb/go-itertools"
	"iter"
)

type Stream[V any] iter.Seq[V]

func FromIterator[V any](iter iter.Seq[V]) Stream[V] {
	return Stream[V](iter)
}

func FromElements[V any](elems ...V) Stream[V] {
	return FromIterator(it.FromElements(elems...))
}

func (s Stream[V]) Iterator() iter.Seq[V] {
	return iter.Seq[V](s)
}

func (s Stream[V]) Slice(start, stop, step int) Stream[V] {
	return FromIterator(it.Slice[V](s.Iterator(), start, stop, step))
}

func (s Stream[V]) Take(n int) Stream[V] {
	return FromIterator(it.Take[V](s.Iterator(), n))
}

func (s Stream[V]) TakeWhile(predicate func(V) bool) Stream[V] {
	return FromIterator(it.TakeWhile[V](s.Iterator(), predicate))
}

func (s Stream[V]) DropWhile(predicate func(V) bool) Stream[V] {
	return FromIterator(it.DropWhile[V](s.Iterator(), predicate))
}

func (s Stream[V]) Filter(predicate func(V) bool) Stream[V] {
	return FromIterator(it.Filter[V](s.Iterator(), predicate))
}

func (s Stream[V]) Map(mapper func(V) V) Stream[V] {
	return FromIterator(it.Map[V](s.Iterator(), mapper))
}

func MapStream[V, R any](s Stream[V], mapper func(V) R) Stream[R] {
	return FromIterator(it.Map(s.Iterator(), mapper))
}
