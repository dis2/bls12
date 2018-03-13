package bls12

import (
	"bytes"
	"io/ioutil"
	"testing"
	"math/big"
	"crypto/sha512"
	"crypto/rand"
)

func BenchmarkUncompressG1(b *testing.B) {
	p1 := new(G1).HashToPoint([]byte("test2"))
	b.ResetTimer()
	m := p1.Marshal()
	for i := 0; i < b.N; i++ {
		p1.Unmarshal(m)
	}
}

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


func TestSetZero(t *testing.T) {
	new(G1).SetZero()
}

func TestSetOne(t *testing.T) {
	new(G1).SetOne()
}

func TestHashToPoint(t *testing.T) {
	// Via relic (hardcoded sha384)
	p1 := new(G1).HashToPoint([]byte("test2"))

	// Manually via custom hash
	buf := sha512.Sum384([]byte("test2"))
	x := new(big.Int).SetBytes(buf[:])
	p2 := new(G1).HashIntToPoint(x)
	if !p1.Equal(p2) {
		t.Fatal("misbehaving hashtopoint")
	}
}


func TestVectorG1Compressed(t *testing.T) {
	//	t.Run("Compressed", func(t *testing.T) {
	var (
		data = readFile(t, "testdata/g1_compressed_valid_test_vectors.dat")
		ep   = (&G1{}).SetZero()
		a    = &G1{}
		one  = (&G1{}).SetOne()
		d    = data
	)
	for i := 0; i < 1000; i++ {
		t.Logf("%d <- %x", i, d[:G1Size])
		ok := a.Unmarshal(d[:G1Size])
		if ok == nil {
			t.Errorf("%d: failed decoding", i)
		}
		if !ep.Equal(a) {
			t.Errorf("%d: different point", i)
		}
		buf := ep.Marshal()
		t.Logf("%d -> %x", i, buf)
		if !bytes.Equal(buf, d[:G1Size]) {
			t.Errorf("%d: different encoding", i)
		}
		d = d[G1Size:]
		ep.Add(one)
	}
	//	})
}
func TestVectorG1Uncompressed(t *testing.T) {
	//	t.Run("Uncompressed", func(t *testing.T) {
	t.Logf("setup\n")
	var (
		data = readFile(t, "testdata/g1_uncompressed_valid_test_vectors.dat")
		ep   = (&G1{}).SetZero()
		a    = &G1{}
		one  = (&G1{}).SetOne()
		d    = data
	)
	for i := 0; i < 1000; i++ {
		t.Logf("%d <- %x", i, d[:G1UncompressedSize])
		ok := a.Unmarshal(d[:G1UncompressedSize])
		if ok == nil {
			t.Errorf("%d: failed decoding",i)
		}
		if !ep.Equal(a) {
			t.Errorf("%d: different point", i)
		}
		buf := ep.MarshalUncompressed()
		t.Logf("%d -> %x", i, buf)
		if !bytes.Equal(buf, d[:G1UncompressedSize]) {
			t.Errorf("%d: different encoding", i)
		}
		d = d[G1UncompressedSize:]
		ep.Add(one)
	}
	//	})
}

func readFile(t *testing.T, name string) []byte {
	t.Helper()
	res, err := ioutil.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return res
}
