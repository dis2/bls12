package bls12

// #include "relic_pp.h"
// #include "relic_ep.h"
// #include "relic_fpx.h"
// void _fp12_mul(fp12_t c, fp12_t a, fp12_t b) { fp12_mul(c,a,b); }
import "C"

type FP12 struct {
	// Go copes with typedef arrays very poorly, and resolves type
	// off-by-one level deeper than it should.
	fp [2]C.fp6_t
}

func (f *FP12) Pair(ep *EP, ep2 *EP2) (*FP12) {
	C.pp_map_oatep_k12(&f.fp[0], &ep.st, &ep2.t)
	return f
}

func (f *FP12) Mul(a,b *FP12) (*FP12) {
	C._fp12_mul(&f.fp[0], &a.fp[0], &b.fp[0])
	return f
}
