// +build cgo

package bls12

// #include "relic_core.h"
// #include "relic_fp.h"
// #include "relic_ep.h"
// void _ep_add(ep_t r, const ep_t p, const ep_t q) { ep_add(r, p, q); }
// void _ep_neg(ep_t r, const ep_t p) { ep_neg(r, p); }
// void _ep_mul(ep_t r, const ep_t p, const bn_st *k) { ep_mul(r, p, k); }
// void _fp_neg(fp_t r, const fp_t p);
// void ep_mul_cof_b12(ep_t r, ep_t p);
import "C"
import "unsafe"

func (p *G1) l() *C.ep_st {
	return (*C.ep_st)(unsafe.Pointer(&p.X))
}

// Normalize the point into affine coordinates.
func (p *G1) Normalize() G {
	C.ep_norm(p.l(), p.l())
	return p
}

// p = G1(inf)
func (p *G1) SetZero() G {
	C.ep_set_infty(p.l())
	return p
}

// p = G1(G)
func (p *G1) SetOne() G {
	C.ep_curve_get_gen(p.l())
	return p
}

// Create new element set to infinity.
func G1Zero() (res *G1) {
	res = new(G1)
	res.SetZero()
	return
}

// Create new element set to generator.
func G1One() (res *G1) {
	res = new(G1)
	res.SetOne()
	return
}

// p = s * G1(p)
func (p *G1) ScalarMult(s *Scalar) G {
	C._ep_mul(p.l(), p.l(), s)
	return p
}

// p = s * G1(G)
func (p *G1) ScalarBaseMult(s *Scalar) G {
	C.ep_mul_gen(p.l(), s)
	return p
}

// p = p + q
func (p *G1) Add(q G) G {
	C._ep_add(p.l(), p.l(), q.(*G1).l())
	return p
}

func (p *G1) HashToPointRelic(msg []byte) G {
	C.ep_map(p.l(), (*C.uint8_t)(&msg[0]), C.int(len(msg)))
	return p
}

// Check if points are the same. This is needed when the points are not
// in normalized form - there's an algebraic trick in relic to do the comparison
// faster than normalizing first. If you're sure the points are normalized, it's
// possible to compare directly with ==.
func (p *G1) Equal(q G) bool {
	return C.ep_cmp(p.l(), q.(*G1).l()) == C.CMP_EQ
}

// p == G1(inf)
func (p *G1) IsZero() bool {
	return C.ep_is_infty(p.l()) == 1
}
