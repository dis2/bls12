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
import "fmt"
import "math/big"

type Fq = C.fp_st

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

// e = a + b
func (e *Fq) Add(a,b *Fq) *Fq {
	C._fp_add(&e[0], &a[0], &b[0])
	return e
}

// e = a - b
func (e *Fq) Sub(a,b *Fq) *Fq {
	C._fp_sub(&e[0], &a[0], &b[0])
	return e
}


// e = x^-1
func (e *Fq) Inverse(x *Fq) *Fq {
	C._fp_inv(&e[0], &x[0])
	return e
}

// e = -x
func (e *Fq) Neg(x *Fq) *Fq {
	C._fp_neg(&e[0], &x[0])
	return e
}


func (e *Fq) Copy() Fq {
	return *e
}

// e == x
func (e *Fq) Equal(x *Fq) bool {
	return C.fp_cmp(&e[0], &x[0]) == C.CMP_EQ
}

// e > x
func (e *Fq) GreaterThan(x *Fq) bool {
	return C.fp_cmp(&e[0], &x[0]) == C.CMP_GT
}


func (e *Fq) IsZero() bool {
	return e.Equal(&Zero)
}

// e = a * b
func (e *Fq) Mul(a, b *Fq) *Fq {
	C._fp_mul(&e[0], &a[0], &b[0])
	return e
}

// e = x^3
func (e *Fq) Cube(x *Fq) *Fq {
	e.Square(x)
	e.Mul(e, x)
	return e
}

// e = 64 bit immediate n
func (e *Fq) SetInt64(n int64) *Fq {
	C.fp_set_dig(&e[0], C.dig_t(n))
	return e
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

func (e *Fq) String() string {
	return fmt.Sprintf("Fq(%d)", e.ToInt())
}
