package sweep

import (
	"strconv"
)

// Int2 represents a 2 byte integer.
type Int2 [2]byte

// Int4 represents a 4 byte integer.
type Int4 [4]byte

// Int6 represents a 6 byte integer.
type Int6 [6]byte

// NewInt2 returns a new 2 byte integer with the given integer.
func NewInt2(n int) Int2 {
	if n < 0 || n > 99 {
		panic("sweep: NewInt2: 0 <= n <= 99 must be true")
	}

	const size = 2
	res := strconv.Itoa(n)
	ret := Int2([size]byte{'0', '0'})
	c := 0
	for i := size - len(res); i < size; i++ {
		ret[i] = res[c]
		c++
	}

	return ret
}

// NewInt4 returns a new 4 byte integer with the given integer.
func NewInt4(n int) Int4 {
	if n < 0 || n > 9999 {
		panic("sweep: NewInt4: 0 <= n <= 9999 must be true")
	}

	const size = 4
	res := strconv.Itoa(n)
	ret := Int4([size]byte{'0', '0'})
	c := 0
	for i := size - len(res); i < size; i++ {
		ret[i] = res[c]
		c++
	}

	return ret
}

// NewInt6 returns a new 6 byte integer with the given integer.
func NewInt6(n int) Int6 {
	if n < 0 || n > 999999 {
		panic("sweep: NewInt4: 0 <= n <= 999999 must be true")
	}

	const size = 6
	res := strconv.Itoa(n)
	ret := Int6([size]byte{'0', '0'})
	c := 0
	for i := size - len(res); i < size; i++ {
		ret[i] = res[c]
		c++
	}

	return ret
}

// String returns the 2 byte integer in string representation.
func (n Int2) String() string {
	return string(n[:])
}

// String returns the 4 byte integer in string representation.
func (n Int4) String() string {
	return string(n[:])
}

// String returns the 6 byte integer in string representation.
func (n Int6) String() string {
	return string(n[:])
}

// Int returns the 2 byte integer in integer representation.
func (n Int2) Int() int {
	num, _ := strconv.Atoi(n.String())
	return num
}

// Int returns the 4 byte integer in integer representation.
func (n Int4) Int() int {
	num, _ := strconv.Atoi(n.String())
	return num
}

// Int returns the 6 byte integer in integer representation.
func (n Int6) Int() int {
	num, _ := strconv.Atoi(n.String())
	return num
}
