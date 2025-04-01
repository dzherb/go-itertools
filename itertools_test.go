package itertools

import (
	"iter"
	"reflect"
	"testing"
)

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
			if got := unwrapIterator(count, tt.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func unwrapIterator[T any](iter iter.Seq[T], limit int) []T {
	i := 0
	res := make([]T, 0, limit)
	for item := range iter {
		if i >= limit {
			break
		}
		res = append(res, item)
		i++
	}
	return res
}
