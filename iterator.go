package graphmatrix

// This file contains code for iterating over graphmatrices.

// NZIter is an iterator over the defined points in the graphmatrix.
// If NZIter.Done is true, there are no more points defined.
// Changing the graphmatrix in the middle of an iteration will lead
// to undefined (and almost certainly unwanted) behavior.
type NZIter struct {
	g                GraphMatrix
	Done             bool 	 // `true` if the iterator has exhausted all nonzero values.
	rowIndex         uint32  // index into the row.
	colIndex         uint32  // index into the column value within a given row.
	rowStart, rowEnd uint64  // from IndPtr, the index of the row's start and end.
}

// NewNZIter creates a new graphmatrix iterator over a graphmatrix.
func (g GraphMatrix) NewNZIter() NZIter {
	rowEnd := g.IndPtr[1]
	return NZIter{g: g, Done: false, rowEnd: rowEnd}
}

// advance moves the iterator to the next nonzero entry, returning boolean
// if we're at the end.
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

// Next returns the next nonzero entry in the iterator, returning its row and column.
// The iterator state is modified so that subsequent calls to `Next()` will retrieve
// successive nonzero values. Once all values are produced, `Next()` will set the
// iterator's `Done` field to `true` and will return `0, 0`.
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
