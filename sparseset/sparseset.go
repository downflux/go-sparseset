package sparseset

type ID uint32

const (
	// mPage is the page mask. The page ID is the first 19 bits of the ID. A
	// page therefore comprises of 2 ** 13 entities.
	mPage = 0xffffe000

	// mID is the ID  mask. This is the last 13 bits of the ID.
	mID = 0x00001fff

	sPage = 1 << 13
)

func (x ID) Page() uint32 { return uint32((x & mPage) >> 19) }
func (x ID) ID() uint32   { return uint32(x & mID) }

type S[T any] struct {
	// ids is a dense array of data IDs. This is used for the normal sparse
	// set lookups.
	ids  []ID
	data []T

	// sparse is a lookup array. The indices of this array are ID values;
	// the values of this array are the indices of the packed ID array. Note
	// that in our implmentation, this sparse array is actually paged to
	// help cut down on needless allocs.
	sparse [][]int
}

func New[T any]() *S[T] {
	return &S[T]{
		ids:    nil,
		data:   nil,
		sparse: nil,
	}
}

func (s *S[T]) Insert(x ID, data T) {
	page, id := x.Page(), x.ID()

	if len(s.sparse) < int(page)+1 {
		dest := make([][]int, page+1)
		copy(dest, s.sparse)
		// pre-allocate the 32kB page.
		dest[page] = make([]int, sPage)
		s.sparse = dest
	}

	s.ids = append(s.ids, x)
	s.data = append(s.data, data)
	s.sparse[page][int(id)] = len(s.ids) - 1
}

func (s *S[T]) Remove(x ID) {
	page, id := x.Page(), x.ID()

	if len(s.sparse) <= int(page) {
		return
	}

	expected := s.sparse[page][id]
	var none T
	s.data[expected] = none
	s.sparse[page][id] = 0
}

func (s *S[T]) In(x ID) bool {
	page, id := x.Page(), x.ID()
	if len(s.sparse) <= int(page) {
		return false
	}

	expected := s.sparse[page][id]
	if len(s.data) < expected {
		return false
	}
	return s.ids[expected] == x
}
