package templatte_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"github.com/andrewwkimm/templatte-go"
)

func TestGreet(t *testing.T) {
	t.Parallel()

	got, err := templatte.Greet("Gopher")

	require.NoError(t, err)
	want := templatte.Greeting{Recipient: "Gopher", Message: "Hello, Gopher!"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Greet(%q) mismatch (-want +got):\n%s", "Gopher", diff)
	}
}

func TestGreet_EmptyName(t *testing.T) {
	t.Parallel()

	_, err := templatte.Greet("")

	require.ErrorIs(t, err, templatte.ErrEmptyName)
}
