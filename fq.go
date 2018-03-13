package bls12

// #include "relic_core.h"
// #include "relic_fp.h"
// void _fp_mul(fp_t c,fp_t a,fp_t b) { fp_mul(c,a,b); }
// void _fp_sqr(fp_t c,fp_t a) { fp_sqr(c,a); }
// void _fp_exp(fp_t c,fp_t a,bn_t b) { fp_exp(c,a,b); }
import "C"
import "math/big"

type Fq = C.fp_st

// e = a ^ ((q-1)/2)A
func (e *Fq) Legendre(a *Fq) {
}

// e = a^n
func (e *Fq) Exp(a *Fq, n *Scalar) {
	if e == nil {
		panic("shouldnt happen")
	}
	C._fp_exp(&e[0], &a[0], n)
}

// e = x^2
func (e *Fq) Square(x *Fq) *Fq {
	C._fp_sqr(&e[0], &x[0])
	return e
}

func (e *Fq) Copy() Fq {
	return *e
}

// e == a
func (e *Fq) Equal(a *Fq) bool {
	return C.fp_cmp(&e[0], &a[0]) == C.CMP_EQ
}

// e = x * y
func (e *Fq) Mul(x, y *Fq) *Fq {
	C._fp_mul(&e[0], &x[0], &y[0])
	return e
}

// e = x^3
func (e *Fq) Cube(x *Fq) *Fq {
	e.Square(x)
	e.Mul(e, x)
	return e
}

// e = 64 bit immediate n
func (e *Fq) SetInt64(n int64) {
	C.fp_set_dig(&e[0], C.dig_t(n))
}

// e = a + 64 bit immediate n
func (e *Fq) AddInt64(a *Fq, n int64) {
	C.fp_add_dig(&e[0], &a[0], C.dig_t(n))
}

func pad(buf []byte) []byte {
	n := len(buf)
	if n > 48 {
		return buf
	}
	return append(make([]byte, 48-n), buf...)
}

func (e *Fq) Unmarshal(b []byte) []byte {
	if len(b) < 48 {
		return nil
	}
	C.fp_read_bin(&e[0], (*C.uint8_t)(&b[0]), 48)
	return b[48:]
}

func (e *Fq) Marshal() []byte {
	var buf [48]byte
	C.fp_write_bin((*C.uint8_t)(&buf[0]), 48, &e[0])
	return buf[:]
}


func (e *Fq) FromInt(b *big.Int) *Fq {
	if e.Unmarshal(pad(b.Bytes())) == nil {
		return nil
	}
	return e
}

func (e *Fq) ToInt() *big.Int {
	return new(big.Int).SetBytes(e.Marshal())
}


