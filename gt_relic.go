// +build cgo

package bls12

// #include "relic_core.h"
// #include "relic_pp.h"
// #include "relic_ep.h"
// #include "relic_fpx.h"
// void _fp12_mul(fp12_t c, fp12_t a, fp12_t b) { fp12_mul(c,a,b); }
import "C"
import "unsafe"

func (p *GT) l() *C.fp6_t {
	return (*C.fp6_t)(unsafe.Pointer(&p[0]))
}

// Find optimal ate pairing for p1 and p2, q = e(p1,p2)
func (q *GT) Pair(p1 *G1, p2 *G2) *GT {
	C.pp_map_oatep_k12(q.l(), p1.l(), p2.l())
	return q
}

// p = p + q
func (p *GT) Add(q *GT) *GT {
	C._fp12_mul(p.l(), p.l(), q.l())
	return p
}

// q = s * GT(p)
func (p *GT) ScalarMult(s *Scalar) (q *GT) {
	q = &GT{}
	C.fp12_exp(q.l(), p.l(), s)
	return
}

// q = -p
func (p *GT) Neg() (q *GT) {
	q = &GT{}
	C.fp12_inv(q.l(), p.l())
	return
}

// p == q
func (p *GT) Equal(q *GT) bool {
	return C.fp12_cmp(p.l(), q.l()) == C.CMP_EQ
}

// Marshal GT into a byte buffer.
func (p *GT) Marshal() []byte {
	var bin [GTSize]byte
	C.fp12_write_bin((*C.uint8_t)(&bin[0]), GTSize, p.l(), 1)
	return bin[:]
}


// Unmarshal GT from a byte buffer.
func (p *GT) Unmarshal(in []byte) ([]byte) {
	if len(in) < GTSize {
		return nil
	}
	C.fp12_read_bin(p.l(), (*C.uint8_t)(&in[0]), GTSize)
	return in[GTSize:]
}


