package graphmatrix

import (
	"fmt"
	"testing"
)

func TestGraphMatrix(t *testing.T) {
	z, err := NewZero(6)
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

	i := []uint32{0, 0, 1, 2, 3, 8}
	j := []uint32{1, 2, 2, 3, 2, 7}
	z, err = NewFromSortedIJ(i, j)
	if err != nil {
		t.Errorf("Error creating NewFromSortedIJ()")
	}
	// fmt.Println("should be t, t, t, f")
	ss := []uint32{0, 0, 1, 3, 2, 2, 8, 7}
	ds := []uint32{1, 2, 2, 2, 3, 1, 7, 8}
	rs := []bool{true, true, true, true, true, false, true, false}

	for i := 0; i < len(ss); i++ {
		if rx := z.GetIndex(ss[i], ds[i]); rx != rs[i] {
			t.Errorf("Error in GetIndex(%d,%d): got %v, want %v", ss[i], ds[i], rx, rs[i])
		}
	}
	if r1, _ := z.GetRow(0); !(len(r1) == 2 && r1[0] == 1 && r1[1] == 2) {
		t.Errorf("Error in GetRow(0): got %v, want %v", r1, []uint32{1, 2})
	} // should be [1, 2]

	if r1, _ := z.GetRow(2); !(len(r1) == 1 && r1[0] == 3) {
		t.Errorf("Error in GetRow(0): got %v, want %v", r1, []uint32{3})
	}
	i = []uint32{1, 2, 3, 0, 0, 2}
	j = []uint32{2, 3, 2, 1, 2, 3}

	_ = SortIJ(&i, &j)
	z, _ = NewFromSortedIJ(i, j)
	if z.N() != 5 {
		t.Errorf("Error in duplicate edge creation: Dim(): got %d, want %d", z.N(), 5)
	}
}

func TestNZIter(t *testing.T) {
	a := []uint32{0, 0, 1, 1, 1, 2, 3, 4, 4, 5, 10}
	b := []uint32{1, 2, 0, 2, 3, 4, 4, 0, 5, 1, 8}

	_ = SortIJ(&a, &b)
	g, _ := NewFromSortedIJ(a, b)

	it := g.NewNZIter()
	itct := 0
	for !it.Done {
		itct++
		r, c := it.Next()
		if !g.GetIndex(r, c) {
			t.Errorf("(%d, %d) not found in graph", r, c)
		}
	}
	if itct != len(a) {
		t.Errorf("Iterator count does not match values: got %d, want %d", itct, len(a))
	}
}
func TestUniqSorted(t *testing.T) {
	a := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 10, 12, 14, 16}
	b := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 10, 12, 14, 16}
	UniqSorted(&a)
	if len(a) != len(b) {
		t.Error("len(a) != len(b)")
	}
	for i := range a {
		if a[i] != b[i] {
			t.Errorf("a[%d] (%d) != b[%d] (%d)", i, a[i], i, b[i])
		}
	}
	a = []uint64{1, 1, 2, 2, 2, 3, 4, 5, 6, 6, 6, 7, 8, 8, 10, 12, 12, 12, 12, 12, 14, 14, 16, 16, 16}
	UniqSorted(&a)
	if len(a) != len(b) {
		t.Error("len(a) != len(b)")
	}
	for i := range a {
		if a[i] != b[i] {
			t.Errorf("a[%d] (%d) != b[%d] (%d)", i, a[i], i, b[i])
		}
	}
}

func benchmarkGraphMatrix(s, d []uint32, b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = NewFromSortedIJ(s, d)
	}
}

func genRandVec(n int) []uint32 {
	v := make([]uint32, n)
	// m := make(map[uint32]bool)
	// for i := 0; i < n; i++ {
	// 	v[i] = rand.Uint32()
	// 	m[v[i]] = true
	// }
	// if n != len(m) {
	// 	fmt.Println("Duplicates found!")
	// }
	return v
}

func benchmarkNZIter(g GraphMatrix, b *testing.B) {
	for n := 0; n < b.N; n++ {
		it := g.NewNZIter()
		for !it.Done {
			_, _ = it.Next()
		}
	}
}

func BenchmarkGraphMatrix(b *testing.B) {
	ns := []int{10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000}
	maxn := int(float64(ns[len(ns)-1]) * 1.2)
	rs := genRandVec(maxn)
	rd := genRandVec(maxn)
	// fmt.Println("pre-sort: len(rs) = ", len(rs), " and len(rd) = ", len(rd))
	if err := SortIJ(&rs, &rd); err != nil {
		b.Errorf("oops: %v", err)
	}

	// fmt.Println("post-sort: len(rs) = ", len(rs), " and len(rd) = ", len(rd))

	for _, n := range ns {
		s := rs[0:n]
		d := rd[0:n]
		name := fmt.Sprintf("n=%d", n)
		b.Run(name, func(b *testing.B) { benchmarkGraphMatrix(s, d, b) })
	}
	fmt.Println("nziter")
	for _, n := range ns {
		s := rs[0:n]
		d := rd[0:n]
		SortIJ(&s, &d)
		g, _ := NewFromSortedIJ(s, d)
		name := fmt.Sprintf("nziter=%d", n)
		b.Run(name, func(b *testing.B) { benchmarkNZIter(g, b) })
	}
}
