package bls12

// #include "relic_core.h"
// #include "relic_fp.h"
// #include "relic_ep.h"
// void _ep_add(ep_t r, const ep_t p, const ep_t q) { ep_add(r, p, q); }
// void _ep_neg(ep_t r, const ep_t p) { ep_neg(r, p); }
// void _ep_mul(ep_t r, const ep_t p, const bn_st *k) { ep_mul(r, p, k); }
// void _fp_neg(fp_t r, const fp_t p);
import "C"
import "fmt"
import "math/big"

// Point on G1, y^2 = x^3 + 4
type G1 struct {
	st C.ep_st
}

// Compute valid points on the curve (but unknown subgroup)
func hashToCurvePoint(x *big.Int) (*big.Int, *big.Int) {
	x.Mod(x, Q)

	four := big.NewInt(4)
	one := big.NewInt(1)
	for {
		xxx := new(big.Int).Mul(x, x)
		xxx.Mul(xxx, x)
		t := new(big.Int).Add(xxx, four)
		y := new(big.Int).ModSqrt(t, Q)
		if y != nil {
			return x, y
		}

		x.Add(x, one)
	}
}

func pad(buf []byte, to int) []byte {
	n := len(buf)
	if n > to {
		return buf
	}
	return append(make([]byte, to-n), buf...)
}

// Set raw affine coordinates X,Y
func (p *G1) SetXY(x, y *big.Int) *G1 {
	C.fp_read_bin(&p.st.x[0], (*C.uint8_t)(&pad(x.Bytes(),G1Size)[0]), G1Size)
	C.fp_read_bin(&p.st.y[0], (*C.uint8_t)(&pad(y.Bytes(),G1Size)[0]), G1Size)

	// Implicitly normalized
	C.fp_set_dig(&p.st.z[0], 1)
	p.st.norm = 1
	return p
}

// Get raw affine coordinates
func (p *G1) GetXY() (x,y *big.Int) {
	var t C.ep_st
	C.ep_norm(&t, &p.st)
	var bx, by [G1Size]byte
	C.fp_write_bin((*C.uint8_t)(&bx[0]), G1Size, &t.x[0]);
	C.fp_write_bin((*C.uint8_t)(&by[0]), G1Size, &t.y[0]);
	x = new(big.Int).SetBytes(bx[:])
	y = new(big.Int).SetBytes(by[:])
	return
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

// p == G1(inf)
func (p *G1) IsZero() bool {
	return C.ep_is_infty(&p.st) == 1
}

// HashToPoint the buffer.
func (p *G1) HashToPoint(b []byte) *G1 {
	C.ep_map(&p.st, (*C.uint8_t)(&b[0]), C.int(len(b)))
	return p
}

// Hash arbitrary integer to a point, use with a custom hash function.
func (p *G1) HashIntToPoint(x *big.Int) *G1 {
	x, y := hashToCurvePoint(x)
	p.SetXY(x,y)
	p.ScalarMult(new(Scalar).FromInt(G1_h))
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

		C.ep_set_infty(&p.st)
		return in[inlen:]
	}

	if compressed {
		bin[0] = 2
		C.ep_read_bin(&p.st, (*C.uint8_t)(&bin[0]), G1Size+1)
		var yneg C.fp_st

		if negativeIsBigger(&yneg[0], &p.st.y[0]) != (in[0]&serializationBigY != 0) {
			p.st.y = yneg
		}
		return in[G1Size:]
	}

	bin[0] = 4
	C.ep_read_bin(&p.st, (*C.uint8_t)(&bin[0]), G1UncompressedSize+1)
	return in[G1UncompressedSize:]
}

// Marshal the point, compressed to X and sign.
func (p *G1) Marshal() (res []byte) {
	var bin [G1Size + 1]byte
	res = bin[1:]

	if C.ep_is_infty(&p.st) == 1 {
		res[0] = serializationInfinity | serializationCompressed
		return
	}
	C.ep_norm(&p.st, &p.st)
	C.ep_write_bin((*C.uint8_t)(&bin[0]), G1Size+1, &p.st, 1)
	res[0] |= serializationCompressed

	var yneg C.fp_st
	if negativeIsBigger(&yneg[0], &p.st.y[0]) {
		res[0] |= serializationBigY
	}
	return
}

func (p *G1) String() string {
	x,y := p.GetXY()
	return fmt.Sprintf("bls12.G1(%d,%d)", x,y)
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
