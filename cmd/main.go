package main

import (
	"fmt"
	"log"

	"github.com/sbromberger/graphmatrix"
)

func main() {
	z, err := graphmatrix.New(6)
	if err != nil {
		log.Fatal("error creating graphmatrix", err)
	}
	fmt.Println(z.Size())
	if err := z.SetIndex(3, 3); err != nil {
		log.Fatal("3,3 error: ", err)
	}
	if z.SetIndex(1, 2) != nil {
		log.Fatal("1,2 error")
	}

	fmt.Println(z)
	y := z.GetIndex(1, 2)
	if !y {
		log.Fatal("1,2 not set (incorrect)")
	} else {
		fmt.Println("1,2 set (correct)")
	}
	if !z.GetIndex(2, 2) {
		fmt.Println("2,2 not set (correct)")
	} else {
		log.Fatal("2,2 set (incorrect)")
	}
	fmt.Println("----- graphmatrix.New(3)")
	z, _ = graphmatrix.New(4)
	z.SetIndex(0, 1)
	z.SetIndex(0, 2)
	z.SetIndex(1, 2)
	z.SetIndex(2, 3)
	fmt.Println(z)

	fmt.Println("string is")
	fmt.Println(z)
	fmt.Println("----- graphmatrix.NewFromSorted()")
	i := []uint32{0, 0, 1, 2}
	j := []uint32{1, 2, 2, 3}
	z, err = graphmatrix.NewFromSortedIJ(i, j)
	if err != nil {
		log.Fatal("NewFromSortedIJ error: ", err)
	}
	fmt.Println(z)
	fmt.Println("should be t, t, t, f")
	fmt.Println(z.GetIndex(0, 1))
	fmt.Println(z.GetIndex(0, 2))
	fmt.Println(z.GetIndex(1, 2))
	fmt.Println(z.GetIndex(2, 1))
}
