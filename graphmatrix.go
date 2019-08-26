package graphmatrix

// graphmatrix provides csc sparse matrices up to 2^32-1 rows/columns.

import (
	"errors"
	"fmt"
	"sort"
)

// GraphMatrix holds a row index and vector of column pointers.
// If a point is defined at a particular row i and column j, an
// edde exists between vertex i and vertex j.
// GraphMatrices thus represent directed graphs; undirected graphs
// must explicitly set the reverse edge from j to i.
type GraphMatrix struct {
	rowidx []uint32 // contains the row values for each column. A stride represents the outneighbors of a vertex at col j.
	colptr []uint64 // indexes into rowidx - must be twice the width of rowidx.
}

// NewGraphMatrix creates an m x m sparse vector.
func New(m int) (GraphMatrix, error) {
	if m < 0 {
		return GraphMatrix{}, errors.New("dimensions must be non-negative")
	}
	r := make([]uint32, 0)
	c := make([]uint64, m+1)
	return GraphMatrix{rowidx: r, colptr: c}, nil
}

func NewFromRC(r []uint32, c []uint64) (GraphMatrix, error) {
	return GraphMatrix{rowidx: r, colptr: c}, nil
}

func (v GraphMatrix) String() string {
	return fmt.Sprintf("GraphMatrix %v, %v, size %d", v.rowidx, v.colptr, v.Dim())
}

// returns the maximum uint and its position in the vector.
// returns -1 as position if the vector is empty.
func maxUint32(v []uint32) (max uint32, maxPos int) {
	if len(v) == 0 {
		return 0, -1
	}
	max = 0
	maxPos = 0
	for i, n := range v {
		if n > max {
			max = n
			maxPos = i
		}
	}
	return max, maxPos
}

// NewFromSortedIJ creates a graph matrix using i,j as src, dst
// of edges. Assumes i and j are already sorted, j first.
func NewFromSortedIJ(s, d []uint32) (GraphMatrix, error) {
	if len(s) != len(d) {
		return GraphMatrix{}, errors.New("graph inputs must be of the same length")
	}
	m1 := s[len(s)-1]     // max s - this is O(1)
	m2, _ := maxUint32(d) // max d - this is O(n)
	m := m1
	if m2 > m1 {
		m = m2
	}
	m++

	colptr := make([]uint64, s[0]+1, m)
	currval := s[0]
	for i, n := range s {
		if n > currval { // the row has changed
			colptr = append(colptr, uint64(i))
			currval = n
		}
	}
	fmt.Println("m = ", m)
	colptr = append(colptr, uint64(m)+1)
	return GraphMatrix{rowidx: d, colptr: colptr}, nil
}

// SortIJ sorts two vectors s and d by s, then by d. Modifies s and d.
func SortIJ(s, d []uint32) error {
	if len(s) != len(d) {
		return errors.New("inputs must be of the same length")
	}
	sd := make([]uint64, len(s))
	for i := 0; i < len(s); i++ {
		sd[i] = uint64(s[i])<<32 + uint64(d[i])
	}
	sort.Slice(sd, func(i, j int) bool { return sd[i] < sd[j] })

	for i := 0; i < len(sd); i++ {
		s[i] = uint32(sd[i] >> 32)
		d[i] = uint32(sd[i] & 0x00000000ffffffff)
	}

	return nil
}

// inRange returns true if (r, c) is a valid index into v.
func (v *GraphMatrix) inRange(r, c uint32) bool {
	n := v.Dim()
	return (c < n) && (r < n)
}

// Dim returns the (single-axis) dimension of the GraphMatrix
func (v *GraphMatrix) Dim() uint32 {
	return uint32(len(v.colptr) - 1)
}

// N returns the number of defined values in the GraphMatrix
func (v *GraphMatrix) N() uint64 {
	return uint64(len(v.rowidx))
}

// searchsorted32 finds a value x in sorted vector v.
// Returns index and true/false indicating found.
// lo and hi constrains search to these indices.
// If lo/hi are out of bounds, return -1 and false unless
// the vector is empty, in which case return 0 and false.
func searchsorted32(v []uint32, x uint32, lo, hi uint64) (int, bool) {
	ulen := uint64(len(v))
	if ulen == 0 {
		return 0, false
	}
	if lo == hi {
		return int(lo), false
	}
	if ulen < lo || ulen < hi || lo > hi {
		return -1, false
	}
	s := sort.Search(int(hi-lo), func(i int) bool { return v[int(lo)+i] >= x }) + int(lo)
	found := (s < len(v)) && (v[s] == x)

	return s, found
}

// GetIndex returns true if the value at (r, c) is defined.
func (v *GraphMatrix) GetIndex(r, c uint32) bool {
	if len(v.colptr) <= int(c)+1 {
		return false
	}

	r1 := v.colptr[r]
	r2 := v.colptr[r+1]
	if r1 >= r2 {
		return false
	}
	_, found := searchsorted32(v.rowidx, c, r1, r2)
	return found
}

// GetRow returns the 'n'th row slice, or an empty slice if empty.
func (v *GraphMatrix) GetRow(n uint32) []uint32 {
	p1 := v.colptr[n]
	p2 := v.colptr[n+1]
	if int(p1) > len(v.rowidx) || int(p2) > len(v.rowidx) {
		return []uint32{}
	}
	return v.rowidx[p1:p2]
}

// SetIndex sets the value at (r, c) to true.
// This can be a relatively expensive operation as it can force
// reallocation as the vectors increase in size.
func (v *GraphMatrix) SetIndex(r, c uint32) error {
	if !v.inRange(r, c) {
		return errors.New("index out of range")
	}
	rowStartIdx := v.colptr[r] // this is the pointer into the rowidx for column c
	rowEndIdx := v.colptr[r+1]

	i, found := searchsorted32(v.rowidx, c, rowStartIdx, rowEndIdx)
	if found { // already set
		return nil
	}
	v.rowidx = append(v.rowidx, 0)
	copy(v.rowidx[i+1:], v.rowidx[i:])
	v.rowidx[i] = c

	for i := int(r + 1); i < len(v.colptr); i++ {
		v.colptr[i]++
	}
	return nil
}
