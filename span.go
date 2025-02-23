package spanset

import (
	"fmt"
	"strings"

	"golang.org/x/exp/constraints"
)

var ()

type SpanOptions byte
type Ordered constraints.Ordered

const (
	None        SpanOptions = 0
	InfNegative SpanOptions = 1 << iota
	SpanStart   SpanOptions = 1 << iota
	SpanEnd     SpanOptions = 1 << iota
	InfPositive SpanOptions = 1 << iota

	bothEnds = SpanStart | SpanEnd
	infinite = bothEnds | InfNegative | InfPositive
)

var ()

type Span[T Ordered] interface {
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

type span[T Ordered] struct {
	start, end T
	includes   SpanOptions
}

func (s span[T]) GetStart() T         { return (s.start) }
func (s span[T]) GetEnd() T           { return (s.end) }
func (s span[T]) IncludeBefore() bool { return (s.includes & InfNegative) != 0 }
func (s span[T]) IncludeStart() bool  { return (s.includes & SpanStart) != 0 }
func (s span[T]) IncludeEnd() bool    { return (s.includes & SpanEnd) != 0 }
func (s span[T]) IncludeAfter() bool  { return (s.includes & InfPositive) != 0 }
func (s span[T]) IsSingleton() bool   { return ((s.includes & bothEnds) != 0) && s.start == s.end }
func (s span[T]) IsUniversal() bool   { return (s.includes & infinite) != 0 }
func (s span[T]) IsEmpty() bool       { return ((s.includes & infinite) == 0) && s.start == s.end }
func (s span[T]) IsInfNegative() bool { return ((s.includes & InfNegative) != 0) }
func (s span[T]) IsInfPositive() bool { return ((s.includes & InfPositive) != 0) }

func (s span[T]) String() string {
	var sb strings.Builder
	if s.IsInfNegative() {
		sb.WriteString("<-- ")
	}
	if s.IncludeStart() {
		sb.WriteString(fmt.Sprintf("%v", s.start))
	} else {
		sb.WriteString(fmt.Sprintf("(%v)", s.start))
	}
	if s.start != s.end {
		sb.WriteString(" --- ")
		if s.IncludeEnd() {
			sb.WriteString(fmt.Sprintf("%v", s.end))
		} else {
			sb.WriteString(fmt.Sprintf("(%v)", s.end))
		}
	}
	if s.IsInfPositive() {
		sb.WriteString(" -->")
	}
	return sb.String()
}

func includes(includeStart, IncludeEnd bool) SpanOptions {
	i := None
	if includeStart {
		i |= SpanStart
	}
	if IncludeEnd {
		i |= SpanEnd
	}
	return i
}

func to_span[T Ordered](s Span[T]) span[T] {
	if as_span, ok := s.(span[T]); ok {
		return as_span
	}
	panic("something ain't right")
}

func makeEmpty[T Ordered]() Span[T] {
	return span[T]{start: *new(T), end: *new(T), includes: None}
}
func makeUniversal[T Ordered]() Span[T] {
	return span[T]{start: *new(T), end: *new(T), includes: infinite}
}

func (s *span[T]) contains_singleton(item T) bool {
	if item > s.start && item < s.end {
		return true
	} else if item == s.start {
		return s.IncludeStart()
	} else if item == s.end {
		return s.IncludeEnd()
	} else {
		return false
	}
}
func (a *span[T]) intersection(b *span[T]) span[T] {
	/*
		1.0:	a-----a
						b-----b
		1.1:	a---a
					b---b
		1.2:	a-----a
					b-----b
		1.3:	a-------a
					b---b
		1.4:	a-------a
				  b---b

		2.2:	a---a
				b-------b
		2.3:	a---a
				b---b
		2.4:	a------a
				b---b

		3.2:	  a---a
				b-------b
		3.3			a---a
				b-------b
		3.4:	  a-----a
				b-----b

		4.4:		a---a
				b---b

		5.4:		   a---a
				b---b

	*/

	if a.end < b.start || b.end < a.start {
		// Rule out cases 1.0 and 5.4
	} else if a.start < b.start {
		// Cases 1.1-1.4
		if a.end == b.start {
			// Case 1.1
			if a.IncludeEnd() && b.IncludeStart() {
				return span[T]{
					start:    b.start,
					end:      a.end,
					includes: includes(a.IncludeStart(), b.IncludeEnd()),
				}
			}
		} else if a.end < b.end {
			// Case 1.2
		}

	}
	return to_span(makeEmpty[T]())
}
func (a *span[T]) union(b *span[T]) (span[T], span[T]) {
	// Check out the different cases under intersection()
	if a.end < b.start {
		// Case 0a
		return *a, *b
	} else if b.end < a.start {
		// Case 2a
		return *b, *a
	} else if a.start < b.start {
		// Case 0b
		if a.end == b.start {
			if a.IncludeEnd() || b.IncludeStart() {
				return span[T]{
					start:    a.start,
					end:      b.end,
					includes: includes(a.IncludeStart(), b.IncludeEnd()),
				}, to_span(makeEmpty[T]())
			} else {
				return *a, *b
			}
		} else if a.end > b.end {
			// Case 0c
			return *a, to_span(makeEmpty[T]())
		} else if a.end > b.start {
			// Case 0d
			return span[T]{
				start:    a.start,
				end:      b.end,
				includes: includes(a.IncludeStart(), b.IncludeEnd()),
			}, to_span(makeEmpty[T]())
		} else if a.end == b.end {
			// Case 0e
			return span[T]{
				start:    a.start,
				end:      a.end,
				includes: includes(a.IncludeStart(), a.IncludeEnd() || b.IncludeEnd()),
			}, to_span(makeEmpty[T]())
		}
	} else if a.start > b.start {
		if b.end == a.start {
			// Case 2b
			if b.IncludeEnd() || a.IncludeStart() {
				return span[T]{
					start:    b.start,
					end:      b.end,
					includes: includes(b.IncludeStart(), a.IncludeEnd()),
				}, to_span(makeEmpty[T]())
			} else {
				return *b, *a
			}
		} else if a.end < b.end {
			// Case 2c
			return *b, to_span(makeEmpty[T]())
		} else if a.end > b.end {
			// Case 2d
			return span[T]{
				start:    b.start,
				end:      a.end,
				includes: includes(b.IncludeStart(), a.IncludeEnd()),
			}, to_span(makeEmpty[T]())
		} else if b.end == a.end {
			// Case 2e
			return span[T]{
				start:    b.start,
				end:      b.end,
				includes: includes(b.IncludeStart(), a.IncludeEnd() || b.IncludeEnd()),
			}, to_span(makeEmpty[T]())
		}
	} else if a.start == b.start {
		if a.end < b.end {
			return span[T]{
				start:    a.start,
				end:      b.end,
				includes: includes(a.IncludeStart() || b.IncludeStart(), b.IncludeEnd()),
			}, to_span(makeEmpty[T]())
		} else if a.end > b.end {
			return span[T]{
				start:    b.start,
				end:      a.end,
				includes: includes(a.IncludeStart() || b.IncludeStart(), a.IncludeEnd()),
			}, to_span(makeEmpty[T]())
		} else if a.end == b.end {
			return span[T]{
				start:    a.start,
				end:      a.end,
				includes: includes(a.IncludeStart() || b.IncludeStart(), a.IncludeEnd() || b.IncludeEnd()),
			}, to_span(makeEmpty[T]())
		}
	}
	panic("This should never happen")
}
func (s *span[T]) inverse() (span[T], span[T]) {
	a := span[T]{
		start:    s.start,
		end:      s.start,
		includes: includes(!s.IncludeStart(), !s.IncludeStart()),
	}
	a.includes ^= InfNegative

	b := span[T]{
		start:    s.end,
		end:      s.end,
		includes: includes(!s.IncludeEnd(), !s.IncludeEnd()),
	}
	b.includes ^= InfPositive

	return a, b
}
