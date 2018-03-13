package bls12

import (
	"crypto/rand"
	"testing"
)

func BenchmarkUncompressGT(b *testing.B) {
	x, _ := rand.Int(rand.Reader, Order)
	x2, _ := rand.Int(rand.Reader, Order)
	s := new(Scalar).FromInt(x)
	s2 := new(Scalar).FromInt(x2)
	g1 := new(G1).ScalarBaseMult(s)
	g2 := new(G2).ScalarBaseMult(s2)
	gt := GT{}
	gt.Pair(g1,g2)
	gt2 := GT{}
	m := gt.Marshal()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gt2.Unmarshal(m)
	}
}

func TestMarshalGT(t *testing.T) {
	for i:= 1; i < 1000; i++ {
		x, _ := rand.Int(rand.Reader, Order)
		x2, _ := rand.Int(rand.Reader, Order)
		s := new(Scalar).FromInt(x)
		s2 := new(Scalar).FromInt(x2)
		g1 := new(G1).ScalarBaseMult(s)
		g2 := new(G2).ScalarBaseMult(s2)
		gt := GT{}
		gt.Pair(g1,g2)
		gt2 := GT{}
		gt2.Unmarshal(gt.Marshal())
		if gt != gt2 {
			t.Logf("gt1: %v", gt)
			t.Logf("gt2: %v", gt2)
			t.Fatal("lossy marshalling")
		}
	}
}

func BenchmarkPairGT(b *testing.B) {
	g1 := new(G1).HashToPoint([]byte("x"))
	g2 := new(G2).HashToPoint([]byte("x"))
	e := GT{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Pair(g1,g2)
	}
}
