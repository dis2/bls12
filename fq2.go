package bls12

import "math/big"

// Fq holds C0, C1, where C1 * u + C0 is the coordinate
type Fq2 struct {
	C [2]Fq
}

func (e *Fq2) opt(x Field) *Fq2 {
	v, ok := x.(*Fq2)
	if !ok && x != nil {
		panic("invalid field type passed")
	}
	if v == nil {
		v = e
	}
	return v
}

func (e *Fq2) GetB() Field {
	e.C[0].GetB()
	e.C[1].GetB()
	return e
}

func (e *Fq2) New() Field {
	return &Fq2{}
}

func (e *Fq2) Set(x Field) {
	*e = *x.(*Fq2)
}

// parity(e) = parity(n)
// Where parity is a "sign" of y coordinate, defined as:
// parity 1 - neg(x).C1 > x.C1
// parity 0 - neg(x).C1 <= x.C1
func (e *Fq2) CopyParity(y Field) Field {
	var nege, negy Fq2
	nege.Neg(e)
	negy.Neg(y)
	if nege.GreaterThan(e) != negy.GreaterThan(y) {
		*e = nege
	}
	return e
}

// Note that this is lexicographic "greater than", it bears no relevance
// to actual order (as result of scalarmult etc). Notably, you can't just
// skimp on it and use QMinus1Half to decide parity with this, you have to
// perform actual negation in the group, and then compare those two.
func (e *Fq2) GreaterThan(y Field) bool {
	return e.C[1].GreaterThan(&y.(*Fq2).C[1])
}

// Ensures parity p. It also returns the parity e had prior to this call.
func (e *Fq2) EnsureParity(p bool) bool {
	var nege Fq2
	nege.Neg(e)
	// The negative is larger
	if nege.GreaterThan(e) {
		if p {
			// And we want it set
			*e = nege
		}
		return false
		// The negative is smaller
	} else {
		if !p {
			// And we want it set
			*e = nege
		}
		return true
	}
}

func (e *Fq2) Copy() Field {
	t := *e
	return &t
}

// e == x
func (e *Fq2) Equal(x Field) bool {
	return e.C[0].Equal(&x.(*Fq2).C[0]) && e.C[1].Equal(&x.(*Fq2).C[1])
}

func (e *Fq2) IsZero() bool {
	var tmp Fq2
	return e.Equal(tmp.Cast(&Zero))
}

// e = x^3
func (e *Fq2) Cube(x Field) Field {
	xx := *e.opt(x)
	e.Square(&xx)
	e.Mul(e, &xx)
	return e
}

// e = 64 bit immediate n
func (e *Fq2) SetInt64(n int64) Field {
	e.C[0].SetInt64(n)
	e.C[1].SetInt64(n)
	return e
}

func (e *Fq2) Y2FromX(x Field) Field {
	xx := *e.opt(x)
	var tmp Fq2
	e.Square(&xx)
	e.Mul(e, &xx)
	e.Add(e, tmp.GetB())
	return e
}

// Cast Fq element into Fq2. This is used to resolve constants in Fq2.
// The contract is to never depend on value of 'e', and always use only
// whatever is returned.
func (tmp *Fq2) Cast(v *Fq) Field {
	tmp.C[0] = *v
	tmp.C[1] = Zero
	return tmp
}

func (e *Fq2) Unmarshal(b []byte) []byte {
	if len(b) < 96 {
		return nil
	}
	return e.C[0].Unmarshal(e.C[1].Unmarshal(b[:]))
}

func (e *Fq2) Marshal() []byte {
	return append(e.C[1].Marshal(), e.C[0].Marshal()...)
}

func (e *Fq2) FromInt(i []*big.Int) Field {
	e.C[0].FromInt(i[0:1])
	e.C[1].FromInt(i[1:2])
	return e
}

func (e *Fq2) ToInt() []*big.Int {
	return append(e.C[0].ToInt(), e.C[1].ToInt()[0])
}

func (e *Fq2) IsResidue() bool {
	var t0, t1 Fq
	t0.Square(&e.C[0])
	t1.Square(&e.C[1])
	return t0.Add(nil,&t1).Mul(nil, &QMinus1Half).Equal(&One)
}

