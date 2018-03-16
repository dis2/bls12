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

// Hash to point (uses SHA3), reasonably fast.
func (p *G1) HashToPointFast(msg []byte) G {
	state := sha3.NewShake256()
	state.Write(msg)
	state.Write([]byte("BLS12-381 G1"))
	var buf [G1Size]byte
	state.Read(buf[:])
	return p.HashToPointBytes(&buf)
}

// Hash to point using custom 48 bytes of hash.
// This input MUST be uniformly random output of a secure hash function!
// The function can mutate the passed buffer.
func (p *G1) HashToPointBytes(buf *[G1Size]byte) G {
	// Trim to 380bits
	buf[0] &= 0xf
	var t Fq
	t.Unmarshal(buf[:])
	FouqueMapXtoY(&t, &p.X, &p.Y)
	// match parity of y with t
	if t.Limbs[0]&1 != p.Y.Limbs[0]&1 {
		p.Y.Neg(&p.Y)
	}
	p.SetNormalized()
	p.ScaleByCofactor()
	return p
}

// Map arbitrary integer to a point, for use with custom hash function.
func (p *G1) MapIntToPoint(in Field) bool {
	if MapXtoY(in.(*Fq), &p.X, &p.Y) {
		p.ScaleByCofactor()
		return true
	}
	return false
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
