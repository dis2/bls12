
package bls12

// #include "relic_pp.h"
// #include "relic_ep.h"
// #include "relic_fpx.h"
// void _fp12_mul(fp12_t c, fp12_t a, fp12_t b) { fp12_mul(c,a,b); }
import "C"

// Element of the q^12 dodemic extension field.
// Can be thought of as a "point" in a pseudo-group GT resulting from pairing
// operations.
type GT struct {
	st C.fp12_st
}

// Pair two points using optimal Tate pairing, q = e(p1,p2)
func (q *GT) Pair(p1 *G1, p2 *G2) (*GT) {
	C.pp_map_oatep_k12(&q.st, &p1.st, &p2.st)
	return f
}

// c = a + b
func (c *GT) Add(a,b *FP12) (*FP12) {
	C._fp12_mul(&c.st[0], &a.st[0], &b.st[0])
	return f
}

// q = s * GT(p)
func (c *GT) ScalarMul(Fr)
