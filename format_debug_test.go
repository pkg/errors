package errors

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"
)

func TestDebugFormatNew(t *testing.T) {
	testDebugs := []struct {
		error
		format string
		want   string
	}{{
		New("error"),
		"%s",
		"error",
	}, {
		New("error"),
		"%v",
		"error\n" +
			"github.com/pkg/errors.TestDebugFormatNew\n" +
			"\t.+/github.com/pkg/errors/format_debug_test.go:22",
	}, {
		New("error"),
		"%+v",
		"error\n" +
			"github.com/pkg/errors.TestDebugFormatNew\n" +
			"\t.+/github.com/pkg/errors/format_debug_test.go:28",
	}, {
		New("error"),
		"%q",
		`"error"`,
	}}

	for i, tt := range testDebugs {
		testDebugFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestDebugFormatErrorf(t *testing.T) {
	testDebugs := []struct {
		error
		format string
		want   string
	}{{
		Errorf("%s", "error"),
		"%s",
		"error",
	}, {
		Errorf("%s", "error"),
		"%v",
		"error\n" +
			"github.com/pkg/errors.TestDebugFormatErrorf\n" +
			"\t.+/github.com/pkg/errors/format_debug_test.go:54",
	}, {
		Errorf("%s", "error"),
		"%+v",
		"error\n" +
			"github.com/pkg/errors.TestDebugFormatErrorf\n" +
			"\t.+/github.com/pkg/errors/format_debug_test.go:60",
	}}

	for i, tt := range testDebugs {
		testDebugFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestDebugFormatWrap(t *testing.T) {
	testDebugs := []struct {
		error
		format string
		want   string
	}{{
		Wrap(New("error"), "error2"),
		"%s",
		"error2: error",
	}, {
		Wrap(New("error"), "error2"),
		"%v",
		"error\n" +
			"github.com/pkg/errors.TestDebugFormatWrap\n" +
			"\t.+/github.com/pkg/errors/format_debug_test.go:82",
	}, {
		Wrap(New("error"), "error2"),
		"%+v",
		"error\n" +
			"github.com/pkg/errors.TestDebugFormatWrap\n" +
			"\t.+/github.com/pkg/errors/format_debug_test.go:88",
	}, {
		Wrap(io.EOF, "error"),
		"%s",
		"error: EOF",
	}, {
		Wrap(io.EOF, "error"),
		"%v",
		"EOF\n" +
			"error\n" +
			"github.com/pkg/errors.TestDebugFormatWrap\n" +
			"\t.+/github.com/pkg/errors/format_debug_test.go:98",
	}, {
		Wrap(io.EOF, "error"),
		"%+v",
		"EOF\n" +
			"error\n" +
			"github.com/pkg/errors.TestDebugFormatWrap\n" +
			"\t.+/github.com/pkg/errors/format_debug_test.go:105",
	}, {
		Wrap(Wrap(io.EOF, "error1"), "error2"),
		"%+v",
		"EOF\n" +
			"error1\n" +
			"github.com/pkg/errors.TestDebugFormatWrap\n" +
			"\t.+/github.com/pkg/errors/format_debug_test.go:112\n",
	}, {
		Wrap(New("error with space"), "context"),
		"%q",
		`"context: error with space"`,
	}}

	for i, tt := range testDebugs {
		testDebugFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestDebugFormatWrapf(t *testing.T) {
	testDebugs := []struct {
		error
		format string
		want   string
	}{{
		Wrapf(io.EOF, "error%d", 2),
		"%s",
		"error2: EOF",
	}, {
		Wrapf(io.EOF, "error%d", 2),
		"%v",
		"EOF\n" +
			"error2\n" +
			"github.com/pkg/errors.TestDebugFormatWrapf\n" +
			"\t.+/github.com/pkg/errors/format_debug_test.go:139",
	}, {
		Wrapf(io.EOF, "error%d", 2),
		"%+v",
		"EOF\n" +
			"error2\n" +
			"github.com/pkg/errors.TestDebugFormatWrapf\n" +
			"\t.+/github.com/pkg/errors/format_debug_test.go:146",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%s",
		"error2: error",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%v",
		"error\n" +
			"github.com/pkg/errors.TestDebugFormatWrapf\n" +
			"\t.+/github.com/pkg/errors/format_debug_test.go:157",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%+v",
		"error\n" +
			"github.com/pkg/errors.TestDebugFormatWrapf\n" +
			"\t.+/github.com/pkg/errors/format_debug_test.go:163",
	}}

	for i, tt := range testDebugs {
		testDebugFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestDebugFormatWithStack(t *testing.T) {
	testDebugs := []struct {
		error
		format string
		want   []string
	}{{
		WithStack(io.EOF),
		"%s",
		[]string{"EOF"},
	}, {
		WithStack(io.EOF),
		"%v",
		[]string{"EOF",
			"github.com/pkg/errors.TestDebugFormatWithStack\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:185"},
	}, {
		WithStack(io.EOF),
		"%+v",
		[]string{"EOF",
			"github.com/pkg/errors.TestDebugFormatWithStack\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:191"},
	}, {
		WithStack(New("error")),
		"%s",
		[]string{"error"},
	}, {
		WithStack(New("error")),
		"%v",
		[]string{"error",
			"github.com/pkg/errors.TestDebugFormatWithStack\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:201",
			"github.com/pkg/errors.TestDebugFormatWithStack\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:201"},
	}, {
		WithStack(New("error")),
		"%+v",
		[]string{"error",
			"github.com/pkg/errors.TestDebugFormatWithStack\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:209",
			"github.com/pkg/errors.TestDebugFormatWithStack\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:209"},
	}, {
		WithStack(WithStack(io.EOF)),
		"%+v",
		[]string{"EOF",
			"github.com/pkg/errors.TestDebugFormatWithStack\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:217",
			"github.com/pkg/errors.TestDebugFormatWithStack\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:217"},
	}, {
		WithStack(WithStack(Wrapf(io.EOF, "message"))),
		"%+v",
		[]string{"EOF",
			"message",
			"github.com/pkg/errors.TestDebugFormatWithStack\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:225",
			"github.com/pkg/errors.TestDebugFormatWithStack\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:225",
			"github.com/pkg/errors.TestDebugFormatWithStack\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:225"},
	}, {
		WithStack(Errorf("error%d", 1)),
		"%+v",
		[]string{"error1",
			"github.com/pkg/errors.TestDebugFormatWithStack\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:236",
			"github.com/pkg/errors.TestDebugFormatWithStack\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:236"},
	}}

	for i, tt := range testDebugs {
		testDebugFormatCompleteCompare(t, i, tt.error, tt.format, tt.want, true)
	}
}

func TestDebugFormatWithMessage(t *testing.T) {
	testDebugs := []struct {
		error
		format string
		want   []string
	}{{
		WithMessage(New("error"), "error2"),
		"%s",
		[]string{"error2: error"},
	}, {
		WithMessage(New("error"), "error2"),
		"%v",
		[]string{
			"error",
			"github.com/pkg/errors.TestDebugFormatWithMessage\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:260",
			"error2"},
	}, {
		WithMessage(New("error"), "error2"),
		"%+v",
		[]string{
			"error",
			"github.com/pkg/errors.TestDebugFormatWithMessage\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:268",
			"error2"},
	}, {
		WithMessage(io.EOF, "addition1"),
		"%s",
		[]string{"addition1: EOF"},
	}, {
		WithMessage(io.EOF, "addition1"),
		"%v",
		[]string{"EOF", "addition1"},
	}, {
		WithMessage(io.EOF, "addition1"),
		"%+v",
		[]string{"EOF", "addition1"},
	}, {
		WithMessage(WithMessage(io.EOF, "addition1"), "addition2"),
		"%v",
		[]string{"EOF", "addition1", "addition2"},
	}, {
		WithMessage(WithMessage(io.EOF, "addition1"), "addition2"),
		"%+v",
		[]string{"EOF", "addition1", "addition2"},
	}, {
		Wrap(WithMessage(io.EOF, "error1"), "error2"),
		"%+v",
		[]string{"EOF", "error1", "error2",
			"github.com/pkg/errors.TestDebugFormatWithMessage\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:296"},
	}, {
		WithMessage(Errorf("error%d", 1), "error2"),
		"%+v",
		[]string{"error1",
			"github.com/pkg/errors.TestDebugFormatWithMessage\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:302",
			"error2"},
	}, {
		WithMessage(WithStack(io.EOF), "error"),
		"%+v",
		[]string{
			"EOF",
			"github.com/pkg/errors.TestDebugFormatWithMessage\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:309",
			"error"},
	}, {
		WithMessage(Wrap(WithStack(io.EOF), "inside-error"), "outside-error"),
		"%+v",
		[]string{
			"EOF",
			"github.com/pkg/errors.TestDebugFormatWithMessage\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:317",
			"inside-error",
			"github.com/pkg/errors.TestDebugFormatWithMessage\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:317",
			"outside-error"},
	}}

	for i, tt := range testDebugs {
		testDebugFormatCompleteCompare(t, i, tt.error, tt.format, tt.want, true)
	}
}

func TestDebugFormatGeneric(t *testing.T) {
	starts := []struct {
		err  error
		want []string
	}{
		{New("new-error"), []string{
			"new-error",
			"github.com/pkg/errors.TestDebugFormatGeneric\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:339"},
		}, {Errorf("errorf-error"), []string{
			"errorf-error",
			"github.com/pkg/errors.TestDebugFormatGeneric\n" +
				"\t.+/github.com/pkg/errors/format_debug_test.go:343"},
		}, {errors.New("errors-new-error"), []string{
			"errors-new-error"},
		},
	}

	debugWrappers := []debugWrapper{
		{
			func(err error) error { return WithMessage(err, "with-message") },
			[]string{"with-message"},
		}, {
			func(err error) error { return WithStack(err) },
			[]string{
				"github.com/pkg/errors.(func·002|TestDebugFormatGeneric.func2)\n\t" +
					".+/github.com/pkg/errors/format_debug_test.go:357",
			},
		}, {
			func(err error) error { return Wrap(err, "wrap-error") },
			[]string{
				"wrap-error",
				"github.com/pkg/errors.(func·003|TestDebugFormatGeneric.func3)\n\t" +
					".+/github.com/pkg/errors/format_debug_test.go:363",
			},
		}, {
			func(err error) error { return Wrapf(err, "wrapf-error%d", 1) },
			[]string{
				"wrapf-error1",
				"github.com/pkg/errors.(func·004|TestDebugFormatGeneric.func4)\n\t" +
					".+/github.com/pkg/errors/format_debug_test.go:370",
			},
		},
	}

	for s := range starts {
		err := starts[s].err
		want := starts[s].want
		testDebugFormatCompleteCompare(t, s, err, "%+v", want, false)
		testDebugGenericRecursive(t, err, want, debugWrappers, 3)
	}
}

func testDebugFormatRegexp(t *testing.T, n int, arg interface{}, format, want string) {
	Debug(true)
	got := fmt.Sprintf(format, arg)
	Debug(false)
	gotLines := strings.SplitN(got, "\n", -1)
	wantLines := strings.SplitN(want, "\n", -1)

	if len(wantLines) > len(gotLines) {
		t.Errorf("testDebug %d: wantLines(%d) > gotLines(%d):\n got: %q\nwant: %q", n+1, len(wantLines), len(gotLines), got, want)
		return
	}

	for i, w := range wantLines {
		match, err := regexp.MatchString(w, gotLines[i])
		if err != nil {
			t.Fatal(err)
		}
		if !match {
			t.Errorf("testDebug %d: line %d: fmt.Sprintf(%q, err):\n got: %q\nwant: %q", n+1, i+1, format, got, want)
		}
	}
}

func testDebugFormatCompleteCompare(t *testing.T, n int, arg interface{}, format string, want []string, detectStackBoundaries bool) {
	Debug(true)
	gotStr := fmt.Sprintf(format, arg)
	Debug(false)

	got, err := parseBlocks(gotStr, detectStackBoundaries)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != len(want) {
		t.Fatalf("testDebug %d: fmt.Sprintf(%s, err) -> wrong number of blocks: got(%d) want(%d)\n got: %s\nwant: %s\ngotStr: %q",
			n+1, format, len(got), len(want), prettyBlocks(got), prettyBlocks(want), gotStr)
	}

	for i := range got {
		if strings.ContainsAny(want[i], "\n") {
			// Match as stack
			match, err := regexp.MatchString(want[i], got[i])
			if err != nil {
				t.Fatal(err)
			}
			if !match {
				t.Fatalf("testDebug %d: block %d: fmt.Sprintf(%q, err):\ngot:\n%q\nwant:\n%q\nall-got:\n%s\nall-want:\n%s\n",
					n+1, i+1, format, got[i], want[i], prettyBlocks(got), prettyBlocks(want))
			}
		} else {
			// Match as message
			if got[i] != want[i] {
				t.Fatalf("testDebug %d: fmt.Sprintf(%s, err) at block %d got != want:\n got: %q\nwant: %q", n+1, format, i+1, got[i], want[i])
			}
		}
	}
}

type debugWrapper struct {
	wrap func(err error) error
	want []string
}

func testDebugGenericRecursive(t *testing.T, beforeErr error, beforeWant []string, list []debugWrapper, maxDepth int) {
	if len(beforeWant) == 0 {
		panic("beforeWant must not be empty")
	}
	for _, w := range list {
		if len(w.want) == 0 {
			panic("want must not be empty")
		}

		err := w.wrap(beforeErr)

		// Copy required cause append(beforeWant, ..) modified beforeWant subtly.
		beforeCopy := make([]string, len(beforeWant))
		copy(beforeCopy, beforeWant)

		beforeWant := beforeCopy
		last := len(beforeWant) - 1
		var want []string

		// Merge two stacks behind each other.
		if strings.ContainsAny(beforeWant[last], "\n") && strings.ContainsAny(w.want[0], "\n") {
			want = append(beforeWant[:last], append([]string{beforeWant[last] + "((?s).*)" + w.want[0]}, w.want[1:]...)...)
		} else {
			want = append(beforeWant, w.want...)
		}

		testDebugFormatCompleteCompare(t, maxDepth, err, "%+v", want, false)
		if maxDepth > 0 {
			testDebugGenericRecursive(t, err, want, list, maxDepth-1)
		}
	}
}
