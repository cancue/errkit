package errkit

import (
	"runtime"
	"testing"
)

var initpc = caller()

type X struct{}

// val returns a Frame pointing to itself.
func (x X) val() Frame {
	return caller()
}

// ptr returns a Frame pointing to itself.
func (x *X) ptr() Frame {
	return caller()
}

func TestFrameFormat(t *testing.T) {
	var tests = []struct {
		Frame
		format string
		want   string
	}{{
		initpc,
		"%s",
		"stack_test.go",
	}, {
		initpc,
		"%+s",
		"github.com/cancue/errkit/errkit.init\n" +
			"\t.+/errkit/stack_test.go",
	}, {
		0,
		"%s",
		"unknown",
	}, {
		0,
		"%+s",
		"unknown",
	}, {
		initpc,
		"%d",
		"8",
	}, {
		0,
		"%d",
		"0",
	}, {
		initpc,
		"%n",
		"init",
	}, {
		func() Frame {
			var x X
			return x.ptr()
		}(),
		"%n",
		`\(\*X\).ptr`,
	}, {
		func() Frame {
			var x X
			return x.val()
		}(),
		"%n",
		"X.val",
	}, {
		0,
		"%n",
		"",
	}, {
		initpc,
		"%v",
		"stack_test.go:8",
	}, {
		initpc,
		"%+v",
		"github.com/cancue/errkit/errkit.init\n" +
			"\t.+/errkit/errkit/stack_test.go:8",
	}, {
		0,
		"%v",
		"unknown:0",
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.Frame, tt.format, tt.want)
	}
}

func TestFuncname(t *testing.T) {
	tests := []struct {
		name, want string
	}{
		{"", ""},
		{"runtime.main", "main"},
		{"github.com/cancue/errkit/errkit.funcname", "funcname"},
		{"funcname", "funcname"},
		{"io.copyBuffer", "copyBuffer"},
		{"main.(*R).Write", "(*R).Write"},
	}

	for _, tt := range tests {
		got := funcname(tt.name)
		want := tt.want
		if got != want {
			t.Errorf("funcname(%q): want: %q, got %q", tt.name, want, got)
		}
	}
}

func TestStackTrace(t *testing.T) {
	tests := []struct {
		err  error
		want []string
	}{{
		New("ooh"), []string{
			"github.com/cancue/errkit/errkit.TestStackTrace\n" +
				"\t.+/errkit/stack_test.go:120",
		},
	}, {
		Wrap(New("ooh"), "ahh"), []string{
			"github.com/cancue/errkit/errkit.TestStackTrace\n" +
				"\t.+/errkit/stack_test.go:125", // this is the stack of Wrap, not New
		},
	}, {
		func() error { return New("ooh") }(), []string{
			`github.com/cancue/errkit/errkit.TestStackTrace.func1` +
				"\n\t.+/errkit/stack_test.go:130", // this is the stack of New
			"github.com/cancue/errkit/errkit.TestStackTrace\n" +
				"\t.+/errkit/stack_test.go:130", // this is the stack of New's caller
		},
	}}
	for i, tt := range tests {
		x, ok := tt.err.(StackTracer)
		if !ok {
			t.Errorf("expected %#v to implement StackTrace() StackTrace", tt.err)
			continue
		}
		st := x.StackTrace()
		for j, want := range tt.want {
			testFormatRegexp(t, i, st[j], "%+v", want)
		}
	}
}

func stackTrace() StackTrace {
	const depth = 8
	var pcs [depth]uintptr
	n := runtime.Callers(1, pcs[:])
	var st stack = pcs[0:n]
	return st.StackTrace()
}

func TestStackTraceFormat(t *testing.T) {
	tests := []struct {
		StackTrace
		format string
		want   string
	}{{
		nil,
		"%s",
		`\[\]`,
	}, {
		nil,
		"%v",
		`\[\]`,
	}, {
		nil,
		"%+v",
		"",
	}, {
		nil,
		"%#v",
		`\[\]errkit.Frame\(nil\)`,
	}, {
		make(StackTrace, 0),
		"%s",
		`\[\]`,
	}, {
		make(StackTrace, 0),
		"%v",
		`\[\]`,
	}, {
		make(StackTrace, 0),
		"%+v",
		"",
	}, {
		make(StackTrace, 0),
		"%#v",
		`\[\]errkit.Frame{}`,
	}, {
		stackTrace()[:2],
		"%s",
		`\[stack_test.go stack_test.go\]`,
	}, {
		stackTrace()[:2],
		"%v",
		`\[stack_test.go:153 stack_test.go:200\]`,
	}, {
		stackTrace()[:2],
		"%+v",
		"\n" +
			"github.com/cancue/errkit/errkit.stackTrace\n" +
			"\t.+/errkit/stack_test.go:153\n" +
			"github.com/cancue/errkit/errkit.TestStackTraceFormat\n" +
			"\t.+/errkit/stack_test.go:204",
	}, {
		stackTrace()[:2],
		"%#v",
		`\[\]errkit.Frame{stack_test.go:153, stack_test.go:212}`,
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.StackTrace, tt.format, tt.want)
	}
}

// a version of runtime.Caller that returns a Frame, not a uintptr.
func caller() Frame {
	var pcs [3]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	frame, _ := frames.Next()
	return Frame(frame.PC)
}
