package bls12

import (
	"testing"
	"crypto/rand"
)

func BenchmarkBaseMultG1(b *testing.B) {
	x, _ := rand.Int(rand.Reader, Order)
	s := new(Scalar).FromInt(x)
	g1 := G1{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g1.ScalarBaseMult(s)
	}
}

func BenchmarkMultG1(b *testing.B) {
	x, _ := rand.Int(rand.Reader, Order)
	s := new(Scalar).FromInt(x)
	g1 := new(G2).HashToPoint([]byte("yxxx"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g1.ScalarMult(s)
	}
}

func BenchmarkMultG2(b *testing.B) {
	x, _ := rand.Int(rand.Reader, Order)
	s := new(Scalar).FromInt(x)
	g1 := new(G2).HashToPoint([]byte("xxx"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g1.ScalarMult(s)
	}
}

func BenchmarkPair(b *testing.B) {
	g1 := new(G1).HashToPoint([]byte("x"))
	g2 := new(G2).HashToPoint([]byte("x"))
	e := GT{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Pair(g1,g2)
	}
}
