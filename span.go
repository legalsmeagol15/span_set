package spanset

import (
	"fmt"
)

type SpanOptions byte

const (
	InfNegative SpanOptions = 1 << iota
	SpanStart
	SpanEnd
	InfPositive
	bothEnds = SpanStart | SpanEnd
	infinite = bothEnds | InfNegative | InfPositive
)

type Span[T ordered] interface {
	GetStart() T
	GetEnd() T

	IncludeStart() bool
	IncludeEnd() bool
	IncludeBefore() bool
	IncludeAfter() bool

	IsSingleton() bool
	IsUniversal() bool
	IsEmpty() bool
}

type span[T ordered] struct {
	start, end T
	includes   byte
}

func (s span[T]) GetStart() T          { return (s.start) }
func (s span[T]) GetEnd() T            { return (s.end) }
func (s *span[T]) IncludeBefore() bool { return (s.includes & byte(InfNegative)) != 0 }
func (s *span[T]) IncludeStart() bool  { return (s.includes & byte(SpanStart)) != 0 }
func (s *span[T]) IncludeEnd() bool    { return (s.includes & byte(SpanEnd)) != 0 }
func (s *span[T]) IncludeAfter() bool  { return (s.includes & byte(InfPositive)) != 0 }
func (s *span[T]) IsSingleton() bool   { return ((s.includes & byte(bothEnds)) != 0) && s.start == s.end }
func (s *span[T]) IsUniversal() bool   { return (s.includes & byte(infinite)) != 0 }
func (s *span[T]) IsEmpty() bool       { return (s.includes & byte(infinite)) == 0 }

func newIntSpan[T ordered](start T, end T, buildInclusions SpanOptions) span[int] {
	if a, ok := any(start).(int); ok {
		if b, ok := any(end).(int); ok {
			s := span[int]{
				start:    a,
				end:      b,
				includes: byte(buildInclusions),
			}
			return s
		}
	}
	panic(fmt.Sprintf("invalid span types: %v, %v", start, end))
}
