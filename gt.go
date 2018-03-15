package bls12

const (
	GTSize = 384
)

// Element of the q^12 extension field.
// Can be thought of as a pseudo-point in a group GT resulting from pairing
// operations.
//
// The group law is shifted from the underlying Fq12 field. Ie multiplication
// is field exponentation etc.
type Fq6 [3]Fq2
type GT [2]Fq6

func (p *GT) Copy() *GT {
	c := *p
	return &c
}
