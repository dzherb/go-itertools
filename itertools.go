package itertools

import "iter"

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

func Repeat[T any](elem T, n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < n; i++ {
			if !yield(elem) {
				return
			}
		}
	}
}

func Limit[V any](iter iter.Seq[V], limit int) iter.Seq[V] {
	return func(yield func(V) bool) {
		i := 0
		for v := range iter {
			i++
			if !yield(v) || i >= limit {
				return
			}
		}
	}
}
