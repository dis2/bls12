package bls12

// #include "relic_core.h"
// #include "relic_fp.h"
// #include "relic_ep.h"
// void _ep_add(ep_t r, const ep_t p, const ep_t q) { ep_add(r, p, q); }
// void _ep_neg(ep_t r, const ep_t p) { ep_neg(r, p); }
// void _ep_mul(ep_t r, const ep_t p, const bn_st *k) { ep_mul(r, p, k); }
// void _fp_neg(fp_t r, const fp_t p);
import "C"
import "errors"

// Point in G1 backed by a relic p.st.
type G1 struct {
	st C.ep_st
}

// p = G1(inf)
func (p *G1) SetZero() *G1 {
	C.ep_set_infty(&p.st)
	return p
}

// p = G1(G)
func (p *G1) SetOne() *G1 {
	C.ep_curve_get_gen(&p.st)
	return p
}

// p = s * G1(p)
func (p *G1) ScalarMult(s *Scalar) *G1 {
	C._ep_mul(&p.st, &p.st, &s.st)
	return p
}

// p = s * G1(G)
func (p *G1) ScalarBaseMult(s *Scalar) *G1 {
	C.ep_mul_gen(&p.st, &s.st)
	return p
}

// p = p + q
func (p *G1) Add(q *G1) *G1 {
	C._ep_add(&p.st, &p.st, &q.st)
	return p
}

// Check if points are the same. This is needed when the points are not
// in normalized form - there's an algebraic trick in relic to do the comparison
// faster than normalizing first. If you're sure the points are normalized, it's
// possible to compare directly with ==.
func (p *G1) Equal(q *G1) bool {
	return C.ep_cmp(&p.st, &q.st) == C.CMP_EQ
}

const (
	G1Size             = 48
	G1UncompressedSize = 2 * G1Size

	// https://github.com/ebfull/pairing/tree/master/src/bls12_381#serialization
	serializationMask       = (1 << 5) - 1
	serializationCompressed = 1 << 7
	serializationInfinity   = 1 << 6
	serializationBigY       = 1 << 5
)

// Unmarshal a point on G1. It consumes either G1Size or
// G1UncompressedSize, depending on how the point was marshalled.
func (p *G1) Unmarshal(in []byte) ([]byte, error) {
	if len(in) < G1Size {
		return nil, errors.New("wrong encoded point size")
	}
	compressed := in[0]&serializationCompressed != 0
	inlen := G1UncompressedSize
	if compressed {
		inlen = G1Size
	}
	if !compressed && len(in) < G1UncompressedSize {
		return nil, errors.New("insufficient data to decode point")
	}
	var bin [G1UncompressedSize + 1]byte
	copy(bin[1:], in[:inlen])
	bin[1] &= serializationMask

	// Big Y, but we're not compressed, or infinity is serialized
	if (in[0]&serializationBigY != 0) == !compressed || (in[0]&serializationInfinity != 0) {
		return nil, errors.New("high Y bit improperly set")
	}

	if in[0]&serializationInfinity != 0 {
		// Check that rest is zero
		for _, v := range bin[1 : inlen+1] {
			if v != 0 {
				return nil, errors.New("invalid infinity encoding")
			}
		}

		C.ep_set_infty(&p.st)
		return in[inlen:], nil
	}

	if compressed {
		bin[0] = 2
		C.ep_read_bin(&p.st, (*C.uint8_t)(&bin[0]), G1Size+1)

		var yneg C.fp_st
		C._fp_neg(&yneg[0], &p.st.y[0])
		// yneg > y?
		if (C.fp_cmp(&yneg[0], &p.st.y[0]) == C.CMP_GT) == (in[0]&serializationBigY != 0) {
			p.st.y = yneg
		}
		return in[G1Size:], nil
	}

	bin[0] = 4
	C.ep_read_bin(&p.st, (*C.uint8_t)(&bin[0]), G1UncompressedSize+1)
	return in[G1UncompressedSize:], nil
}

// Marshal the point, compressed to X and sign.
func (p *G1) Marshal() (res []byte) {
	var bin [G1Size + 1]byte
	res = bin[1:]
	res[0] |= serializationCompressed

	if C.ep_is_infty(&p.st) == 1 {
		res[0] |= serializationInfinity
		return
	}
	C.ep_norm(&p.st, &p.st)
	C.ep_write_bin((*C.uint8_t)(&bin[0]), G1Size+1, &p.st, 1)

	var yneg C.fp_st
	C._fp_neg(&yneg[0], &p.st.y[0])
	if C.fp_cmp(&yneg[0], &p.st.y[0]) == C.CMP_GT {
		res[0] |= serializationBigY
	}
	return
}

// Marshal the point, as uncompressed XY.
func (p *G1) MarshalUncompressed() (res []byte) {
	var bin [G1UncompressedSize + 1]byte
	res = bin[1:]

	if C.ep_is_infty(&p.st) == 1 {
		res[0] |= serializationInfinity
		return
	}
	C.ep_write_bin((*C.uint8_t)(&bin[0]), G1UncompressedSize+1, &p.st, 0)
	return
}
