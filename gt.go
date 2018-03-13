package bls12

// #include "relic_core.h"
// #include "relic_pp.h"
// #include "relic_ep.h"
// #include "relic_fpx.h"
// void _fp12_mul(fp12_t c, fp12_t a, fp12_t b) { fp12_mul(c,a,b); }
import "C"

const (
	GTSize = 384
)

// Element of the q^12 extension field.
// Can be thought of as a pseudo-point in a group GT resulting from pairing
// operations.
//
// The group law is shifted from the underlying Fq12 field. Ie multiplication
// is field exponentation etc.
type GT struct {
	st [2]C.fp6_t // Workaround for Go type resolution bugs.
}

// Find optimal ate pairing for p1 and p2, q = e(p1,p2)
func (q *GT) Pair(p1 *G1, p2 *G2) *GT {
	C.pp_map_oatep_k12(&q.st[0], p1, p2)
	return q
}

// p = p + q
func (p *GT) Add(q *GT) *GT {
	C._fp12_mul(&p.st[0], &p.st[0], &q.st[0])
	return p
}

// q = s * GT(p)
func (p *GT) ScalarMult(s *Scalar) (q *GT) {
	q = &GT{}
	C.fp12_exp(&q.st[0], &p.st[0], s)
	return
}

// q = -p
func (p *GT) Neg() (q *GT) {
	q = &GT{}
	C.fp12_inv(&q.st[0], &p.st[0])
	return
}

// p == q
func (p *GT) Equal(q *GT) bool {
	return C.fp12_cmp(&p.st[0], &q.st[0]) == C.CMP_EQ
}

// Marshal GT into a byte buffer.
func (p *GT) Marshal() []byte {
	var bin [GTSize]byte
	C.fp12_write_bin((*C.uint8_t)(&bin[0]), GTSize, &p.st[0], 1)
	return bin[:]
}


// Unmarshal GT from a byte buffer.
func (p *GT) Unmarshal(in []byte) ([]byte) {
	if len(in) < GTSize {
		return nil
	}
	C.fp12_read_bin(&p.st[0], (*C.uint8_t)(&in[0]), GTSize)
	return in[GTSize:]
}
