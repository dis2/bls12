// Package bls12 implements bilinear pairing curve BLS12-381
package bls12 // import "github.com/dis2/bls12"

// #include "relic_core.h"
// #include "relic_ep.h"
// void _fp_neg(fp_t r, const fp_t p) { fp_neg(r, p); }
// void _bn_new(bn_t bn) { bn_new(bn); }
import "C"
import "bytes"
import "math/big"
import "encoding/hex"

// Check if encoding of negated coordinate is lexicographically bigger.
// Also returns the given negate.
func negativeIsBigger(neg, a *C.dig_t) bool {
	C._fp_neg(neg, a)

	var abuf, bbuf [48]byte
	C.fp_write_bin((*C.uint8_t)(&abuf[0]), 48, a)
	C.fp_write_bin((*C.uint8_t)(&bbuf[0]), 48, neg)
	return bytes.Compare(abuf[:], bbuf[:]) > 0
}

func hexConst(s string) (ret *big.Int) {
	ret, _ = new(big.Int).SetString(s, 16)
	return
}

var init_done = false

func init_pending() {
	if init_done {
		return
	}
	init_done = true
	C.core_init()
	C.ep_param_set_any_pairf()
}

func QConst(s string) (f Fq) {
	if len(s)%2 != 0 {
		panic("bad const padding for "+s)
	}
	init_pending()
	var buf [48]byte
	pad := 48 - len(s)/2
	_, err := hex.Decode(buf[pad:], []byte(s))
	if err != nil || f.Unmarshal(buf[:]) == nil {
		panic("invalid const " + s)
	}
	return
}


