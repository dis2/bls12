package bls12

// #include "relic_core.h"
// #include "relic_fp.h"
// #include "relic_ep.h"
// void _ep_add(ep_t r, const ep_t p, const ep_t q) { ep_add(r, p, q); }
// void _ep_neg(ep_t r, const ep_t p) { ep_neg(r, p); }
// void _ep_mul(ep_t r, const ep_t p, const bn_st *k) { ep_mul(r, p, k); }
// void _fp_neg(fp_t r, const fp_t p);
// void ep_mul_cof_b12(ep_t r, ep_t p);
import "C"
import "fmt"

// Point on G1, y^2 = x^3 + 4
type G1 = C.ep_st

// Computes x,y for nearest point in Fq, given the x.
// We land in unknown subgroup, caller is responsible for scaling by cofactor.
func g1MapXtoY(in *Fq) (x, y Fq) {
	x = *in
	for {
		var y2, ytest Fq
		// y2 = y^2 = x^3 + 4
		y2.Square(&x)
		y2.Mul(&y2,&x)
		y2.AddInt64(&y2, 4)

		// y = y2 ^ ((q+1)/4)
		y.Exp(&y2, &QPlus1Quarter)

		// if y^2 == y2
		if y2.Equal(ytest.Square(&y)) {
			return
		}

		x.AddInt64(&x,1)
	}
}

// Set affine coordinates X,Y with implicit Z=1
func (p *G1) SetXY(x, y *Fq) *G1 {
	p.x = *x
	p.y = *y

	// Implicitly normalized
	p.z.SetInt64(1)
	p.norm = 1
	return p
}

// Get pointers to raw coordinates of the element.
func (p *G1) GetXYZ() (x,y,z *Fq) {
	return &p.x, &p.y, &p.z
}

// Normalize XY to Z=1
func (p *G1) Normalize() {
	C.ep_norm(p, p)
}

// p = G1_h * G1(p)
func (p *G1) ScaleByCofactor() {
	p.ScalarMult(&G1_h)
}

// p = G1(inf)
func (p *G1) SetZero() *G1 {
	C.ep_set_infty(p)
	return p
}

// p = G1(G)
func (p *G1) SetOne() *G1 {
	C.ep_curve_get_gen(p)
	return p
}

// p = s * G1(p)
func (p *G1) ScalarMult(s *Scalar) *G1 {
	C._ep_mul(p, p, s)
	return p
}

// p = s * G1(G)
func (p *G1) ScalarBaseMult(s *Scalar) *G1 {
	C.ep_mul_gen(p, s)
	return p
}

// p = p + q
func (p *G1) Add(q *G1) *G1 {
	C._ep_add(p, p, q)
	return p
}

// Check if points are the same. This is needed when the points are not
// in normalized form - there's an algebraic trick in relic to do the comparison
// faster than normalizing first. If you're sure the points are normalized, it's
// possible to compare directly with ==.
func (p *G1) Equal(q *G1) bool {
	return C.ep_cmp(p, q) == C.CMP_EQ
}

// p == G1(inf)
func (p *G1) IsZero() bool {
	return C.ep_is_infty(p) == 1
}

// HashToPoint the buffer, using whatever relic does.
func (p *G1) HashToPointRelic(b []byte) *G1 {
	C.ep_map(p, (*C.uint8_t)(&b[0]), C.int(len(b)))
	return p
}

func (p *G1) HashToPoint(b []byte) *G1 {
	C.ep_map(p, (*C.uint8_t)(&b[0]), C.int(len(b)))
	return p
}

// Map arbitrary integer to a point, for use with custom hash function.
func (p *G1) MapIntToPoint(in *Fq) *G1 {
	x, y := g1MapXtoY(in)
	p.SetXY(&x,&y)
	p.ScaleByCofactor()
	return p
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
func (p *G1) Unmarshal(in []byte) []byte {
	if len(in) < G1Size {
		return nil
	}
	compressed := in[0]&serializationCompressed != 0
	inlen := G1UncompressedSize
	if compressed {
		inlen = G1Size
	}
	if !compressed && len(in) < G1UncompressedSize {
		return nil
	}
	var bin [G1UncompressedSize + 1]byte
	copy(bin[1:], in[:inlen])
	bin[1] &= serializationMask

	// Big Y, but we're not compressed, or infinity is serialized
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

		C.ep_set_infty(p)
		return in[inlen:]
	}

	if compressed {
		bin[0] = 2
		C.ep_read_bin(p, (*C.uint8_t)(&bin[0]), G1Size+1)
		var yneg C.fp_st

		if negativeIsBigger(&yneg[0], &p.y[0]) != (in[0]&serializationBigY != 0) {
			p.y = yneg
		}
		return in[G1Size:]
	}

	bin[0] = 4
	C.ep_read_bin(p, (*C.uint8_t)(&bin[0]), G1UncompressedSize+1)
	return in[G1UncompressedSize:]
}

// Marshal the point, compressed to X and sign.
func (p *G1) Marshal() (res []byte) {
	var bin [G1Size + 1]byte
	res = bin[1:]

	if C.ep_is_infty(p) == 1 {
		res[0] = serializationInfinity | serializationCompressed
		return
	}
	C.ep_norm(p, p)
	C.ep_write_bin((*C.uint8_t)(&bin[0]), G1Size+1, p, 1)
	res[0] |= serializationCompressed

	var yneg C.fp_st
	if negativeIsBigger(&yneg[0], &p.y[0]) {
		res[0] |= serializationBigY
	}
	return
}

func (p *G1) String() string {
	x,y,_ := p.GetXYZ()
	return fmt.Sprintf("bls12.G1(%d,%d)", x.ToInt(),y.ToInt())
}

// Marshal the point, as uncompressed XY.
func (p *G1) MarshalUncompressed() (res []byte) {
	var bin [G1UncompressedSize + 1]byte
	res = bin[1:]

	if C.ep_is_infty(p) == 1 {
		res[0] |= serializationInfinity
		return
	}
	C.ep_write_bin((*C.uint8_t)(&bin[0]), G1UncompressedSize+1, p, 0)
	return
}
