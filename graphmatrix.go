package graphmatrix

import (
	// "errors"
	"errors"
	"fmt"
	"sort"
	// "sort"
)

// GraphMatrix holds a row index and vector of column pointers.
// If a point is defined at a particular row i and column j, an
// edge exists between vertex i and vertex j.
// GraphMatrices thus represent directed graphs; undirected graphs
// must explicitly set the reverse edge from j to i.
type GraphMatrix struct {
	IndPtr  []uint64 // indexes into Indices - must be twice the width of Indices.
	Indices []uint32 // contains the row values for each column. A stride represents the outneighbors of a vertex at col j.
}

// String is used for pretty printing.
func (g GraphMatrix) String() string {
	return fmt.Sprintf("GraphMatrix %v, %v, size %d", g.IndPtr, g.Indices, g.Dim())
}

// GetRow returns the 'n'th row slice, or an empty slice if empty.
func (g GraphMatrix) GetRow(r uint32) ([]uint32, error) {
	if r > uint32(len(g.IndPtr))-1 {
		return []uint32{}, fmt.Errorf("Row %d out of bounds (max %d)", r, len(g.IndPtr))
	}
	rowStart := g.IndPtr[r]
	rowEnd := g.IndPtr[r+1]
	return g.Indices[rowStart:rowEnd], nil
}

// NZIter is an iterator over the defined points in the graphmatrix.
// If NZIter.Done is true, there are no more points defined.
// Changing the graphmatrix in the middle of an iteration will lead
// to undefined (and almost certainly unwanted) behavior.
type NZIter struct {
	g                GraphMatrix
	Done             bool
	rowIndex         uint32
	colIndex         uint32
	rowStart, rowEnd uint64
}

func (g GraphMatrix) NewNZIter() NZIter {
	rowEnd := g.IndPtr[1]

	return NZIter{g: g, Done: false, rowEnd: rowEnd}
}

func (it *NZIter) advance() bool {
	rowLen := uint32(it.rowEnd - it.rowStart)
	it.colIndex++
	if it.colIndex >= rowLen {
		it.colIndex = 0
		it.rowIndex++
		if it.rowIndex < uint32(len(it.g.IndPtr)-1) {
			it.rowStart = it.g.IndPtr[it.rowIndex]
			it.rowEnd = it.g.IndPtr[it.rowIndex+1]
			return false
		}
		return true
	}
	return false
}

func (it *NZIter) Next() (uint32, uint32) {
	if it.Done {
		return 0, 0
	}
	rowLen := uint32(it.rowEnd - it.rowStart)
	for rowLen == 0 {
		it.Done = it.advance()
		rowLen = uint32(it.rowEnd - it.rowStart)
	}
	r := it.rowIndex
	index := it.rowStart + uint64(it.colIndex)
	c := it.g.Indices[index]

	it.Done = it.advance()
	return r, c
}
func cumsum(p []uint64, c []uint64, n uint64) uint64 {
	nz := uint64(0)
	for i := nz; i < n; i++ {
		p[i] = nz
		nz += uint64(c[i])
		c[i] = p[i]
	}
	p[n] = nz
	return nz
}

func compress(row []uint32, col []uint32, n uint64) (ia []uint64, ja []uint32) {
	w := make([]uint64, n+1)
	ia = make([]uint64, n+1)
	ja = make([]uint32, len(col))

	for _, v := range row {
		w[v]++
	}
	cumsum(ia, w, n)

	for j, v := range col {
		p := w[row[j]]
		ja[p] = v
		w[row[j]]++
	}
	return
}

// inRange returns true if (r, c) is a valid index into v.
func (g *GraphMatrix) inRange(r, c uint32) bool {
	n := g.Dim()
	return (c < n) && (r < n)
}

// Dim returns the (single-axis) dimension of the GraphMatrix
func (g GraphMatrix) Dim() uint32 {
	return uint32(len(g.IndPtr) - 1)
}

