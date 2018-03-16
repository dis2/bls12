// +build cgo

package bls12

// #include "relic_core.h"
// #include "relic_fp.h"
// void _fp2_neg(fp_t r,fp_t p) { fp2_neg(r, p); }
// void _fp2_mul(fp_t c,fp_t a,fp_t b) { fp2_mul(c,a,b); }
// void _fp2_add(fp_t c,fp_t a,fp_t b) { fp2_add(c,a,b); }
// void _fp2_sub(fp_t c,fp_t a,fp_t b) { fp2_sub(c,a,b); }
// void _fp2_sqr(fp_t c,fp_t a) { fp2_sqr(c,a); }
// void _fp2_inv(fp_t c,fp_t a) { fp2_inv(c,a); }
// int _fp2_srt(fp_t c,fp_t a) { return fp2_srt(c,a); }
// void _fp2_exp(fp_t c,fp_t a,bn_t b) { fp2_exp(c,a,b); }
import "C"

func (e *Fq2) l() *Limb {
	return e.C[0].l()
}

func (e *Fq2) le(x Field) *Limb {
	v, ok := x.(*Fq2)
	if !ok && x != nil {
		panic("invalid field type passed")
	}
	if v == nil {
		v = e
	}
	return v.l()
}

// e = a^n
func (e *Fq2) Exp(a Field, n *Scalar) Field {
	if e == nil {
		panic("shouldnt happen")
	}
	C._fp2_exp(e.l(), e.le(a), n)
	return e
}

// e = x^2
func (e *Fq2) Square(x Field) Field {
	C._fp2_sqr(e.l(), e.le(x))
	return e
}

// e = x^-1
func (e *Fq2) Sqrt(x Field) bool {
	return C._fp2_srt(e.l(), e.le(x)) == C.int(1)
}

// e = a + b
func (e *Fq2) Add(a, b Field) Field {
	C._fp2_add(e.l(), e.le(a), b.(*Fq2).l())
	return e
}

// e = a - b
func (e *Fq2) Sub(a, b Field) Field {
	C._fp2_sub(e.l(), e.le(a), b.(*Fq2).l())
	return e
}

// e = 1/x
func (e *Fq2) Inverse(x Field) Field {
	C._fp2_inv(e.l(), e.le(x))
	return e
}

// e = -x
func (e *Fq2) Neg(x Field) Field {
	C._fp2_neg(e.l(), e.le(x))
	return e
}

// e = a * b
func (e *Fq2) Mul(a, b Field) Field {
	C._fp2_mul(e.l(), e.le(a), b.(*Fq2).l())
	return e
}
