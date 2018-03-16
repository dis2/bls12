package bls12

// Explicit formulas only.

// Fancy hash:
// https://www.di.ens.fr/~fouque/pub/latincrypt12.pdf
// https://github.com/herumi/mcl/blob/9fbc54305d01b984e39d83e96bfa94bb17648a86/include/mcl/bn.hpp#L76
//
// This computes only the raw mapping in a random subgroup. Scaling by cofactor
// and setting parity of y according to t is responsibility of the caller.
// input: t, output: x,y
func FouqueMapXtoY(t, x, y Field) {
	w, y2, c := t.Copy(), t.New(), t.New()
	// w = (t^2 + 4u + 1)^(-1) * sqrt(-3) * t
	//if t.IsZero() { panic("degenerate t=0") }
	w.Square(w)
	w.Add(w, c.GetB()) // 4u
	w.Add(w, c.Cast(&One))
	//if w.IsZero() { panic("degenerate t^2=-5") }
	w.Inverse(w)
	w.Mul(w, c.Cast(&QSqrtMinus3))
	w.Mul(w, t)

	for i := 0; i < 3; i++ {
		switch i {
		// x = (sqrt(-3) - 1) / 2 - (w * t)
		case 0:
			x.Mul(w, t)
			x.Sub(c.Cast(&QSqrtMinus3Minus1Half), x)
		// x = -1 - x
		case 1:
			x.Sub(c.Cast(&QMinus1), x)
		// x = 1/w^2 + 1
		case 2:
			x.Square(w)
			x.Inverse(x)
			x.Add(x, c.Cast(&One))
		}
		// y2 = y^2 = x^3 + 4u
		y2.Square(x)
		y2.Mul(y2, x)
		y2.Add(y2, c.GetB())
		// y = sqrt(y2)
		if y.Sqrt(y2) {
			return
		}
	}
	panic("Uh oh.")
}

// Simple map (RELIC): Bump x until y satisfies curve equation.
// This computes only the raw mapping in a random subgroup. Scaling by cofactor
// and setting parity of y according to t is responsibility of the caller.
// input: t, output: x,y
func MapXtoY(t, x, y Field) bool {
	y2, c := t.New(), t.New()
	x.Set(t)
	// There is no bound, make up arbitrary one
	for i := 0; i < 1000; i++ {
		// y2 = y^2 = x^3 + 4
		y2.Square(x)
		y2.Mul(y2, x)
		y2.Add(y2, c.GetB()) // 4u

		if y.Sqrt(y2) {
			return true
		}
		x.Add(x, c.Cast(&One))
	}
	return false
}
