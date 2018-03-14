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
import "fmt"

type G2 = C.ep2_st

func (p *G2) Copy() *G2 {
	c := *p
	return &c
}

// p = G2(inf)
func (p *G2) SetZero() *G2 {
	C.ep2_set_infty(p)
	return p
}

// p = G2(G)
func (p *G2) SetOne() *G2 {
	C.ep2_curve_get_gen(p)
	return p
}

// p = G_h * G(p)
func (p *G2) ScaleByCofactor() {
}

// p = s * G2(p)
func (p *G2) ScalarMult(s *Scalar) *G2 {
	C._ep2_mul(p, p, s)
	return p
}

// p = s * G2(G)
func (p *G2) ScalarBaseMult(s *Scalar) *G2 {
	C._ep2_mul(p, new(G2).SetOne(), s)
	return p
}

// p = p + q
func (p *G2) Add(q *G2) *G2 {
	C._ep2_add(p, p, q)
	return p
}

// p == q
func (p *G2) Equal(q *G2) bool {
	return C.ep2_cmp(p, q) == C.CMP_EQ
}

// p == G2(inf)
func (p *G2) IsZero() bool {
	return C.ep2_is_infty(p) == 1
}

// HashToPoint the buffer.
func (p *G2) HashToPoint(b []byte) *G2 {
	C.ep2_map(p, (*C.uint8_t)(&b[0]), C.int(len(b)))
	return p
}

const (
	G2Size             = 96
	G2UncompressedSize = 2 * G2Size
)

// Unmarshal a point on G2. It consumes either G2Size or
// G2UncompressedSize, depending on how the point was marshalled.
func (p *G2) Unmarshal(in []byte) []byte {
	if len(in) < G2Size {
		return nil
	}
	compressed := in[0]&serializationCompressed != 0
	inlen := G2UncompressedSize
	if compressed {
		inlen = G2Size
	}
	if !compressed && len(in) < G2UncompressedSize {
		return nil
	}
	var bin [G2UncompressedSize + 1]byte

	// Big Y set, but we're not compressed, or infinity is serialized
	if (in[0]&serializationBigY != 0) && (!compressed || (in[0]&serializationInfinity != 0)) {
		return nil
	}

	if in[0]&serializationInfinity != 0 {
		// Check that rest is zero
		for _, v := range bin[1 : inlen+1] {
			if v != 0 {
				return nil
			}
		}

		C.ep2_set_infty(p)
		return in[inlen:]
	}

	// swap c0 and c1
	bin[0] = 4
	copy(bin[1:], in[G2Size/2:G2Size])
	copy(bin[1+G2Size/2:], in[:G2Size/2])
	bin[1+G2Size/2] &= serializationMask

	if compressed {
		C.ep2_read_x(p, (*C.uint8_t)(&bin[1]), G2Size)
		if C.ep2_upk(p, p) == 0 {
			return nil
		}

		var yneg C.fp_st
		if negativeIsBigger(&yneg[0], &p.y[1][0]) != (in[0]&serializationBigY != 0) {
			p.y[1] = yneg
			// negate c0 too
			C._fp_neg(&p.y[0][0], &p.y[0][0])
		}

		return in[G2Size:]
	}
	copy(bin[1+G2Size:], in[G2Size+G2Size/2:])
	copy(bin[1+G2Size+G2Size/2:], in[G2Size:])
	C.ep2_read_bin(p, (*C.uint8_t)(&bin[0]), G2UncompressedSize+1)
	return in[G2UncompressedSize:]
}

// Marshal the point, compressed to X and sign.
func (p *G2) Marshal() (res []byte) {
	var bin [G2Size + 1]byte
	res = bin[1:]
	if C.ep2_is_infty(p) == 1 {
		res[0] |= serializationInfinity | serializationCompressed
		return
	}
	C.ep2_norm(p, p)
	C.ep2_write_bin((*C.uint8_t)(&bin[0]), G2Size+1, p, 1)

	var bin2 [G2Size + 1]byte
	copy(bin2[1:], res[G2Size/2:G2Size])
	copy(bin2[1+G2Size/2:], res[:G2Size/2])
	res = bin2[1:]
	res[0] |= serializationCompressed
	var yneg C.fp_st
	if negativeIsBigger(&yneg[0], &p.y[1][0]) {
		res[0] |= serializationBigY
	}
	return
}

func (p *G2) String() string {
	return fmt.Sprintf("bls12.G2(%x)", p.Marshal())
}

// Marshal the point, as uncompressed XY.
func (p *G2) MarshalUncompressed() (res []byte) {
	var bin [G2UncompressedSize + 1]byte
	res = bin[1:]

	if C.ep2_is_infty(p) == 1 {
		res[0] |= serializationInfinity
		return
	}
	C.ep2_write_bin((*C.uint8_t)(&bin[0]), G2UncompressedSize+1, p, 0)
	var bin2 [G2UncompressedSize + 1]byte
	copy(bin2[1:], res[G2Size/2:G2Size])
	copy(bin2[1+G2Size/2:], res[:G2Size/2])
	copy(bin2[1+G2Size:], res[G2Size+G2Size/2:])
	copy(bin2[1+G2Size+G2Size/2:], res[G2Size:])
	return bin2[1:]
}
