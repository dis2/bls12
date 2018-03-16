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
func (p *G2) HashToPoint(msg []byte) G {
	state := sha3.NewShake256()
	state.Write([]byte("BLS12-381 G2"))
	state.Write(msg)

	var t Fq2
	var h [96]byte
	// Trim to 380 bits.
	state.Read(h[:])
	h[0] &= 0x0f
	h[48] &= 0x0f
	t.C[1].Unmarshal(t.C[0].Unmarshal(h[:]))
	x, y := FouqueMapXtoY(&t)
	p.SetXY(x, y)
	p.ScaleByCofactorFast()
	return p
}

func (p *G2) MapIntToPoint(in Field) G {
	x, y := MapXtoY(in.(*Fq2))
	p.SetXY(x, y)
	p.ScaleByCofactorFast()
	return p
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

