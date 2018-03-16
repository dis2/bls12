// +build cgo

package bls12

// #include "relic_core.h"
// #include "relic_fp.h"
// void _fp_neg(fp_t r, const fp_t p);
// void _fp_mul(fp_t c,fp_t a,fp_t b) { fp_mul(c,a,b); }
// void _fp_add(fp_t c,fp_t a,fp_t b) { fp_add(c,a,b); }
// void _fp_sub(fp_t c,fp_t a,fp_t b) { fp_sub(c,a,b); }
// void _fp_sqr(fp_t c,fp_t a) { fp_sqr(c,a); }
// void _fp_inv(fp_t c,fp_t a) { fp_inv(c,a); }
// void _fp_exp(fp_t c,fp_t a,bn_t b) { fp_exp(c,a,b); }
import "C"
import "bytes"

type Limbs [NLimbs]Limb

type Fq struct {
	Limbs
}


func (e *Fq) l() *Limb {
	return &e.Limbs[0]
}

func (e *Fq) le(x Field) *Limb {
	v, ok := x.(*Fq)
	if !ok && x != nil {
		panic("invalid field type passed")
	}
	if v == nil {
		v = e
	}
	return v.l()
}

// e = a^n
func (e *Fq) Exp(a Field, n *Scalar) Field {
	C._fp_exp(e.l(), e.le(a), n)
	return e
}

// e = x^2
func (e *Fq) Square(x Field) Field {
	C._fp_sqr(e.l(), e.le(x))
	return e
}

// e = a + b
func (e *Fq) Add(a, b Field) Field {
	C._fp_add(e.l(), e.le(a), b.(*Fq).l())
	return e
}

// e = a - b
func (e *Fq) Sub(a, b Field) Field {
	C._fp_sub(e.l(), e.le(a), b.(*Fq).l())
	return e
}

// e = 1/x
func (e *Fq) Inverse(x Field) Field {
	C._fp_inv(e.l(), e.le(x))
	return e
}

// e = -x
func (e *Fq) Neg(x Field) Field {
	C._fp_neg(e.l(), e.le(x))
	return e
}

// e == x
func (e *Fq) Equal(x Field) bool {
	return C.fp_cmp(e.l(), x.(*Fq).l()) == C.CMP_EQ
}

// e > x
func (e *Fq) GreaterThan(x Field) bool {
	var buf1, buf2 [48]byte
	C.fp_write_bin((*C.uint8_t)(&buf1[0]), 48, e.l())
	C.fp_write_bin((*C.uint8_t)(&buf2[0]), 48, x.(*Fq).l())
	return bytes.Compare(buf1[:], buf2[:]) == 1
}

// e = a * b
func (e *Fq) Mul(a, b Field) Field {
	C._fp_mul(e.l(), e.le(a), b.(*Fq).l())
	return e
}

// e = 64 bit immediate n
func (e *Fq) SetInt64(n int64) Field {
	C.fp_set_dig(e.l(), Limb(n))
	return e
}

func (e *Fq) Unmarshal(b []byte) []byte {
	if len(b) < 48 {
		return nil
	}
	C.fp_read_bin(e.l(), (*C.uint8_t)(&b[0]), 48)
	return b[48:]
}

func (e *Fq) Marshal() []byte {
	var buf [48]byte
	C.fp_write_bin((*C.uint8_t)(&buf[0]), 48, e.l())
	return buf[:]
}
