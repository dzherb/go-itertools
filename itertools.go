package itertools

import (
	"iter"
)

func Count(start, step int) iter.Seq[int] {
	return func(yield func(int) bool) {
		cur := start
		for {
			if !yield(cur) {
				return
			}
			cur += step
		}
	}
}

func Cycle[V any](iter iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for {
			for v := range iter {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Cycle2[K, V any](iter iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for {
			for k, v := range iter {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

func Repeat[V any](elem V, n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		for i := 0; i < n; i++ {
			if !yield(elem) {
				return
			}
		}
	}
}

func Chain[V any](iters ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, it := range iters {
			for v := range it {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Chain2[K, V any](iters ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, it := range iters {
			for k, v := range it {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

func Zip[K, V any](first iter.Seq[K], second iter.Seq[V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		ch1 := make(chan K)
		ch2 := make(chan V)
		done := make(chan struct{})

		go func() {
			defer close(ch1)
			first(func(k K) bool {
				select {
				case ch1 <- k:
					return true
				case <-done:
					return false
				}
			})
		}()

		go func() {
			defer close(ch2)
			second(func(v V) bool {
				select {
				case ch2 <- v:
					return true
				case <-done:
					return false
				}
			})
		}()

		for {
			var k K
			var v V
			var ok1, ok2 bool

			select {
			case k, ok1 = <-ch1:
			case <-done:
				return
			}

			select {
			case v, ok2 = <-ch2:
			case <-done:
				return
			}

			if !ok1 || !ok2 {
				close(done)
				return
			}

			if !yield(k, v) {
				close(done)
				return
			}
		}
	}
}

const slicePanicMessage = "itertools: step argument can't be negative"

func Slice[V any](iter iter.Seq[V], start, stop, step int) iter.Seq[V] {
	if step < 0 {
		panic(slicePanicMessage)
	}
	return func(yield func(V) bool) {
		i := 0
		for v := range iter {
			if i >= start && i%step == 0 {
				if !yield(v) {
					return
				}
			}
			i++
			if i >= stop {
				return
			}
		}
	}
}

func Slice2[K, V any](iter iter.Seq2[K, V], start, stop, step int) iter.Seq2[K, V] {
	if step < 0 {
		panic(slicePanicMessage)
	}
	return func(yield func(K, V) bool) {
		i := 0
		for k, v := range iter {
			if i >= start && i%step == 0 {
				if !yield(k, v) {
					return
				}
			}
			i++
			if i >= stop {
				return
			}
		}
	}
}

func Limit[V any](iter iter.Seq[V], limit int) iter.Seq[V] {
	return Slice(iter, 0, limit, 1)
}

func Limit2[K, V any](iter iter.Seq2[K, V], limit int) iter.Seq2[K, V] {
	return Slice2(iter, 0, limit, 1)
}

func TakeWhile[V any](predicate func(V) bool, iter iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range iter {
			if !predicate(v) || !yield(v) {
				return
			}
		}
	}
}

func DropWhile[V any](predicate func(V) bool, iter iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		skip := true
		iter(func(v V) bool {
			if skip {
				if predicate(v) {
					return true
				}
				skip = false
			}
			return yield(v)
		})
	}
}
