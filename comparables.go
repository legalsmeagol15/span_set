package spanset

import "golang.org/x/exp/constraints"

type integerType interface {
	int | int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | uint
}
type compable[T any] struct {
	Compare  func(a, b compable[T]) int
	absolute T
}

func intCompare(a, b compable[int]) int { return a.absolute - b.absolute }

func newInt[T integerType](item T) compable[int] {
	return compable[int]{absolute: int(item), Compare: intCompare}
}

func smooo[T constraints.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}
