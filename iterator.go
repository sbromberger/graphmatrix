package graphmatrix

// This file contains code for iterating over graphmatrices.

// NZIter is an iterator over the defined points in the graphmatrix.
// If NZIter.Done is true, there are no more points defined.
// Changing the graphmatrix in the middle of an iteration will lead
// to undefined (and almost certainly unwanted) behavior.
type NZIter struct {
	g                GraphMatrix
	rowIndex         uint32 // index into the row.
	colIndex         uint32 // index into the column value within a given row.
	rowStart, rowEnd uint64 // from IndPtr, the index of the row's start and end.
	n                uint64 // number of elements already seen.
}

// NewNZIter creates a new graphmatrix iterator over a graphmatrix.
func (g GraphMatrix) NewNZIter() NZIter {
	rowEnd := g.IndPtr[1]
	return NZIter{g: g, rowEnd: rowEnd}
}

// advance moves the iterator to the next nonzero entry, returning boolean
// if we're at the end.
func (it *NZIter) advance() bool {
	rowLen := uint32(it.rowEnd - it.rowStart)
	it.colIndex++
	if it.colIndex >= rowLen { // if we're at the end of the current row
		it.colIndex = 0
		it.rowIndex++                                 // move to the next row
		if it.rowIndex < uint32(len(it.g.IndPtr)-1) { // if we're not on the last row
			it.rowStart = it.g.IndPtr[it.rowIndex]
			it.rowEnd = it.g.IndPtr[it.rowIndex+1]
			return false
		}
		return true // we're on the last row, and we're at the end of the current row.
	}
	return false // we're not at the end of the current row.
}

// Done returns true if the iterator has exhausted all defined points.
func (it *NZIter) Done() bool {
	return it.n >= it.g.N()
}

// Next returns the next nonzero entry in the iterator, returning its row and column.
// The iterator state is modified so that subsequent calls to `Next()` will retrieve
// successive nonzero values. Once all values are produced, `Next()` will set the
// iterator's `Done` field to `true` and will return `0, 0`.
func (it *NZIter) Next() (uint32, uint32, bool) {
	if it.Done() {
		return 0, 0, true
	}
	rowLen := uint32(it.rowEnd - it.rowStart)
	for rowLen == 0 {
		if it.advance() { // we've hit the end.
			return 0, 0, true
		}
		rowLen = uint32(it.rowEnd - it.rowStart)
	}

	r := it.rowIndex
	index := it.rowStart + uint64(it.colIndex)
	c := it.g.Indices[index]

	it.n++
	return r, c, it.advance()
}
