// Package bls12 implements bilinear pairing curve BLS12-381
package bls12 // import "github.com/dis2/bls12"

import "math/big"
import "encoding/hex"

func hexConst(s string) (ret *big.Int) {
	ret, _ = new(big.Int).SetString(s, 16)
	return
}

func QConst(s string) (f Fq) {
	if len(s)%2 != 0 {
		panic("bad const padding for " + s)
	}
	initPending()
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
