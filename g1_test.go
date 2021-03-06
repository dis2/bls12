package bls12

import (
	"bytes"
	"fmt"
	"crypto/rand"
	"io/ioutil"
	"testing"
)

func BenchmarkUncompressG1(b *testing.B) {
	p1 := new(G1).HashToPointFast([]byte("test2"))
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
	g1 := new(G2).HashToPointFast([]byte("yxxx"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g1.ScalarMult(s)
	}
}

func BenchmarkHashToPointG1(b *testing.B) {
	var buf [512]byte
	var g1 G1
	for i := 0; i < b.N; i++ {
		buf[0] = byte(i)
		buf[1] = byte(i >> 8)
		buf[2] = byte(i >> 16)
		g1.HashToPoint(buf[:])
	}
}

func BenchmarkHashToPointFastG1(b *testing.B) {
	var buf [512]byte
	var g1 G1
	for i := 0; i < b.N; i++ {
		buf[0] = byte(i)
		buf[1] = byte(i >> 8)
		buf[2] = byte(i >> 16)
		g1.HashToPointFast(buf[:])
	}
}

func BenchmarkHashToPointRelicG1(b *testing.B) {
	var buf [512]byte
	var g1 G1
	for i := 0; i < b.N; i++ {
		buf[0] = byte(i)
		buf[1] = byte(i >> 8)
		buf[2] = byte(i >> 16)
		g1.HashToPointRelic(buf[:])
	}
}

func TestHashToPointG1(t *testing.T) {
	var p G1
	var buf [512]byte
	for i := 0; i < 1000; i++ {
		buf[0] = byte(i)
		buf[1] = byte(i >> 8)
		buf[2] = byte(i >> 16)
		p.HashToPoint(buf[:])
		if !p.Check() {
			t.Fatalf("point landed in wrong subgroup for %d\n", i)
		}
	}
}

func TestHashToPointFastG1(t *testing.T) {
	var p G1
	var buf [512]byte
	for i := 0; i < 1000; i++ {
		buf[0] = byte(i)
		buf[1] = byte(i >> 8)
		buf[2] = byte(i >> 16)
		p.HashToPointFast(buf[:])
		if !p.Check() {
			t.Fatalf("point landed in wrong subgroup for %d\n", i)
		}
	}
}

func TestVectorG1HashToPoint(t *testing.T) {
	data := readFile(t, "testdata/g1_hashtopoint.dat")
	for i := 0; i < 1000; i++ {
		var g1, g1t G1
		data = g1.Unmarshal(data)
		if data == nil {
			t.Fatal("failed to unmarshal test data")
		}
		g1t.HashToPoint([]byte(fmt.Sprintf("%d", i)))
		if !g1t.Equal(&g1) {
			t.Fatalf("wrong hash at %d\n", i)
		}
	}
}

func TestVectorG1Compressed(t *testing.T) {
	//	t.Run("Compressed", func(t *testing.T) {
	var (
		data = readFile(t, "testdata/g1_compressed_valid_test_vectors.dat")
		ep   = G1Zero()
		a    = &G1{}
		one  = G1One()
		d    = data
	)
	for i := 0; i < 1000; i++ {
		ok := a.Unmarshal(d[:G1Size])
		if ok == nil {
			t.Errorf("%d: failed decoding", i)
		}
		if !ep.Equal(a) {
			t.Errorf("%d: different point", i)
		}
		buf := ep.Marshal()
		if !bytes.Equal(buf, d[:G1Size]) {
			t.Logf("%d <- %x", i, d[:G1Size])
			t.Logf("%d -> %x", i, buf)
			t.Errorf("%d: different encoding", i)
		}
		d = d[G1Size:]
		ep.Add(one)
	}
	//	})
}
func TestVectorG1Uncompressed(t *testing.T) {
	//	t.Run("Uncompressed", func(t *testing.T) {
	var (
		data = readFile(t, "testdata/g1_uncompressed_valid_test_vectors.dat")
		ep   = G1Zero()
		a    = &G1{}
		one  = G1One()
		d    = data
	)
	for i := 0; i < 1000; i++ {
		ok := a.Unmarshal(d[:G1UncompressedSize])
		if ok == nil {
			t.Errorf("%d: failed decoding", i)
		}
		if !ep.Equal(a) {
			t.Errorf("%d: different point", i)
		}
		buf := ep.MarshalUncompressed()
		if !bytes.Equal(buf, d[:G1UncompressedSize]) {
			t.Logf("%d <- %x", i, d[:G1UncompressedSize])
			t.Logf("%d -> %x", i, buf)
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
