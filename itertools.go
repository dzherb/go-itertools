package itertools

import (
	"iter"
)

// FromElements creates an iterator from a fixed list of elements.
//
// Example:
//
//	elems := FromElements("a", "b", "c")
//	for v := range elems {
//	    fmt.Println(v) // "a", "b", "c"
//	}
func FromElements[V any](elems ...V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range elems {
			if !yield(v) {
				return
			}
		}
	}
}

// FromPairs creates an iterator from a slice of key-value pairs.
//
// Example:
//
//	pairs := FromPairs([][2]string{{"a", "1"}, {"b", "2"}})
//	for k, v := range pairs {
//	    fmt.Println(k, v) // "a" "1", "b" "2"
//	}
func FromPairs[V any](pairs [][2]V) iter.Seq2[V, V] {
	return func(yield func(V, V) bool) {
		for _, pair := range pairs {
			if !yield(pair[0], pair[1]) {
				return
			}
		}
	}
}

// Count returns an infinite sequence of numbers starting from `start`, incremented by `step`.
//
// Example:
//
//	count := Count(10, 2)
//	for n := range Take(count, 5) {
//		fmt.Println(n) // 10, 12, 14, 16, 18
//	}
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

// Cycle returns an infinite sequence cycling through elements of `iter`.
//
// Example:
//
//	cycle := Cycle(FromElements("a", "b", "c"))
//	for v := range Take(cycle, 5) {
//		fmt.Println(v) // "a", "b", "c", "a", "b"
//	}
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

// Cycle2 returns an infinite sequence cycling through key-value pairs of `iter`.
//
// Example:
//
//	cycle := Cycle2(FromPairs([][2]string{{"a", "1"}, {"b", "2"}}))
//	for k, v := range Take2(cycle, 4) {
//		fmt.Println(k, v) // "a" "1", "b" "2", "a" "1", "b" "2"
//	}
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

// Repeat returns a sequence repeating `elem` exactly `n` times.
//
// Example:
//
//	repeat := Repeat("hello", 3)
//	for v := range repeat {
//		fmt.Println(v) // "hello", "hello", "hello"
//	}
func Repeat[V any](elem V, n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		for i := 0; i < n; i++ {
			if !yield(elem) {
				return
			}
		}
	}
}

// Chain returns a sequence that combines multiple sequences into one.
//
// Example:
//
//	chain := Chain(FromElements(1, 2), FromElements(3, 4))
//	for v := range chain {
//		fmt.Println(v) // 1, 2, 3, 4
//	}
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

// Chain2 returns a sequence that combines multiple key-value pair sequences into one.
//
// Example:
//
//	chain := Chain2(
//		FromPairs([][2]string{{"a", "1"}, {"b", "2"}}),
//		FromPairs([][2]string{{"c", "3"}}),
//	)
//	for k, v := range chain {
//		fmt.Println(k, v) // "a" "1", "b" "2", "c" "3"
//	}
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

// Zip combines two iterables into a sequence of key-value pairs.
// The function yields pairs from the two sequences in parallel,
// until one or both sequences are exhausted.
//
// Example:
//
//	first := FromElements(1, 2, 3)
//	second := FromElements("a", "b", "c")
//	zip := Zip(first, second)
//	for k, v := range zip {
//		fmt.Println(k, v) // 1 "a", 2 "b", 3 "c"
//	}
//
// Returns a sequence of key-value pairs formed by combining elements from the two sequences.
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

// Slice returns a sequence that includes elements from the input sequence `iter`
// starting from index `start` to `stop`, with a specified `step`.
//
// It panics if `step` is negative.
//
// Example:
//
//	iter := FromElements(1, 2, 3, 4, 5, 6)
//	slice := Slice(iter, 1, 5, 2)
//	for v := range slice {
//		fmt.Println(v) // 2, 4
//	}
func Slice[V any](iter iter.Seq[V], start, stop, step int) iter.Seq[V] {
	if step < 0 {
		panic(slicePanicMessage)
	}

	return func(yield func(V) bool) {
		i := 0
		for v := range iter {
			if i >= start && i < stop && (i-start)%step == 0 {
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

// Slice2 returns a sequence of key-value pairs from the input sequence `iter`
// starting from index `start` to `stop`, with a specified `step`.
//
// It panics if `step` is negative.
//
// Example:
//
//	iter := FromPairs([][2]string{{"a", "1"}, {"b", "2"}, {"c", "3"}})
//	slice := Slice2(iter, 0, 2, 1)
//	for k, v := range slice {
//		fmt.Println(k, v) // "a" "1", "b" "2"
//	}
func Slice2[K, V any](iter iter.Seq2[K, V], start, stop, step int) iter.Seq2[K, V] {
	if step < 0 {
		panic(slicePanicMessage)
	}
	return func(yield func(K, V) bool) {
		i := 0
		for k, v := range iter {
			if i >= start && i < stop && (i-start)%step == 0 {
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

// Take returns a sequence containing the first n elements from the input sequence `iter`.
// This is equivalent to applying Slice with start = 0, stop = n, and step = 1.
//
// Example:
//
//	iter := FromElements(1, 2, 3, 4, 5)
//	taken := Take(iter, 3)
//	for v := range taken {
//		fmt.Println(v) // Output: 1, 2, 3
//	}
func Take[V any](iter iter.Seq[V], n int) iter.Seq[V] {
	return Slice(iter, 0, n, 1)
}

// Take2 returns a sequence containing the first n key-value pairs from the input sequence `iter`.
// This is equivalent to applying Slice2 with start = 0, stop = n, and step = 1.
//
// Example:
//
//	iter := FromPairs([][2]string{{"a", "1"}, {"b", "2"}, {"c", "3"}})
//	taken := Take2(iter, 2)
//	for k, v := range taken {
//		fmt.Println(k, v) // Output: "a" "1", "b" "2"
//	}
func Take2[K, V any](iter iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	return Slice2(iter, 0, n, 1)
}

// TakeWhile returns a sequence of elements from `iter` as long
// as the predicate function returns true for each element.
// The sequence stops as soon as the predicate returns false.
//
// Example:
//
//	iter := FromElements(1, 2, 3, 4, 5)
//	taken := TakeWhile(func(v int) bool { return v < 4 }, iter)
//	for v := range taken {
//		fmt.Println(v) // Output: 1, 2, 3
//	}
func TakeWhile[V any](predicate func(V) bool, iter iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range iter {
			if !predicate(v) || !yield(v) {
				return
			}
		}
	}
}

// DropWhile returns a sequence of elements from `iter`
// after the predicate function first returns false.
// Initially, elements are skipped until the predicate no longer holds,
// and then the rest of the sequence is returned.
//
// Example:
//
//	iter := iter.FromElements(1, 2, 3, 4, 5)
//	dropped := DropWhile(func(v int) bool { return v < 3 }, iter)
//	for v := range dropped {
//		fmt.Println(v) // Output: 3, 4, 5
//	}
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
