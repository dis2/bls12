package bls12

import "fmt"
import "math/big"

type Limbs [NLimbs]Limb

type Fq struct {
	Limbs
}

func (e *Fq) opt(x Field) *Fq {
	v, ok := x.(*Fq)
	if !ok && x != nil {
		panic("invalid field type passed")
	}
	if v == nil {
		v = e
	}
	return v
}

// parity(e) = parity(n)
// Where parity is a "sign" of y coordinate, defined as:
// parity 1 - neg(x) > x
// parity 0 - neg(x) <= x
func (e *Fq) CopyParity(y Field) Field {
	if e.GreaterThan(&QMinus1Half) != e.GreaterThan(&QMinus1Half) {
		e.Neg(e)
	}
	return e
}

// Ensures parity p. It also returns the parity e had prior to this call.
func (e *Fq) EnsureParity(p bool) bool {
	var t Fq
	t.Neg(e)
	// The negative is larger
	if e.GreaterThan(&t) {
		if p {
			// And we want it set
			*e = t
		}
		return false
	// The negative is smaller
	} else {
		if !p {
			// And we want it set
			*e = t
		}
		return true
	}
}

func (e *Fq) Copy() Field {
	t := *e
	return &t
}

func (e *Fq) New() Field {
	return &Fq{}
}

func (e *Fq) Set(x Field) {
	*e = *x.(*Fq)
}


// Cast Fq into .. Fq. Silly, but needed to satisfy interface.
// The contract is to never depend on value of 'tmp', and always use only
// whatever is returned.
func (tmp *Fq) Cast(v *Fq) Field {
	return v
}

func (e *Fq) IsZero() bool {
	return e.Equal(&Zero)
}

func (e *Fq) Sqrt(a Field) bool {
	aa := e.opt(a)
	chk := *aa
	e.Exp(a, &QPlus1Quarter)
	return chk.Equal(e.Copy().Square(nil))
}

func (e *Fq) Y2FromX(x Field) Field {
	xx := *e.opt(x)
	e.Square(&xx)
	e.Mul(e, &xx)
	e.Add(e, &Four)
	return e
}


func pad(buf []byte) []byte {
	n := len(buf)
	if n > 48 {
		return buf
	}
	return append(make([]byte, 48-n), buf...)
}

func (e *Fq) FromInt(b *big.Int) *Fq {
	if e.Unmarshal(pad(b.Bytes())) == nil {
		return nil
	}
	return e
}

func (e *Fq) ToInt() *big.Int {
	return new(big.Int).SetBytes(e.Marshal())
}

func (e *Fq) String() string {
	return fmt.Sprintf("Fq(%d)", e.ToInt())
}
