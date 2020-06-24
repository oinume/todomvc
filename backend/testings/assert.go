package testings

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
)

var (
	// These options are applied to cmp.Diff when there are no options specified.
	defaultCmpOptions = []cmp.Option{
		// compare if two protocol buffer messages are equal
		cmp.Comparer(proto.Equal),
	}
)

// AssertEqual ensures that got and want are equal by cmp.Diff.
// If they are not equal, it reports failure by t.Errorf with given `message`.
// If options are empty, it applies default options, which are defined in `defaultCmpOptions`.
// If any options are given, `defaultCmpOptions` are not applied.
func AssertEqual(t *testing.T, want, got interface{}, message string, options ...cmp.Option) {
	t.Helper()
	// The default options are applied only when there are no given options
	// because there is a problem that cmpopts.IgnoreFields in the given options for protbuf is ignored.
	// We need to wait until protobuf apiv2 is introduced.
	if len(options) == 0 {
		options = append(options, defaultCmpOptions...)
	}

	if diff := cmp.Diff(want, got, options...); diff != "" {
		if message == "" {
			t.Errorf("AssertEqual failed (-want +got):\n%s", diff)
		} else {
			t.Errorf("AssertEqual failed: %q: (-want +got):\n%s", message, diff)
		}
	}
}

// RequireEqual ensures that got and want are equal by cmp.Diff.
// If they are not equal, it reports failure by t.Fatalf with given `message`.
// If options are empty, it applies default options, which are defined in `defaultCmpOptions`.
// If any options are given, `defaultCmpOptions` are not applied.
func RequireEqual(t *testing.T, want, got interface{}, message string, options ...cmp.Option) {
	t.Helper()
	// The default options are applied only when there are no given options
	// because there is a problem that cmpopts.IgnoreFields in the given options for protbuf is ignored.
	// We need to wait until protobuf apiv2 is introduced.
	if len(options) == 0 {
		options = append(options, defaultCmpOptions...)
	}
	if diff := cmp.Diff(want, got, options...); diff != "" {
		if message == "" {
			t.Fatalf("RequireEqual failed (-want +got):\n%s", diff)
		} else {
			t.Fatalf("RequireEqual failed: %q: (-want +got):\n%s", message, diff)
		}
	}
}
