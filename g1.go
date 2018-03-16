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

func (p *G1) Copy() G {
	c := *p
	return &c
}

// Set affine coordinates X,Y with implicit Z=1
func (p *G1) SetXY(x, y Field) G {
	p.X = *x.(*Fq)
	p.Y = *y.(*Fq)

	// Implicitly normalized
	p.SetNormalized()
	return p
}

// Make the point explicitly normalized (ie after manually editing X/Y)
func (p *G1) SetNormalized() G {
	p.Z = One
	p.Norm = true
	return p
}

// Get a copy of coordinates of the element.
// You may need to call Normalize() first if you want affine and ignore Z.
func (p *G1) GetXYZ() (x, y, z Field) {
	return &p.X, &p.Y, &p.Z
}

// p = G1_h * G1(p)
func (p *G1) ScaleByCofactor() G {
	p.ScalarMult(&G1_h)
	return p
}

// Hash to point
func (p *G1) HashToPoint(msg []byte) G {
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
	p.SetXY(x, y)
	p.ScaleByCofactor()
	return p
}

// Map arbitrary integer to a point, for use with custom hash function.
func (p *G1) MapIntToPoint(in Field) G {
	x, y := MapXtoY(in.(*Fq))
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

// Unmarshal point from input slice, returns unconsumed remainder of the slice
// (depends on compression flag).
func (p *G1) Unmarshal(in []byte) []byte {
	return GUnmarshal(p, in)
}

func (p *G1) Marshal() []byte {
	return GMarshal(p, 1)
}

func (p *G1) MarshalUncompressed() []byte {
	return GMarshal(p, 2)
}

// Get (compressed) size, uncompressed is twice that. For interfaces.
func (p *G1) GetSize() int {
	return G1Size
}


