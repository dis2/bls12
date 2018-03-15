package bls12

import "fmt"
import "golang.org/x/crypto/sha3"

// Point on G1 y^2 = x^3 + 4 represented by x and y Fq coords.
// This structure holds either affine or (extended) projective form on ad-hoc
// basis signalled by the Norm flag. All operations can deal with either form.
// You can force point to become affine with Normalize().
type G1 struct {
	X, Y, Z Fq
	Norm    bool
}

func (p *G1) Copy() *G1 {
	c := *p
	return &c
}

// Set affine coordinates X,Y with implicit Z=1
func (p *G1) SetXY(x, y Field) {
	p.X = *x.(*Fq)
	p.Y = *y.(*Fq)

	// Implicitly normalized
	p.SetNormalized()
}

// Make the point explicitly normalized (ie after manually editing X/Y)
func (p *G1) SetNormalized() {
	p.Z = One
	p.Norm = true
}

// Get a copy of coordinates of the element.
// You may need to call Normalize() first if you want affine and ignore Z.
func (p *G1) GetXYZ() (x, y, z Field) {
	return &p.X, &p.Y, &p.Z
}

// p = G1_h * G1(p)
func (p *G1) ScaleByCofactor() {
	p.ScalarMult(&G1_h)
}

// Hash to point
func (p *G1) HashToPoint(msg []byte) *G1 {
	var h [48]byte
	var t Fq
	state := sha3.NewShake256()
	state.Write([]byte("BLS12-381 G1"))
	state.Write(msg)
	state.Read(h[:])
	// trim to 380 bits
	h[0] &= 0x0f
	t.Unmarshal(h[:])
	x, y := FouqueMapXtoY(&t)
	y.CopyParity(&t)
	p.SetXY(x, y)
	p.ScaleByCofactor()
	return p
}

// Map arbitrary integer to a point, for use with custom hash function.
func (p *G1) MapIntToPoint(in *Fq) *G1 {
	x, y := MapXtoY(in)
	p.SetXY(x, y)
	p.ScaleByCofactor()
	return p
}

// Check that the point is on the curve and in correct subgroup.
func (p *G1) Check() bool {
	p.Normalize()
	y2 := new(Fq).Y2FromX(&p.X)
	return new(Fq).Square(&p.Y).Equal(y2) && p.Copy().ScalarMult(&R).IsZero()
}

const (
	G1Size             = 48
	G1UncompressedSize = 2 * G1Size
)

func (p *G1) String() string {
	return fmt.Sprintf("bls12.G1(%x)", p.Marshal())
}

func (p *G1) Unmarshal(in []byte) []byte {
	return unmarshalG(marshallerG(p), in)
}

func (p *G1) Marshal() []byte {
	return marshalG(marshallerG(p), 1)
}

func (p *G1) MarshalUncompressed() []byte {
	return marshalG(marshallerG(p), 2)
}

func (p *G1) getSize() int {
	return G1Size
}

/*
// Unmarshal a point on G1. It consumes either G1Size or
// G1UncompressedSize, depending on how the point was marshalled.
func (p *G1) Unmarshal(in []byte) (res []byte) {
	if len(in) < G1Size {
		return nil
	}
	var bin [G1Size]byte
	copy(bin[:], in)
	flags := bin[0]
	bin[0] &= serializationMask

	compressed := flags&serializationCompressed != 0
	inlen := G1UncompressedSize
	if compressed {
		inlen = G1Size
		res = in[inlen:]
	}
	if !compressed && len(in) < G1UncompressedSize {
		return nil
	}


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

	if compressed {
		p.X.Unmarshal(bin[:])
		if !p.Y.Y2FromX(&p.X).Sqrt(nil) {
			return nil
		}
		p.Y.EnsureParity(flags&serializationBigY!=0)
	} else {
		p.Y.Unmarshal(in[G1Size:G1UncompressedSize])
	}
	if !p.Check() {
		return nil
	}
	return res
}

// Marshal the point, compressed to X and sign.
func (p *G1) Marshal() (res []byte) {
	var bin [G1Size + 1]byte
	res = bin[1:]

	if p.IsZero() {
		res[0] = serializationInfinity | serializationCompressed
		return
	}
	p.Normalize()
	res = p.X.Marshal()
	res[0] |= serializationCompressed

	if p.Y.Copy().EnsureParity(false) {
		res[0] |= serializationBigY
	}
	return
}

// Marshal the point, as uncompressed XY.
func (p *G1) MarshalUncompressed() (res []byte) {
	if p.IsZero() {
		var buf [G1UncompressedSize]byte
		buf[0] = serializationInfinity
		return buf[:]
	}
	p.Normalize()
	return append(p.X.Marshal(), p.Y.Marshal()...)
}
*/
