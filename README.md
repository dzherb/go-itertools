# Miscellaneous utils for iteration

`go-itertools` is a versatile library inspired by Python's [itertools](https://docs.python.org/3/library/itertools.html), designed to provide elegant and flexible utilities for working with sequences in Go. It offers a wide variety of functions to help you manipulate, combine, and transform iterators in a simple and efficient way.

## Why this is a useful tool

Go's standard library provides some built-in iteration capabilities, but working with sequences can become cumbersome for more advanced operations. `go-itertools` bridges this gap by offering powerful tools that allow you to write cleaner and more expressive code, especially when dealing with complex sequences.

This library is beneficial in scenarios where you need to:

- Combine multiple sequences in a memory-efficient manner (e.g., chaining or zipping sequences).
- Create infinite or cyclic sequences.
- Apply transformations or filters on data with minimal overhead.
- Write clean and readable code that avoids excessive boilerplate.

### Go and Custom Iterators

Starting from Go 1.23, the language introduced the `iter.Seq` interface, which makes it easy to create [custom iterators](https://go.dev/blog/range-functions) that can be used with Go's `range` loop. This powerful feature allows developers to define their own iteration logic while still leveraging Go's native iteration capabilities. By using type parameters (generics) and the `iter.Seq` interface, we can write reusable and type-safe code that works with sequences of any type. This approach is the foundation of `go-itertools`, enabling more elegant and flexible ways to work with iterators in Go.

---

## Examples

Here are a couple of examples showing how you can use `go-itertools` to work with sequences.

### Basic Iteration with Itertools

```go
package main

import (
	"fmt"
	it "github.com/dzherb/go-itertools"
)

func main() {
	seq1 := it.Repeat("hi", 3)
	seq2 := it.Count(1, 1)

	for k, v := range it.Zip(seq1, seq2) {
		fmt.Println(k, v) // hi 1, hi 2, hi 3
	}

	seq3 := it.FromElements(10, 20, 30)
	seq4 := it.FromElements(40, 50)

	for s := range it.Chain(seq3, seq4) {
		fmt.Println(s) // 10, 20, 30, 40, 50
	}
}
```

In this example:
- `Repeat` creates a sequence repeating the value `"hi"` three times.
- `Count` generates a sequence of numbers starting at 1 and incrementing by 1.
- `Zip` combines two sequences element by element.
- `Chain` concatenates multiple sequences into a single one.

### Stream API

```go
package main

import (
	"fmt"
	it "github.com/dzherb/go-itertools"
	"github.com/dzherb/go-itertools/stream"
)

func main() {
	seq := it.Cycle(it.FromElements(1, 2, 3))

	st := stream.FromIterator(seq).
		Filter(func(i int) bool { return i != 2 }).
		Map(func(i int) int { return i * 2 }).
		Take(6)

	for s := range st {
		fmt.Println(s) // 2, 6, 2, 6, 2, 6
	}
}
```

In this example, the `stream` API is used to create a pipeline that:
- Starts with a cyclic sequence `[1, 2, 3]`.
- Filters out the value `2`.
- Maps the remaining values to their doubles.
- Limits the sequence to the first 6 elements.

This shows the power of the stream API, which lets you easily chain multiple operations together to create complex transformations in a clear and readable manner.

--- 

With `go-itertools`, you can easily build complex pipelines for processing sequences, making your Go code more refined, modular, and easier to understand. Whether you are working with simple sequences or need to create intricate transformations, this library provides all the tools you need.