package bls12

// #include "relic_core.h"
// #include "relic_fp.h"
// #include "relic_epx.h"
// void _ep2_add(ep2_t r, const ep2_t p, const ep2_t q) { ep2_add(r, p, q); }
// void _ep2_neg(ep2_t r, const ep2_t p) { ep2_neg(r, p); }
// void _ep2_mul(ep2_t r, const ep2_t p, const bn_t k) { ep2_mul(r, p, k); }
// void _fp_neg(fp_t r, const fp_t p);
// void _fp2_neg(fp_t r, const fp_t p);
// void ep2_mul_cof_b12(ep2_t r, ep2_t p); // unexported, don't @ me
// void ep2_scale_by_cofactor(ep2_t p);
// void ep2_read_x(ep2_t a, uint8_t* bin, int len) {
//     a->norm = 1;
//     fp_set_dig(a->z[0], 1);
//     fp_zero(a->z[1]);
//     fp2_read_bin(a->x, bin, len);
//     fp2_zero(a->y);
// }
import "C"
import "errors"

type G2 struct {
	st C.ep2_st
}

// p = G2(inf)
func (p *G2) SetZero() *G2 {
	C.ep2_set_infty(&p.st)
	return p
}

// p = G2(G)
func (p *G2) SetOne() *G2 {
	C.ep2_curve_get_gen(&p.st)
	return p
}

// p = s * G2(p)
func (p *G2) ScalarMult(s *Fr) *G2 {
	C._ep2_mul(&p.st, &p.st, &s.st)
	return p
}

// p = p + q
func (p *G2) Add(q *G2) *G2 {
	C._ep2_add(&p.st, &p.st, &q.st)
	return p
}

// p == q
func (p *G2) Equal(q *G2) bool {
	return C.ep2_cmp(&p.st, &q.st) == C.CMP_EQ
}

// p == G2(inf)
func (p *G2) IsZero() bool {
	return C.ep2_is_infty(&p.st) == 1
}

const (
	Fq2ElementSize     = 96
	G2CompressedSize   = Fq2ElementSize
	G2UncompressedSize = 2 * Fq2ElementSize
)

// Unmarshal a point on G2. It consumes either G2CompressedSize or
// G2UncompressedSize, depending on how the point was marshalled.
func (p *G2) Unmarshal(in []byte) ([]byte, error) {
	if len(in) < G2CompressedSize {
		return nil, errors.New("wrong encoded point size")
	}
	compressed := in[0]&serializationCompressed != 0
	inlen := G2UncompressedSize
	if compressed {
		inlen = G2CompressedSize
	}
	if !compressed && len(in) < G2UncompressedSize {
		return nil, errors.New("insufficient data to decode point")
	}
	var bin [G2UncompressedSize+1]byte

	// swap c0 and c1
	copy(bin[:], in[Fq2ElementSize/2:Fq2ElementSize])
	copy(bin[Fq2ElementSize/2:], in[:Fq2ElementSize/2])
	bin[Fq2ElementSize/2] &= serializationMask

	// Big Y, but we're not compressed, or infinity is serialized
	if ((in[0]&serializationBigY != 0) == !compressed || (in[0]&serializationInfinity != 0))  {
		return nil, errors.New("high Y bit improperly set")
	}

	if in[0]&serializationInfinity != 0 {
		// Check that rest is zero
		for _, v := range bin[1:inlen+1] {
			if v != 0 {
				return nil, errors.New("invalid infinity encoding")
			}
		}

		C.ep2_set_infty(&p.st)
		return in[inlen:], nil
	}

	if compressed {
		bin[0] = 2
		C.ep2_read_x(&p.st, (*C.uint8_t)(&bin[0]), G2CompressedSize+1)
		if C.ep2_upk(&p.st, &p.st) == 0 {
			return nil, errors.New("no square root found")
		}

		var yneg C.fp_st
		C._fp_neg(&yneg[0], &p.st.y[1][0])
		// yneg > y, compares only c1
		if (C.fp_cmp(&yneg[0], &p.st.y[1][0]) == C.CMP_GT) == (in[0]&serializationBigY != 0) {
			p.st.y[1] = yneg
			// negate c0 too
			C._fp_neg(&p.st.y[0][0], &p.st.y[0][0])
		}
		return in[G2CompressedSize:], nil
	}

	C.ep2_read_bin(&p.st, (*C.uint8_t)(&bin[0]), G2UncompressedSize+1)
	return in[G2UncompressedSize:], nil
}

// Marshal the point, compressed to X and sign.
func (p *G2) Marshal() (res []byte) {
	var bin [G2CompressedSize+1]byte
	res = bin[1:]
	if C.ep2_is_infty(&p.st) == 1 {
		res[0] |= serializationInfinity | serializationCompressed
		return
	}
	C.ep2_norm(&p.st, &p.st)
	C.ep2_write_bin((*C.uint8_t)(&bin[0]), G2CompressedSize+1, &p.st, 1)

	var bin2 [G2CompressedSize+1]byte
	copy(bin2[1:], res[Fq2ElementSize/2:Fq2ElementSize])
	copy(bin2[1+Fq2ElementSize/2:], res[:Fq2ElementSize/2])
	res = bin2[1:]
	res[0] |= serializationCompressed

	var yneg C.fp_st
	C._fp_neg(&yneg[0], &p.st.y[1][0])
	if C.fp_cmp(&yneg[0], &p.st.y[1][0]) == C.CMP_GT {
		res[0] |= serializationBigY
	}
	return
}

// Marshal the point, as uncompressed XY.
func (p *G2) MarshalUncompressed() (res []byte) {
	var bin [G2UncompressedSize+1]byte
	res = bin[1:]

	if C.ep2_is_infty(&p.st) == 1 {
		res[0] |= serializationInfinity
		return
	}
	C.ep2_write_bin((*C.uint8_t)(&bin[0]), G2UncompressedSize+1, &p.st, 0)
	var bin2 [G2CompressedSize+1]byte
	copy(bin2[1:], res[Fq2ElementSize/2:Fq2ElementSize])
	copy(bin2[1+Fq2ElementSize/2:], res[:Fq2ElementSize/2])
	res = bin2[1:]
	return
}

