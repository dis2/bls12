// +build cgo

package bls12

// #include "relic_core.h"
// #include "relic_fp.h"
// #include "relic_epx.h"
// void _ep2_add(ep2_t r, const ep2_t p, const ep2_t q) { ep2_add(r, p, q); }
// void _ep2_neg(ep2_t r, const ep2_t p) { ep2_neg(r, p); }
// void _ep2_mul(ep2_t r, const ep2_t p, const bn_t k) { ep2_mul(r, p, k); }
// void _fp_neg(fp_t r, const fp_t p);
// void _fp2_neg(fp_t r, const fp_t p);
// void ep2_mul_cof_b12(ep2_t r, ep2_t p);
// void ep2_scale_by_cofactor(ep2_t p);
// void ep2_read_x(ep2_t a, uint8_t* bin, int len) {
//     a->norm = 1;
//     fp_set_dig(a->z[0], 1);
//     fp_zero(a->z[1]);
//     fp2_read_bin(a->x, bin, len);
//     fp2_zero(a->y);
// }
// void ep2_scale_by_cofactor(ep2_t p) {
//     bn_t k;
//     bn_new(k);
//     bn_read_str(k, "5d543a95414e7f1091d50792876a202cd91de4547085abaa68a205b2e5a7ddfa628f1cb4d9e82ef21537e293a6691ae1616ec6e786f0c70cf1c38e31c7238e5", 127, 16);
//     ep2_mul_basic(p, p, k);
//     bn_free(k);
// }
import "C"
import "unsafe"

func (p *G2) l() *C.ep2_st {
	return (*C.ep2_st)(unsafe.Pointer(&p.X))
}

// p = G2(inf)
func (p *G2) SetZero() G {
	C.ep2_set_infty(p.l())
	return p
}

// Create new element set to infinity.
func G2Zero() (res *G2) {
	res = new(G2)
	res.SetZero()
	return
}

// Create new element set to generator.
func G2One() (res *G2) {
	res = new(G2)
	res.SetOne()
	return
}

// p = G2(G)
func (p *G2) SetOne() G {
	C.ep2_curve_get_gen(p.l())
	return p
}

// p = G2_h * G2(p)
func (p *G2) ScaleByCofactor() G {
	C.ep2_scale_by_cofactor(p.l())
	return p
}

// p = s * G2(p)
func (p *G2) ScalarMult(s *Scalar) G {
	C._ep2_mul(p.l(), p.l(), s)
	return p
}

// p = s * G2(G)
func (p *G2) ScalarBaseMult(s *Scalar) G {
	C._ep2_mul(p.l(), G2One().l(), s)
	return p
}

// p = p + q
func (p *G2) Add(q G) G {
	C._ep2_add(p.l(), p.l(), q.(*G2).l())
	return p
}

// Normalize the point into affine coordinates.
func (p *G2) Normalize() G {
	C.ep2_norm(p.l(), p.l())
	return p
}

// p == q
func (p *G2) Equal(q G) bool {
	return C.ep2_cmp(p.l(), q.(*G2).l()) == C.CMP_EQ
}

// p == G2(inf)
func (p *G2) IsZero() bool {
	return C.ep2_is_infty(p.l()) == 1
}
