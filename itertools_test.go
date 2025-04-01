package itertools

import (
	"iter"
	"reflect"
	"slices"
	"testing"
)

func unwrapIterator2[K, V any](iter iter.Seq2[K, V], limit int) ([]K, []V) {
	i := 0
	keys := make([]K, 0, limit)
	val := make([]V, 0, limit)
	for k, v := range iter {
		if i >= limit {
			break
		}
		keys = append(keys, k)
		val = append(val, v)
		i++
	}
	return keys, val
}

func TestCount(t *testing.T) {
	type args struct {
		start int
		step  int
	}
	tests := []struct {
		name  string
		args  args
		limit int
		want  []int
	}{
		{
			name:  "from 0, step 1",
			args:  args{0, 1},
			limit: 10,
			want:  []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name:  "from -10, step 1",
			args:  args{-10, 1},
			limit: 10,
			want:  []int{-10, -9, -8, -7, -6, -5, -4, -3, -2, -1},
		},
		{
			name:  "from 10, step -1",
			args:  args{10, -1},
			limit: 10,
			want:  []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
		{
			name:  "from 0, step 2",
			args:  args{0, 2},
			limit: 5,
			want:  []int{0, 2, 4, 6, 8},
		},
		{
			name:  "from -100, step 10",
			args:  args{-100, 10},
			limit: 5,
			want:  []int{-100, -90, -80, -70, -60},
		},
		{
			name:  "from 10, step -5",
			args:  args{10, -5},
			limit: 5,
			want:  []int{10, 5, 0, -5, -10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := Count(tt.args.start, tt.args.step)
			if got := slices.Collect(Limit(count, tt.limit)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCycle(t *testing.T) {
	type args[V any] struct {
		iter iter.Seq[V]
	}
	type testCase struct {
		name  string
		args  args[int]
		limit int
		want  []int
	}
	tests := []testCase{
		{
			name: "[1, 2, 3] cycle",
			args: args[int]{
				slices.Values([]int{1, 2, 3}),
			},
			limit: 9,
			want:  []int{1, 2, 3, 1, 2, 3, 1, 2, 3},
		},
		{
			name: "[-10, 0] cycle",
			args: args[int]{
				slices.Values([]int{-10, 0}),
			},
			limit: 5,
			want:  []int{-10, 0, -10, 0, -10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cycle := Cycle(tt.args.iter)
			if got := slices.Collect(Limit(cycle, tt.limit)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cycle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCycle2(t *testing.T) {
	type args[K any, V any] struct {
		iter iter.Seq2[K, V]
	}
	type testCase struct {
		name       string
		args       args[int, string]
		limit      int
		wantKeys   []int
		wantValues []string
	}
	tests := []testCase{
		{
			name: "map cycle",
			args: args[int, string]{
				Zip(slices.Values([]int{2, 4, 6}), slices.Values([]string{"a", "b", "c"})),
			},
			limit:      9,
			wantKeys:   []int{2, 4, 6, 2, 4, 6, 2, 4, 6},
			wantValues: []string{"a", "b", "c", "a", "b", "c", "a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cycle := Cycle2(tt.args.iter)
			gotKeys, gotValues := unwrapIterator2(cycle, tt.limit)

			if got := gotKeys; !reflect.DeepEqual(got, tt.wantKeys) {
				t.Errorf("Cycle2() = %v, want %v", got, tt.wantKeys)
			}

			if got := gotValues; !reflect.DeepEqual(got, tt.wantValues) {
				t.Errorf("Cycle2() = %v, want %v", got, tt.wantValues)
			}
		})
	}
}

func TestRepeat(t *testing.T) {
	type args struct {
		elem int
		n    int
	}
	type testCase struct {
		name string
		args args
		want []int
	}
	tests := []testCase{
		{
			name: "2 for 5 times",
			args: args{2, 5},
			want: []int{2, 2, 2, 2, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rep := Repeat(tt.args.elem, tt.args.n)
			if got := slices.Collect(rep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repeat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZip(t *testing.T) {
	type args struct {
		first  iter.Seq[int]
		second iter.Seq[int]
	}
	type testCase struct {
		name       string
		args       args
		wantKeys   []int
		wantValues []int
	}
	tests := []testCase{
		{
			name: "equal length",
			args: args{
				first:  slices.Values([]int{1, 2, 3, 4, 5}),
				second: slices.Values([]int{6, 7, 8, 9, 10}),
			},
			wantKeys:   []int{1, 2, 3, 4, 5},
			wantValues: []int{6, 7, 8, 9, 10},
		},
		{
			name: "unequal length",
			args: args{
				first:  slices.Values([]int{1, 2, 3}),
				second: slices.Values([]int{6, 7, 8, 9, 10}),
			},
			wantKeys:   []int{1, 2, 3},
			wantValues: []int{6, 7, 8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zip := Zip(tt.args.first, tt.args.second)
			keys, val := unwrapIterator2(zip, 5)
			if !reflect.DeepEqual(keys, tt.wantKeys) {
				t.Errorf("Zip() = %v, want keys %v", keys, tt.wantKeys)
			}

			if !reflect.DeepEqual(val, tt.wantValues) {
				t.Errorf("Zip() = %v, want values %v", val, tt.wantValues)
			}
		})
	}
}

func TestSlice(t *testing.T) {
	type args struct {
		iter  iter.Seq[int]
		start int
		stop  int
		step  int
	}
	type testCase struct {
		name string
		args args
		want []int
	}
	tests := []testCase{
		{
			name: "simple",
			args: args{
				iter:  slices.Values([]int{1, 2, 3, 4, 5}),
				start: 0,
				stop:  5,
				step:  1,
			},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "step 2",
			args: args{
				iter:  slices.Values([]int{1, 2, 3, 4, 5}),
				start: 0,
				stop:  5,
				step:  2,
			},
			want: []int{1, 3, 5},
		},
		{
			name: "start 2, stop 8",
			args: args{
				iter:  slices.Values([]int{1, 2, 3, 4, 5}),
				start: 2,
				stop:  4,
				step:  1,
			},
			want: []int{3, 4},
		},
		{
			name: "start 2, stop 8, step 2",
			args: args{
				iter:  slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
				start: 2,
				stop:  8,
				step:  2,
			},
			want: []int{3, 5, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice := Slice(tt.args.iter, tt.args.start, tt.args.stop, tt.args.step)
			if got := slices.Collect(slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}
