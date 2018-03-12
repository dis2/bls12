// Package bls12 implements bilinear pairing curve BLS12-381
package bls12 // import "github.com/dis2/bls12"

// #include "relic_core.h"
// #include "relic_ep.h"
// void _fp_neg(fp_t r, const fp_t p) { fp_neg(r, p); }
import "C"
import "bytes"

func init() {
	C.core_init()
	C.ep_param_set_any_pairf()
}

// Check if encoding of negated coordinate is lexicographically bigger.
// Also returns the given negate.
func negativeIsBigger(neg, a *C.dig_t) bool {
	C._fp_neg(neg, a)

	var abuf, bbuf [48]byte
	C.fp_write_bin((*C.uint8_t)(&abuf[0]), 48, a)
	C.fp_write_bin((*C.uint8_t)(&bbuf[0]), 48, neg)
	return bytes.Compare(abuf[:], bbuf[:]) > 0
}
