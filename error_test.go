package connect

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/bufconnect/connect/internal/assert"
)

func TestErrorCodeOK(t *testing.T) {
	// Must use == rather than assert.Equal to avoid being fooled by
	// typed nils.
	assert.True(t, Wrap(CodeOK, errors.New("ok")) == nil, "wrap code ok")
	assert.True(t, Errorf(CodeOK, "ok") == nil, "errorf code ok")
}

func TestErrorFormatting(t *testing.T) {
	assert.Equal(t, Errorf(CodeUnavailable, "").Error(), CodeUnavailable.String(), "no message")
	text := Errorf(CodeUnavailable, "foo").Error()
	assert.True(t, strings.Contains(text, CodeUnavailable.String()), "error text should include code")
	assert.True(t, strings.Contains(text, "foo"), "error text should include message")
}

func TestErrorCode(t *testing.T) {
	err := fmt.Errorf("another: %w", Errorf(CodeUnavailable, "foo"))
	cerr, ok := AsError(err)
	assert.True(t, ok, "extract connect error")
	assert.Equal(t, cerr.Code(), CodeUnavailable, "extracted code")
}

func TestCodeOf(t *testing.T) {
	assert.Equal(t, CodeOf(nil), CodeOK, "nil error code")
	assert.Equal(t, CodeOf(Errorf(CodeUnavailable, "foo")), CodeUnavailable, "explicitly-set code")
	assert.Equal(t, CodeOf(errors.New("foo")), CodeUnknown, "fallback code")
}

func TestErrorDetails(t *testing.T) {
	second := durationpb.New(time.Second)
	detail, err := anypb.New(second)
	assert.Nil(t, err, "create anypb.Any")
	rerr := Errorf(CodeUnknown, "details")
	assert.Zero(t, rerr.Details(), "fresh error")
	rerr.AddDetail(detail)
	assert.Equal(t, rerr.Details(), []ErrorDetail{detail}, "retrieve details")
}