// N returns the number of defined values in the GraphMatrix
func (g *GraphMatrix) N() uint64 {
	return uint64(len(g.Indices))
}

// searchSorted32 finds a value x in sorted vector v.
// Returns index and true/false indicating found.
// lo and hi constrains search to these Indices.
// If lo/hi are out of bounds, return -1 and false unless
// the vector is empty, in which case return 0 and false.
func searchSorted32(v []uint32, x uint32, lo, hi uint64) (int, bool) {
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
func (g GraphMatrix) GetIndex(r, c uint32) bool {
	if uint32(len(g.IndPtr)) <= c+1 {
		return false
	}

	r1 := g.IndPtr[r]
	r2 := g.IndPtr[r+1]
	if r1 >= r2 {
		return false
	}
	_, found := searchSorted32(g.Indices, c, r1, r2)
	return found
}

// SetIndex sets the value at (r, c) to true.
// This can be a relatively expensive operation as it can force
// reallocation as the vectors increase in size.
func (g *GraphMatrix) SetIndex(r, c uint32) error {
	if !g.inRange(r, c) {
		return errors.New("index out of range")
	}
	rowStartIdx := g.IndPtr[r] // this is the pointer into the Indices for column c
	rowEndIdx := g.IndPtr[r+1]

	i, found := searchSorted32(g.Indices, c, rowStartIdx, rowEndIdx)
	if found { // already set
		return nil
	}
	g.Indices = append(g.Indices, 0)
	copy(g.Indices[i+1:], g.Indices[i:])
	g.Indices[i] = c

	for i := int(r + 1); i < len(g.IndPtr); i++ {
		g.IndPtr[i]++
	}
	return nil
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

// NewZero creates an m x m sparse matrix.
func NewZero(m int) (GraphMatrix, error) {
	if m < 0 {
		return GraphMatrix{}, errors.New("dimensions must be non-negative")
	}
	i := make([]uint32, 0)
	ip := make([]uint64, m+1)
	return GraphMatrix{IndPtr: ip, Indices: i}, nil
}

func NewFromSortedIJ(s []uint32, d []uint32) (GraphMatrix, error) {
	if len(s) != len(d) {
		return GraphMatrix{}, fmt.Errorf("graph inputs must be of the same length (got %d, %d)", len(s), len(d))
	}
	m1 := s[len(s)-1]     // max s - this is O(1)
	m2, _ := maxUint32(d) // max d - this is O(n)
	m := m1
	if m2 > m1 {
		m = m2
	}
	m++ // m is the number of rows/cols for this matrix.

	ia, ja := compress(s, d, uint64(m))
	// ja, data = dedupe(ia, ja, data, c.r, c.c)
	return GraphMatrix{ia, ja}, nil
}

// SortIJ sorts two vectors s and d by s, then by d, and eliminates any duplicate pairs.
// Modifies s and d.
func SortIJ(s, d *[]uint32) error {
	if len(*s) != len(*d) {
		return errors.New("inputs must be of the same length")
	}
	sd := make([]uint64, len(*s))
	for i := 0; i < len(*s); i++ {
		sd[i] = uint64((*s)[i])<<32 + uint64((*d)[i])
	}
	sort.Slice(sd, func(i, j int) bool { return sd[i] < sd[j] })
	UniqSorted(&sd)

	for i := 0; i < len(sd); i++ {
		(*s)[i] = uint32(sd[i] >> 32)
		(*d)[i] = uint32(sd[i] & 0x00000000ffffffff)
	}

	(*s) = (*s)[:len(sd)]
	(*d) = (*d)[:len(sd)]
	return nil
}

func UniqSorted(a *[]uint64) {
	j := 0
	for i := 1; i < len(*a); i++ {
		if (*a)[j] == (*a)[i] {
			continue
		}
		j++
		// preserve the original data
		// in[i], in[j] = in[j], in[i]
		// only set what is required
		(*a)[j] = (*a)[i]
	}
	(*a) = (*a)[:j+1]
}
