package bls12

import "testing"
import "math/big"
import "crypto/rand"

func rnd() (f Fq, i *big.Int) {
	i,_ = rand.Int(rand.Reader, Q.RawInt())
	f.FromInt(i)
	return
}

func TestFqAdd(t *testing.T) {
	for i := 0; i < 10000; i++ {
		a, ia:= rnd()
		b, ib:= rnd()
		a.Add(nil, &b)
		ia.Add(ia, ib)
		ia.Mod(ia, Q.RawInt())
		if a.ToInt().Cmp(ia) != 0 {
			t.Fatal("failed addition")
		}
	}
}

func TestFqSub(t *testing.T) {
	for i := 0; i < 10000; i++ {
		a, ia:= rnd()
		b, ib:= rnd()
		a.Sub(nil, &b)
		ia.Sub(ia, ib)
		ia.Mod(ia, Q.RawInt())
		if a.ToInt().Cmp(ia) != 0 {
			t.Fatal("failed substraction")
		}
	}
}

func TestFqNeg(t *testing.T) {
	for i := 0; i < 10000; i++ {
		a, ia:= rnd()
		a.Neg(&a)
		ia.Neg(ia)
		ia.Mod(ia, Q.RawInt())
		if a.ToInt().Cmp(ia) != 0 {
			t.Fatal("failed neg")
		}
	}
}

func TestFqMul(t *testing.T) {

	for i := 0; i < 10000; i++ {
		a, ia:= rnd()
		b, ib:= rnd()
		a.Mul(nil, &b)
		ia.Mul(ia, ib)
		ia.Mod(ia, Q.RawInt())
		if a.ToInt().Cmp(ia) != 0 {
			t.Fatal("failed mult")
		}
	}
}

func TestFqSquare(t *testing.T) {
	for i := 0; i < 10000; i++ {
		a, ia:= rnd()
		a.Square(&a)
		ia.Mul(ia, ia)
		ia.Mod(ia, Q.RawInt())
		if a.ToInt().Cmp(ia) != 0 {
			t.Fatal("failed square")
		}
	}
}

