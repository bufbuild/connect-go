// Package assert is a minimal assert package using generics.
//
// This prevents connect from needing additional dependencies.
package assert

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

// Equal asserts that two values are equal.
func Equal[T any](t testing.TB, got, want T, message string, opts ...Option) bool {
	t.Helper()
	params := newParams(got, want, message, opts...)
	if cmp.Equal(got, want, params.cmpOpts...) {
		return true
	}
	report(t, params, "assert.Equal", true /* showWant */)
	return false
}

// NotEqual asserts that two values aren't equal.
func NotEqual[T any](t testing.TB, got, want T, message string, opts ...Option) bool {
	t.Helper()
	params := newParams(got, want, message, opts...)
	if !cmp.Equal(got, want, params.cmpOpts...) {
		return true
	}
	report(t, params, "assert.NotEqual", true /* showWant */)
	return false
}

// Nil asserts that the value is nil.
func Nil(t testing.TB, got any, message string, opts ...Option) bool {
	t.Helper()
	if isNil(got) {
		return true
	}
	params := newParams(got, nil, message, opts...)
	report(t, params, "assert.Nil", false /* showWant */)
	return false
}

// NotNil asserts that the value isn't nil.
func NotNil(t testing.TB, got any, message string, opts ...Option) bool {
	t.Helper()
	if !isNil(got) {
		return true
	}
	params := newParams(got, nil, message, opts...)
	report(t, params, "assert.NotNil", false /* showWant */)
	return false
}

// Zero asserts that the value is its type's zero value.
func Zero[T any](t testing.TB, got T, message string, opts ...Option) bool {
	t.Helper()
	var want T
	params := newParams(got, want, message, opts...)
	if cmp.Equal(got, want, params.cmpOpts...) {
		return true
	}
	report(t, params, fmt.Sprintf("assert.Zero (type %T)", got), false /* showWant */)
	return false
}

// NotZero asserts that the value is non-zero.
func NotZero[T any](t testing.TB, got T, message string, opts ...Option) bool {
	t.Helper()
	var want T
	params := newParams(got, want, message, opts...)
	if !cmp.Equal(got, want, params.cmpOpts...) {
		return true
	}
	report(t, params, fmt.Sprintf("assert.NotZero (type %T)", got), false /* showWant */)
	return false
}

// Match asserts that the value matches a regexp.
func Match(t testing.TB, got, want, message string, opts ...Option) bool {
	t.Helper()
	re, err := regexp.Compile(want)
	if err != nil {
		t.Fatalf("invalid regexp %q: %v", want, err)
	}
	if re.MatchString(got) {
		return true
	}
	params := newParams(got, want, message, opts...)
	report(t, params, "assert.Match", true /* showWant */)
	return false
}

// ErrorIs asserts that "want" is in "got's" error chain. See the standard
// library's errors package for details on error chains. On failure, output is
// identical to Equal.
func ErrorIs(t testing.TB, got, want error, message string, opts ...Option) bool {
	t.Helper()
	if errors.Is(got, want) {
		return true
	}
	params := newParams(got, want, message, opts...)
	report(t, params, "assert.ErrorIs", true /* showWant */)
	return false
}

// False asserts that "got" is false.
func False(t testing.TB, got bool, message string, opts ...Option) bool {
	t.Helper()
	if !got {
		return true
	}
	params := newParams(got, false, message, opts...)
	report(t, params, "assert.False", false)
	return false
}

// True asserts that "got" is true.
func True(t testing.TB, got bool, message string, opts ...Option) bool {
	t.Helper()
	if got {
		return true
	}
	params := newParams(got, false, message, opts...)
	report(t, params, "assert.True", false)
	return false
}

// Panics asserts that the function called panics.
func Panics(t testing.TB, panicker func(), message string, opts ...Option) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			params := newParams("no panic", "panic", message, opts...)
			report(t, params, "assert.Panic", false)
		}
	}()
	panicker()
}

// An Option modifies assertions.
type Option interface {
	apply(*params)
}

// Fmt treats the assertion message as a template, using fmt.Sprintf and the
// supplied arguments to expand it. For example,
//   assert.Equal(t, 0, 1, "failed to parse %q", assert.Fmt("foobar"))
// will print the message
//   failed to parse "foobar"
func Fmt(args ...any) Option {
	return optionFunc(func(p *params) {
		if len(args) > 0 {
			p.message = fmt.Sprintf(p.message, args...)
		}
	})
}

// Diff prints a diff between "got" and "want" on failures.
func Diff() Option {
	return optionFunc(func(p *params) {
		p.printDiff = true
	})
}

type optionFunc func(*params)

func (f optionFunc) apply(p *params) { f(p) }

type params struct {
	got       any
	want      any
	cmpOpts   []cmp.Option // user-supplied equality configuration
	message   string       // user-supplied description of failure
	printDiff bool         // include diff in failure output
}

func newParams(got, want any, message string, opts ...Option) *params {
	p := &params{
		got:     got,
		want:    want,
		cmpOpts: []cmp.Option{protocmp.Transform()},
		message: message,
	}
	for _, opt := range opts {
		opt.apply(p)
	}
	if p.got == nil || p.want == nil {
		// diff panics on nil
		p.printDiff = false
	}
	return p
}

func report(t testing.TB, params *params, desc string, showWant bool) {
	t.Helper()
	w := &bytes.Buffer{}
	if params.message != "" {
		w.WriteString(params.message)
	}
	w.WriteString("\n")
	fmt.Fprintf(w, "assertion:\t%s\n", desc)
	fmt.Fprintf(w, "got:\t%+v\n", params.got)
	if showWant {
		fmt.Fprintf(w, "want:\t%+v\n", params.want)
	}
	if params.printDiff {
		fmt.Fprintf(w, "\ndiff (-want, +got):\n%v", cmp.Diff(params.want, params.got))
	}
	t.Fatal(w.String())
}

func isNil(got any) bool {
	// Simple case, true only when the user directly passes a literal nil.
	if got == nil {
		return true
	}
	// Possibly more complex. Interfaces are a pair of words: a pointer to a type
	// and a pointer to a value. Because we're passing got as an interface, it's
	// likely that we've gotten a non-nil type and a nil value. This makes got
	// itself non-nil, but the user's code passed a nil value.
	val := reflect.ValueOf(got)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return val.IsNil()
	default:
		return false
	}
}
