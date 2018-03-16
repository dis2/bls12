package bls12

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func BenchmarkUncompressG2(b *testing.B) {
	p2 := new(G2).HashToPoint([]byte("test2"))
	b.ResetTimer()
	m := p2.Marshal()
	for i := 0; i < b.N; i++ {
		p2.Unmarshal(m)
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

func BenchmarkHashToPointG2(b *testing.B) {
	var buf [512]byte
	var g2 G2
	for i := 0; i < b.N; i++ {
		buf[0] = byte(i)
		buf[1] = byte(i >> 8)
		buf[2] = byte(i >> 16)
		g2.HashToPoint(buf[:])
	}
}

func BenchmarkHashToPointRelicG2(b *testing.B) {
	var buf [512]byte
	var g2 G2
	for i := 0; i < b.N; i++ {
		buf[0] = byte(i)
		buf[1] = byte(i >> 8)
		buf[2] = byte(i >> 16)
		g2.HashToPointRelic(buf[:])
	}
}


func TestHashToPointG2(t *testing.T) {
	var p G2
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

func TestVectorG2Compressed(t *testing.T) {
	var (
		data = readFile(t, "testdata/g2_compressed_valid_test_vectors.dat")
		ep2  = G2Zero()
		a    = new(G2)
		one  = G2One()
		d    = data
	)
	for i := 0; i < 1000; i++ {
		ok := a.Unmarshal(d[:G2Size])
		if ok == nil {
			t.Errorf("%d: failed decoding", i)
		}
		ep2.Normalize()
		if !ep2.X.Equal(&a.X) {
			t.Errorf("%d: different X", i)
		}
		if !ep2.Y.Equal(&a.Y) {
			t.Errorf("%d: different Y", i)
		}
		if !ep2.Equal(a) {
			t.Errorf("%d: different point", i)
		}
		buf := ep2.Marshal()
		if !bytes.Equal(buf, d[:G2Size]) {
			t.Logf("%d <- %x", i, d[:G2Size])
			t.Logf("%d -> %x", i, buf)
			t.Errorf("%d: different encoding", i)
		}
		d = d[G2Size:]
		ep2.Add(one)
	}
}
func TestVectorG2Uncompressed(t *testing.T) {
	var (
		data = readFile(t, "testdata/g2_uncompressed_valid_test_vectors.dat")
		ep2  = G2Zero()
		a    = new(G2)
		one  = G2One()
		d    = data
	)
	for i := 0; i < 1000; i++ {
		//t.Logf("%d <- %x", i, d[:G2UncompressedSize])
		ok := a.Unmarshal(d[:G2UncompressedSize])
		if ok == nil {
			t.Errorf("%d: failed decoding", i)
		}
		if !ep2.Equal(a) {
			t.Errorf("%d: different point", i)
		}
		buf := ep2.MarshalUncompressed()
		//t.Logf("%d -> %x", i, buf)
		if !bytes.Equal(buf, d[:G2UncompressedSize]) {
			t.Errorf("%d: different encoding", i)
		}
		d = d[G2UncompressedSize:]
		ep2.Add(one)
	}
}
