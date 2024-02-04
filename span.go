package spanset

import (
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
func univConsecutive[T Ordered](a, b T) bool { return false }
func intConsecutive(a, b int) bool           { return a+1 == b }
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
		case 0: a.start < b.start
			case 0a: <--a-->
							 <--b-->
			case 0b: <--a--|
						   |---b--->
			case 0c: <------a------>
						 <--b-->
			case 0d: <--a-->
						<----b----->
			case 0e: <------a------|
							<--b---|
		case 1: a.start == b.start
			case 1a: |--a-->
					 |------b------>
			case 1b: |------a------|
					 |------b------|
			case 1c: |------a------>
					 |--b-->
		case 2: a.start > b.start
			case 2a: <--b-->
							 <--a-->
			case 2b: <--b--|
						   |---a--->
			case 2c: <------b------>
						<--a-->
			case 2d: <--b-->
						<----a----->
			case 2e: <------b------|
							<--a---|

	*/
	if a.end < b.start || b.end < a.start {
		// Rule out cases 0a and 2a
	} else if a.start < b.start {
		// Cases 0b-0e
		if a.end == b.start {
			// Case 0b
			if a.IncludeEnd() && b.IncludeStart() {
				return span[T]{
					start:    a.start,
					end:      a.end,
					includes: includes(a.IncludeStart(), b.IncludeEnd()),
				}
			}
		} else if a.end > b.end {
			// Case 0c
			return *b
		} else if a.end > b.start {
			// Case 0d
			return span[T]{
				start:    b.start,
				end:      a.end,
				includes: includes(b.IncludeStart(), a.IncludeEnd())}
		} else if a.end == b.end {
			// Case 0e
			return span[T]{
				start:    b.start,
				end:      a.end,
				includes: includes(b.IncludeStart(), a.IncludeEnd() && b.IncludeEnd()),
			}
		}
	} else if a.start > b.start {
		// Cases 2b - 2e
		if a.start == b.end {
			// Case 2b
			if a.IncludeStart() && b.IncludeEnd() {
				return span[T]{
					start:    b.end,
					end:      a.start,
					includes: includes(b.IncludeStart(), a.IncludeEnd()),
				}
			}
		} else if a.end < b.end {
			// Case 2c
			return *a
		} else if a.start < b.end {
			// Case 2d
			return span[T]{
				start:    a.start,
				end:      b.end,
				includes: includes(b.IncludeEnd(), a.IncludeStart()),
			}
		} else if a.end == b.end {
			// Case 2e
			return span[T]{
				start:    a.start,
				end:      b.end,
				includes: includes(a.IncludeStart(), a.IncludeEnd() && b.IncludeEnd())}
		}
	} else if a.start == b.start {
		if a.end < b.end {
			// Case 1a
			return *a
		} else if a.end == b.end {
			// Case 1b
			return span[T]{
				start:    a.start,
				end:      a.end,
				includes: a.includes & b.includes}
		} else if a.end > b.end {
			// Case 1c
			return *b
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
