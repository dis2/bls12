package bls12

import "fmt"
import "golang.org/x/crypto/sha3"

// Point on G2 y^2 = x^3 + 4(u+1) represented by x and y Fq2 coords.
// This structure holds either affine or (extended) projective form on ad-hoc
// basis signalled by the Norm flag. All operations can deal with either form.
// You can force point to become affine with Normalize().
type G2 struct {
	X, Y, Z Fq2
	Norm    bool
}

func (p *G2) Copy() G {
	c := *p
	return &c
}

// Set affine coordinates X,Y with implicit Z=1
func (p *G2) SetXY(x, y Field) G {
	p.X = *x.(*Fq2)
	p.Y = *y.(*Fq2)
	p.SetNormalized()
	return p
}

// Make the point explicitly normalized (ie after manually editing X/Y)
func (p *G2) SetNormalized() G {
	p.Z.C[1], p.Z.C[0] = Zero, One
	p.Norm = true
	return p
}

// Get interface pointers to XYZ coordinates.
// You may need to call Normalize() first if you want affine and ignore Z.
func (p *G2) GetXYZ() (x, y, z Field) {
	return &p.X, &p.Y, &p.Z
}

// HashToPoint the message.
func (p *G2) HashToPointFast(msg []byte) G {
	state := sha3.NewShake256()
	state.Write([]byte("BLS12-381 G2"))
	state.Write(msg)

	var t Fq2
	h0 := t.C[0].Buf()
	h1 := t.C[1].Buf()
	state.Read(h0[:])
	state.Read(h1[:])
	// Trim to 380 bits.
	h0[47] &= 0x0f
	h1[47] &= 0x0f
	toEndian(h0)
	toEndian(h1)
	FouqueMapXtoY(&t, &p.X, &p.Y)
	// match parity of y with t
	if t.C[0].Limbs[0]&1 != p.Y.C[0].Limbs[0]&1 {
		p.Y.Neg(&p.Y)
	}
	p.SetNormalized()
	p.ScaleByCofactorFast()
	return p
}

func (p *G2) MapIntToPoint(in Field) bool {
	if MapXtoY(in.(*Fq2), &p.X, &p.Y) {
		p.ScaleByCofactorFast()
		return true
	}
	return false
}

const (
	G2Size             = 96
	G2UncompressedSize = 2 * G2Size
)

// Check that the point is on the curve and in correct subgroup.
func (p *G2) Check() bool {
	p.Normalize()
	y2 := new(Fq2).Y2FromX(&p.X)
	return new(Fq2).Square(&p.Y).Equal(y2) && p.Copy().ScalarMult(&R).IsZero()
}

func (p *G2) String() string {
	return fmt.Sprintf("bls12.G2(%x)", p.Marshal())
}

// Unmarshal point from input slice, returns unconsumed remainder of the slice
// (depends on compression flag).
func (p *G2) Unmarshal(in []byte) []byte {
	return GUnmarshal(p, in)
}

func (p *G2) Marshal() []byte {
	return GMarshal(p, 1)
}

func (p *G2) MarshalUncompressed() []byte {
	return GMarshal(p, 2)
}

// Get (compressed) size, uncompressed is twice that. For interfaces.
func (p *G2) GetSize() int {
	return G2Size
}

