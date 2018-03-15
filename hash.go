package bls12
//import "fmt"
// Explicit formulas only.

// Fancy hash: 
// https://www.di.ens.fr/~fouque/pub/latincrypt12.pdf
// https://github.com/herumi/mcl/blob/9fbc54305d01b984e39d83e96bfa94bb17648a86/include/mcl/bn.hpp#L76
//
// This computes only the raw mapping in a random subgroup. Scaling by cofactor
// and setting parity of y according to t is responsibility of the caller.
func FouqueMapXtoY(t Field) (x, y Field) {
	x, y = t.New(), t.New()
	w, y2, ytest, c := t.Copy(), t.New(), t.New(), t.New()
	// w = (t^2 + 4 + 1)^(-1) * sqrt(-3) * t
	//if t.IsZero() { panic("degenerate t=0") }
	w.Square(w)
	w.Add(w, c.Cast(&Five))
	//if w.IsZero() { panic("degenerate t^2=-5") }
	w.Inverse(w)
	w.Mul(w, c.Cast(&QSqrtMinus3))
	w.Mul(w, t)

	for i := 0; i < 3; i++ {
		switch i {
		// x = (sqrt(-3) - 1) / 2 - (w * t)
		case 0: x.Mul(w, t)
			x.Sub(c.Cast(&QSqrtMinus3Minus1Half), x)
		// x = -1 - x
		case 1: x.Sub(c.Cast(&QMinus1), x)
		// x = 1/w^2 + 1
		case 2: x.Square(w)
			x.Inverse(x)
			x.Add(x,c.Cast(&One))
		}
		// y2 = y^2 = x^3 + 4
		y2.Square(x)
		y2.Mul(y2,x)
		y2.Add(y2, c.Cast(&Four))
		// y = y2 ^ ((q+1)/4)
		y.Exp(y2, &QPlus1Quarter)
		// if y^2 == y2
		if y2.Equal(ytest.Square(y)) {
			return x, y
		}
	}
	panic("Uh oh.")
}

// Simple map (RELIC): Bump x until y satisfies curve equation.
// This computes only the raw mapping in a random subgroup. Scaling by cofactor
// and setting parity of y according to t is responsibility of the caller.
func MapXtoY(t Field) (x, y Field) {
	x, y = t.Copy(), t.New()
	y2, ytest, c := t.New(), t.New(), t.New()
	for {
		// y2 = y^2 = x^3 + 4
		y2.Square(x)
		y2.Mul(y2,x)
		y2.Add(y2, c.Cast(&Four))
		// y = y2 ^ ((q+1)/4)
		y.Exp(y2, &QPlus1Quarter)
		// if y^2 == y2
		if y2.Equal(ytest.Square(y)) {
			return
		}
		x.Add(x, c.Cast(&One))
	}
}


