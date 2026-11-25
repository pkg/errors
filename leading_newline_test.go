package errors

import (
	"fmt"
	"strings"
	"testing"
)

// Regression test for https://github.com/pkg/errors/issues/230:
// formatting a StackTrace directly with %+v must not produce a leading
// newline before the first frame.
func TestStackTraceFormatNoLeadingNewline(t *testing.T) {
	st := stackTrace()[:2]
	got := fmt.Sprintf("%+v", st)
	if strings.HasPrefix(got, "\n") {
		t.Errorf("fmt.Sprintf(\"%%+v\", stackTrace) has leading newline:\n%q", got)
	}

	// Frames must still be separated by newlines.
	first := fmt.Sprintf("%+v", st[0])
	both := fmt.Sprintf("%+v", st[:1])
	if len(st) > 1 && !strings.Contains(got, first) {
		t.Errorf("expected output to contain first frame %q, got %q", first, both)
	}
}
