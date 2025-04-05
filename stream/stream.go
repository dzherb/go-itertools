package stream

import (
	it "github.com/dzherb/go-itertools"
	"iter"
	"slices"
)

type Stream[V any] iter.Seq[V]

func FromIterator[V any](iter iter.Seq[V]) Stream[V] {
	return Stream[V](iter)
}

func FromElements[V any](elems ...V) Stream[V] {
	return FromIterator(it.FromElements(elems...))
}

func FromChan[V any](ch <-chan V) Stream[V] {
	return FromIterator(it.FromChan(ch))
}

func (s Stream[V]) Iterator() iter.Seq[V] {
	return iter.Seq[V](s)
}

func (s Stream[V]) Collect() []V {
	return slices.Collect(s.Iterator())
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

// Map applies a transformation function to each element of a stream,
// producing a new stream of the result type.
//
// Since Go currently does not support type parameters in methods,
// we need a standalone generic function to allow
// mapping from one stream type to another.
func Map[V, R any](s Stream[V], mapper func(V) R) Stream[R] {
	return FromIterator(it.Map(s.Iterator(), mapper))
}

func (s Stream[V]) ForEach(consumer func(V)) {
	it.ForEach(s.Iterator(), consumer)
}
