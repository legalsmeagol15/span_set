package spanset

import (
	"testing"
)

func TestBitmasking(t *testing.T) {
	if i := newIntSpan(1, 10, SpanStart); !i.IncludeStart() || i.IncludeEnd() || i.IncludeBefore() || i.IncludeAfter() {
		t.Errorf("bitmasking didn't work on IncludeStart")
	}

}
