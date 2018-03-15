package bls12

// Abstract algebraic field for Fq (G1 coords), Fq2 (G2 coords) and Fq12 (GT)
type Field interface {
	Exp(a Field, n *Scalar) Field
	Square(x Field) Field
	Add(a,b Field) Field
	Mul(a,b Field) Field
	Sub(a,b Field) Field
	Inverse(x Field) Field
	Neg(x Field) Field
	Sqrt(a Field) bool
	Equal(x Field) bool
	GreaterThan(y Field) bool
	CopyParity(y Field) Field
	EnsureParity(bool) bool
	Cast(x *Fq) Field
	Marshal() []byte
	Unmarshal([]byte) []byte
	Y2FromX(Field) Field
	l() *Limb
	le(Field) *Limb
	Copy() Field
	Set(Field)
	New() Field
}
