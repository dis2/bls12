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


