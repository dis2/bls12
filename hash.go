package bls12
// Explicit formulas only.
//
// Fancy hash: 
// https://www.di.ens.fr/~fouque/pub/latincrypt12.pdf
// https://github.com/herumi/mcl/blob/9fbc54305d01b984e39d83e96bfa94bb17648a86/include/mcl/bn.hpp#L76
//
// This is called twice, with keyed hash. Scale both resulting points by
// cofactor, and add them together for final point.
func fouqueHalfPoint(t *Fq) (x, y Fq) {
	var w, y2, ytest Fq
	// w = (t^2 + 4 + 1)^(-1) * sqrt(-3) * t
	//if t.IsZero() { panic("degenerate t=0") }
	w = *t
	w.Square(&w)
	w.Add(&w, &Five)
	//if w.IsZero() { panic("degenerate t^2=-5") }
	w.Inverse(&w)
	w.Mul(&w, &QSqrtMinus3)
	w.Mul(&w, t)
	for i := 0; i < 3; i++ {
		switch i {
		// x = (sqrt(-3) - 1) / 2 - (w * t)
		case 0: x.Mul(&w, t)
			x.Sub(&QSqrtMinus3Minus1Half, &x)
		// x = -1 - x
		case 1: x.Sub(&QMinus1, &x)
		// x = 1/w^2 + 1
		case 2: x.Square(&w)
			x.Inverse(&x)
			x.Add(&x,&One)
		}
		// y2 = y^2 = x^3 + 4
		y2.Square(&x)
		y2.Mul(&y2,&x)
		y2.Add(&y2, &Four)
		// y = y2 ^ ((q+1)/4)
		y.Exp(&y2, &QPlus1Quarter)
		// if y^2 == y2
		if y2.Equal(ytest.Square(&y)) {
			var residue Fq
			residue.Exp(t, &QMinus1Half)
			//if residue.IsZero() { panic("degenerate residue") }
			if !residue.Equal(&One) { // Must be non-residue.
				y.Neg(&y)
			}
			return x, y
		}
	}
	panic("Uh oh.")
}

// Simple map (RELIC): Bump x until y satisfies curve equation. Again, caller is
// responsible for scaling by h.
func mapXtoY(t *Fq) (x, y Fq) {
	x = *t
	for {
		var y2, ytest Fq
		// y2 = y^2 = x^3 + 4
		y2.Square(&x)
		y2.Mul(&y2,&x)
		y2.Add(&y2, &Four)
		// y = y2 ^ ((q+1)/4)
		y.Exp(&y2, &QPlus1Quarter)
		// if y^2 == y2
		if y2.Equal(ytest.Square(&y)) {
			return
		}
		x.Add(&x, &One)
	}
}


