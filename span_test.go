package spanset

import (
	"testing"
)

func TestSpanFunctions_Case0a(t *testing.T) {
	/*
		case 0: a.start < b.start
			case 0a: <--a-->
							 <--b-->
	*/
	a := span[int]{
		start:    10,
		end:      30,
		includes: bothEnds,
	}
	b := span[int]{
		start:    31,
		end:      50,
		includes: bothEnds,
	}

	c := a.intersection(&b)
	if !c.IsEmpty() {
		t.Fatal("should have been empty")
	}

	c0, c1 := a.union(&b)
	if c1.IsEmpty() {
		t.Fatal("the second span should not have been empty")
	}
	if c0.IsSingleton() || c1.IsSingleton() {
		t.Fatal("both should be infinites, not singletons")
	}
}

func TestSpanFunctions_Case0b(t *testing.T) {
	/*
		case 0: a.start < b.start
			case 0b: <--a--|
						   |---b--->
	*/

	a := span[int]{
		start:    10,
		end:      30,
		includes: bothEnds,
	}
	b := span[int]{
		start:    30,
		end:      50,
		includes: bothEnds,
	}

	c := a.intersection(&b)
	if !c.IsSingleton() {
		t.Fatalf("should come back as a singleton: %v", c)
	} else if c.start != c.end || c.start != 30 {
		t.Fatalf("should be a singleton at 30: %v", c)
	}

	c0, c1 := a.union(&b)
	if !c1.IsEmpty() {
		t.Fatalf("second span should be empty")
	} else if c0.IsEmpty() {
		t.Fatalf("first span shouldn't be empty")
	}
	if c0.start != 10 || c0.end != 50 || !c.IncludeStart() || !c.IncludeEnd() || c.IsInfNegative() || c.IsInfPositive() {
		t.Fatalf("has the wrong shape: %v", c)
	}

	a.includes = SpanStart
	c0, c1 = a.union(&b)
	if !c1.IsEmpty() {
		t.Fatalf("second span should be empty")
	} else if c0.IsEmpty() {
		t.Fatalf("first span shouldn't be empty")
	}
	if c0.start != 10 || c0.end != 50 || !c.IncludeStart() || !c.IncludeEnd() || c.IsInfNegative() || c.IsInfPositive() {
		t.Fatalf("has the wrong shape: %v", c)
	}

	b.includes = SpanEnd
	c0, c1 = a.union(&b)
	if c0.IsEmpty() || c1.IsEmpty() {
		t.Fatalf("neither should be empty: %v    %v", c0, c1)
	}
	if c0 != a || c1 != b {
		t.Fatalf("should be same as originals: %v    %v", c0, c1)
	}
}

func TestSpanFunctions_Case0c(t *testing.T) {
	/*
		case 0: a.start < b.start
			case 0c: <------a------>
						 <--b-->
	*/
	a := span[int]{
		start:    10,
		end:      50,
		includes: bothEnds,
	}
	b := span[int]{
		start:    20,
		end:      40,
		includes: bothEnds,
	}

	c := a.intersection(&b)
	if c != b {
		t.Fatalf("should be identical to 2nd span: %v", c)
	}
	c0, c1 := a.union(&b)
	if !c1.IsEmpty() {
		t.Fatalf("second span should be empty: %v", c1)
	}
	if c0 != a {
		t.Fatalf("should be identical to 1st span: %v", c0)
	}
}

func TestSpanFunctions_Case0d(t *testing.T) {
	/*
		case 0: a.start < b.start
			case 0d: <--a-->
						<----b----->
	*/
	a := span[int]{
		start:    10,
		end:      30,
		includes: bothEnds,
	}
	b := span[int]{
		start:    20,
		end:      40,
		includes: bothEnds,
	}

	expected_and := span[int]{
		start:    20,
		end:      30,
		includes: bothEnds,
	}
	expected_or := span[int]{
		start:    10,
		end:      40,
		includes: bothEnds,
	}

	c := a.intersection(&b)
	if c != expected_and {
		t.Fatalf("wrong shape for intersection: %v", c)
	}

	d0, d1 := a.union(&b)
	if !d1.IsEmpty() {
		t.Fatalf("second span should be empty: %v", d1)
	} else if d0 != expected_or {
		t.Fatalf("wrong shape for union: %v", d0)
	}
}
