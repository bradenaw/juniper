//go:build go1.18

package tseq

import (
	"fmt"
	"testing"

	"github.com/bradenaw/juniper/internal/require2"
)

func TestTSeqBasic(t *testing.T) {
	runs := [][]bool{}
	Run(t, func(tseq *TSeq) {
		runs = append(runs, []bool{tseq.FlipCoin(), tseq.FlipCoin(), tseq.FlipCoin()})
	})

	require2.DeepEqual(t, runs, [][]bool{
		{false, false, false},
		{false, false, true},
		{false, true, false},
		{false, true, true},
		{true, false, false},
		{true, false, true},
		{true, true, false},
		{true, true, true},
	})
}

func TestTSeqDependent(t *testing.T) {
	runs := [][]bool{}
	Run(t, func(tseq *TSeq) {
		if tseq.FlipCoin() {
			runs = append(runs, []bool{true, tseq.FlipCoin()})
		} else {
			runs = append(runs, []bool{false})
		}
	})

	require2.DeepEqual(t, runs, [][]bool{
		{false},
		{true, false},
		{true, true},
	})
}

func TestTSeqChoose(t *testing.T) {
	for i := 1; i < 50; i++ {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			expected := []int{}
			for j := 0; j < i; j++ {
				expected = append(expected, j)
			}
			runs := []int{}
			Run(t, func(tseq *TSeq) {
				runs = append(runs, tseq.Choose(i))
			})

			require2.DeepEqual(t, runs, expected)
		})
	}
}
