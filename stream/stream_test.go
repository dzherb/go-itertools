package stream

import (
	"reflect"
	"testing"
)

func TestFromElementsAndCollect(t *testing.T) {
	s := FromElements(1, 2, 3, 4)
	got := s.Collect()
	want := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FromElements().Collect() = %v; want %v", got, want)
	}
}

func TestMap(t *testing.T) {
	s := FromElements(1, 2, 3)
	mapped := s.Map(func(v int) int { return v * v })
	got := mapped.Collect()
	want := []int{1, 4, 9}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Map().Collect() = %v; want %v", got, want)
	}
}

func TestStandaloneMap(t *testing.T) {
	s := FromElements("a", "bb", "ccc")
	mapped := Map(s, func(s string) int { return len(s) })
	got := mapped.Collect()
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Standalone Map().Collect() = %v; want %v", got, want)
	}
}

func TestFilter(t *testing.T) {
	s := FromElements(1, 2, 3, 4, 5)
	filtered := s.Filter(func(v int) bool { return v%2 == 0 })
	got := filtered.Collect()
	want := []int{2, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Filter().Collect() = %v; want %v", got, want)
	}
}

func TestTake(t *testing.T) {
	s := FromElements(10, 20, 30, 40)
	taken := s.Take(2)
	got := taken.Collect()
	want := []int{10, 20}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Take().Collect() = %v; want %v", got, want)
	}
}

func TestSlice(t *testing.T) {
	s := FromElements(0, 1, 2, 3, 4, 5, 6)
	sliced := s.Slice(1, 6, 2)
	got := sliced.Collect()
	want := []int{1, 3, 5}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Slice().Collect() = %v; want %v", got, want)
	}
}

func TestTakeWhile(t *testing.T) {
	s := FromElements(1, 2, 3, 4, 1, 0)
	tw := s.TakeWhile(func(v int) bool { return v < 4 })
	got := tw.Collect()
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("TakeWhile().Collect() = %v; want %v", got, want)
	}
}

func TestDropWhile(t *testing.T) {
	s := FromElements(0, 0, 1, 2, 0)
	dw := s.DropWhile(func(v int) bool { return v == 0 })
	got := dw.Collect()
	want := []int{1, 2, 0}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DropWhile().Collect() = %v; want %v", got, want)
	}
}

func TestForEach(t *testing.T) {
	s := FromElements(1, 2, 3)
	var result []int
	s.ForEach(func(v int) { result = append(result, v*10) })
	want := []int{10, 20, 30}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("ForEach() result = %v; want %v", result, want)
	}
}

func TestFromChan(t *testing.T) {
	ch := make(chan int, 3)
	ch <- 100
	ch <- 200
	ch <- 300
	close(ch)

	s := FromChan(ch)
	got := s.Collect()
	want := []int{100, 200, 300}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FromChan().Collect() = %v; want %v", got, want)
	}
}
