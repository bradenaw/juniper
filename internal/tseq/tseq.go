package tseq

import (
	"encoding/hex"
	"flag"
	"testing"
)

var seed string

func init() {
	flag.StringVar(&seed, "tseq.seed", seed, "")
}

func Run(t *testing.T, f func(tseq *TSeq)) {
	tseq := TSeq{}

	if seed == "" {
		defer func() {
			if t.Failed() {
				t.Logf("rerun with --tseq.seed=%s", scriptToString(tseq.script))
			}
		}()

		for {
			f(&tseq)
			if !tseq.next() {
				break
			}
		}
	} else {
		var err error
		tseq.script, err = scriptFromString(seed)
		if err != nil {
			t.Fatalf("invalid --tseq.seed %s: %s", seed, err)
		}
		f(&tseq)
	}
}

type TSeq struct {
	script []bool
	i      int
}

func scriptToString(script []bool) string {
	b := make([]byte, (len(script)+7)/8)
	for i := range script {
		if script[i] {
			idx := i / 8
			off := i % 8
			b[idx] |= 1 << off
		}
	}
	return "00" + hex.EncodeToString(b)
}

func scriptFromString(s string) ([]bool, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	b = b[1:]
	script := make([]bool, len(b)*8)
	for i := range script {
		idx := i / 8
		off := i % 8
		script[i] = ((b[idx] >> off) & 1) == 1
	}
	return script, nil
}

func (tseq *TSeq) FlipCoin() bool {
	if tseq.i == len(tseq.script) {
		tseq.script = append(tseq.script, false)
	}
	outcome := tseq.script[tseq.i]
	tseq.i++
	return outcome
}

func (tseq *TSeq) Choose(n int) int {
	if n == 0 {
		panic("can't choose between 0 choices")
	}
	i := 0
	j := n - 1
	for i < j {
		mid := (i + j) / 2
		if tseq.FlipCoin() {
			i = mid + 1
		} else {
			j = mid
		}
	}
	return i
}

func (tseq *TSeq) next() bool {
	for len(tseq.script) > 0 && tseq.script[len(tseq.script)-1] {
		tseq.script = tseq.script[:len(tseq.script)-1]
	}
	if len(tseq.script) == 0 {
		return false
	}
	tseq.script[len(tseq.script)-1] = true
	tseq.i = 0
	return true
}
