package assert

import (
	"errors"
	"fmt"
	"testing"
)

func TestAssertions(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		Equal(t, 1, 1, "%v aren't equal", Diff(), Fmt("ints"))
		NotEqual(t, 1, 2, "")
	})

	t.Run("nil", func(t *testing.T) {
		Nil(t, nil, "")
		Nil(t, (*chan int)(nil), "")
		Nil(t, (*func())(nil), "")
		Nil(t, (*map[int]int)(nil), "")
		Nil(t, (*pair)(nil), "")
		Nil(t, (*[]int)(nil), "")

		NotNil(t, make(chan int), "")
		NotNil(t, func() {}, "")
		NotNil(t, any(1), "")
		NotNil(t, make(map[int]int), "")
		NotNil(t, &pair{}, "")
		NotNil(t, make([]int, 0), "")

		NotNil(t, "foo", "")
		NotNil(t, 0, "")
		NotNil(t, false, "")
		NotNil(t, pair{}, "")
	})

	t.Run("zero", func(t *testing.T) {
		var n *int
		Zero(t, n, "")
		var p pair
		Zero(t, p, "")
		var null *pair
		Zero(t, null, "")
		var s []int
		Zero(t, s, "")
		var m map[string]string
		Zero(t, m, "")
		NotZero(t, 3, "")
	})

	t.Run("error chain", func(t *testing.T) {
		want := errors.New("base error")
		ErrorIs(t, fmt.Errorf("context: %w", want), want, "")
	})

	t.Run("unexported fields", func(t *testing.T) {
		// Two pairs differ only in an unexported field.
		p1 := pair{1, 2}
		p2 := pair{1, 3}
		NotEqual(t, p1, p2, "")
	})

	t.Run("regexp", func(t *testing.T) {
		Match(t, "foobar", `^foo`, "")
	})

	t.Run("panics", func(t *testing.T) {
		Panics(t, func() { panic("testing") }, "")
	})
}

type pair struct {
	First, Second int
}
