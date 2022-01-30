//go:build go1.18

package main

import (
	"testing"

	"github.com/bradenaw/juniper/internal/require2"
)

func TestNonLocalSymbolLink(t *testing.T) {
	require2.Equal(
		t,
		nonLocalSymbolLink(
			"github.com/bradenaw/juniper/container/tree",
			"github.com/bradenaw/juniper/iterator",
			"Iterator",
		),
		"../iterator.html#Iterator",
	)
	require2.Equal(
		t,
		nonLocalSymbolLink(
			"github.com/bradenaw/juniper/container/tree",
			"github.com/bradenaw/juniper/internal/fuzz",
			"Operations",
		),
		"../internal/fuzz.html#Operations",
	)
	require2.Equal(
		t,
		nonLocalSymbolLink(
			"github.com/bradenaw/juniper/xmath",
			"github.com/bradenaw/juniper/xmath/xrand",
			"Sample",
		),
		"xmath/xrand.html#Sample",
	)
	require2.Equal(
		t,
		nonLocalSymbolLink(
			"github.com/bradenaw/juniper/xmath/xrand",
			"github.com/bradenaw/juniper/xmath",
			"Min",
		),
		"../xmath.html#Min",
	)
}
