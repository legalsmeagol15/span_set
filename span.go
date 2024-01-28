package spanset

import "fmt"

type Span[T any] interface {
	GetStart() T
	GetEnd() T
	IncludeStart() bool
	IncludeEnd() bool
	IncludeBefore() bool
	IncludeAfter() bool

	IsSingleton() bool
	IsUniversal() bool
}

func NewSpan[T any](start T, end T) Span {
	if a, ok := any(start).(int); ok {
		if b, ok := any(end).(int); ok {
			return &span{Start: a, End: b, IncludeStart: true, IncludeEnd: true}
		}

	}
	panic(fmt.Sprintf("invalid span types: %v, %v", start, end))
}

type span[T any] struct {
	start, end T

	includeStart, includeEnd, includeBefore, includeAfter, isSingleton, isUniversal bool
}

func GetStart[T any](s *span[T]) T        { return s.start }
func GetEnd[T any](s *span[T]) T          { return s.end }
func IncludeStart[T any](s *span[T]) bool { return s.includeStart }
func IncludeEnd[T any](s *span[T]) bool   { return s.includeEnd }
func IsSingleton[T any](s *span[T]) bool  { return s.start == s.end }
