package graphmatrix

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
