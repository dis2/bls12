package bls12

// #include "relic_pp.h"
// #include "relic_ep.h"
// #include "relic_fpx.h"
// void _fp12_mul(fp12_t c, fp12_t a, fp12_t b) { fp12_mul(c,a,b); }
import "C"
import "errors"

const (
	GTSize = 48
)

// Element of the q^12 extension field.
// Can be thought of as a "point" in a pseudo-group GT resulting from pairing
// operations.
//
// CAVEAT: The operators of the group are shifted from the field.
// Addition in group is multiplication in the field. Scalar multiplication in
// group, is exponentiation in the field. Negation is inverse, and so on...
type GT struct {
	st [2]C.fp6_t // Workaround for Go type resolution bugs.
}

// Pair two points using optimal Tate pairing, q = e(p1,p2)
func (q *GT) Pair(p1 *G1, p2 *G2) *GT {
	C.pp_map_oatep_k12(&q.st[0], &p1.st, &p2.st)
	return q
}

// c = a + b
func (c *GT) Add(a, b *GT) *GT {
	C._fp12_mul(&c.st[0], &a.st[0], &b.st[0])
	return c
}

// q = s * GT(p)
func (p *GT) ScalarMult(s *Scalar) (q *GT) {
	q = &GT{}
	C.fp12_exp(&q.st[0], &p.st[0], &s.st)
	return
}

// q = -p
func (p *GT) Neg() (q *GT) {
	q = &GT{}
	C.fp12_inv(&q.st[0], &p.st[0])
	return
}

// Marshal GT into a byte buffer.
func (p *GT) Marshal() []byte {
	var bin [GTSize]byte
	C.fp12_write_bin((*C.uint8_t)(&bin[0]), G1Size, &p.st[0], 1)
	return bin[:]
}

// Unmarshal GT from a byte buffer.
func (p *GT) Unmarshal(in []byte) ([]byte, error) {
	if len(in) < GTSize {
		return nil, errors.New("wrong encoded GT size")
	}
	C.fp12_read_bin(&p.st[0], (*C.uint8_t)(&in[0]), G1Size)
	return in[GTSize:], nil
}
