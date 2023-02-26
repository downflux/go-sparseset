package sparseset

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInsert(t *testing.T) {
	type config struct {
		name string
		s    *S[float64]
		x    ID
		data float64
		want *S[float64]
	}

	configs := []config{
		{
			name: "Trivial",
			s:    New[float64](),
			x:    0,
			data: 1,
			want: &S[float64]{
				ids:  []ID{0},
				data: []float64{1},
				sparse: [][]int{
					make([]int, sPage),
				},
			},
		},
		{
			name: "Trivial/OutOfOrder",
			s:    New[float64](),
			x:    1,
			data: 1,
			want: &S[float64]{
				ids:  []ID{1},
				data: []float64{1},
				sparse: [][]int{
					make([]int, sPage),
				},
			},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			c.s.Insert(c.x, c.data)
			if diff := cmp.Diff(
				c.want,
				c.s,
				cmp.AllowUnexported(S[float64]{}),
			); diff != "" {
				t.Errorf("Insert() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
