package graphmatrix

import (
	"testing"
)

func TestGraphMatrix(t *testing.T) {
	z, err := New(6)
	if err != nil {
		t.Errorf("Error creating New(6)")
	}
	if x := z.Dim(); x != 6 {
		t.Errorf("Dim(): got %d, want %d", x, 6)
	}
	if err := z.SetIndex(3, 3); err != nil {
		t.Errorf("Error in SetIndex(3,3)")
	}
	if z.SetIndex(1, 2) != nil {
		t.Errorf("Error in SetIndex(1,2)")
	}
	if z.SetIndex(3, 1) != nil {
		t.Errorf("Error in SetIndex(3,1)")
	}
	if z.SetIndex(1, 6) == nil {
		t.Errorf("Error not thrown in SetIndex(1,6): should be out of bounds")
	}

	if y := z.GetIndex(1, 2); !y {
		t.Errorf("GetIndex(1,2): got %v, want %v", y, true)
	}
	if y := z.GetIndex(2, 2); y {
		t.Errorf("GetIndex(2,2): got %v, want %v", y, false)
	}

	i := []uint32{0, 0, 1, 2, 3}
	j := []uint32{1, 2, 2, 3, 2}
	z, err = NewFromSortedIJ(i, j)
	if err != nil {
		t.Errorf("Error creating NewFromSortedIJ()")
	}
	// fmt.Println("should be t, t, t, f")
	ss := []uint32{0, 0, 1, 3, 2, 2}
	ds := []uint32{1, 2, 2, 2, 3, 1}
	rs := []bool{true, true, true, true, true, false}

	for i := 0; i < len(ss); i++ {
		if rx := z.GetIndex(ss[i], ds[i]); rx != rs[i] {
			t.Errorf("Error in GetIndex(%d,%d): got %v, want %v", ss[i], ds[i], rx, rs[i])
		}
	}
	if r1 := z.GetRow(0); !(len(r1) == 2 && r1[0] == 1 && r1[1] == 2) {
		t.Errorf("Error in GetRow(0): got %v, want %v", r1, []uint32{1, 2})
	} // should be [1, 2]

	if r1 := z.GetRow(2); !(len(r1) == 1 && r1[0] == 3) {
		t.Errorf("Error in GetRow(0): got %v, want %v", r1, []uint32{3})
	}
}
