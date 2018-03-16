// Package bls12 implements bilinear pairing curve BLS12-381
package bls12 // import "github.com/dis2/bls12"

// #include "relic_core.h"
// #include "relic_ep.h"
// void _fp_neg(fp_t r, const fp_t p) { fp_neg(r, p); }
// void _bn_new(bn_t bn) { bn_new(bn); }
import "C"
import "bytes"
import "math/big"
import "encoding/hex"

// Check if encoding of negated coordinate is lexicographically bigger.
// Also returns the given negate.
func negativeIsBigger(neg, a *C.dig_t) bool {
	C._fp_neg(neg, a)

	var abuf, bbuf [48]byte
	C.fp_write_bin((*C.uint8_t)(&abuf[0]), 48, a)
	C.fp_write_bin((*C.uint8_t)(&bbuf[0]), 48, neg)
	return bytes.Compare(abuf[:], bbuf[:]) > 0
}

func hexConst(s string) (ret *big.Int) {
	ret, _ = new(big.Int).SetString(s, 16)
	return
}

var init_done = false

func init_pending() {
	if init_done {
		return
	}
	init_done = true
	C.core_init()
	C.ep_param_set_any_pairf()
}

func QConst(s string) (f Fq) {
	if len(s)%2 != 0 {
		panic("bad const padding for " + s)
	}
	init_pending()
	var buf [48]byte
	pad := 48 - len(s)/2
	_, err := hex.Decode(buf[pad:], []byte(s))
	if err != nil || f.Unmarshal(buf[:]) == nil {
		panic("invalid const " + s)
	}
	return
}

const (
	// https://github.com/ebfull/pairing/tree/master/src/bls12_381#serialization
	serializationMask       = (1 << 5) - 1
	serializationCompressed = 1 << 7 // 0x80
	serializationInfinity   = 1 << 6 // 0x40
	serializationBigY       = 1 << 5 // 0x20
)

func unmarshalG(p marshallerG, in []byte) (res []byte) {
	size := p.getSize()
	if len(in) < size {
		return nil
	}
	var bin = make([]byte, size)
	copy(bin[:], in)
	flags := bin[0]
	bin[0] &= serializationMask

	compressed := flags&serializationCompressed != 0
	inlen := size * 2
	if compressed {
		inlen = size
	} else if len(in) < inlen {
		return nil
	}
	res = in[inlen:]

	// Big Y, but we're not compressed, or infinity is serialized
	if (flags&serializationBigY != 0) && (!compressed || (flags&serializationInfinity != 0)) {
		return nil
	}

	if flags&serializationInfinity != 0 {
		// Check that rest is zero
		for _, v := range in[1:inlen] {
			if v != 0 {
				return nil
			}
		}
		p.SetZero()
		return res
	}

	X, Y, _ := p.GetXYZ()
	X.Unmarshal(bin[:])
	if compressed {
		if !Y.Y2FromX(X).Sqrt(nil) {
			return nil
		}
		Y.EnsureParity(flags&serializationBigY != 0)
	} else {
		Y.Unmarshal(in[size : size*2])
	}
	p.SetNormalized()
	if !p.Check() {
		return nil
	}
	return res

}

// Marshal the point, compressed to X and sign.
func marshalG(p marshallerG, comp int) (res []byte) {
	p.Normalize()
	X, Y, _ := p.GetXYZ()
	if p.IsZero() {
		res = make([]byte, p.getSize()*comp)
		res[0] = serializationInfinity
		if comp == 1 {
			res[0] |= serializationCompressed
		}
		return
	}
	res = X.Marshal()
	if comp == 1 {
		if Y.Copy().EnsureParity(false) {
			res[0] |= serializationBigY
		}
		res[0] |= serializationCompressed
	} else {
		res = append(res, Y.Marshal()...)
	}
	return
}

type marshallerG interface {
	getSize() int
	GetXYZ() (x, y, z Field)
	Check() bool
	SetZero()
	SetOne()
	Normalize()
	SetNormalized()
	IsZero() bool
}
