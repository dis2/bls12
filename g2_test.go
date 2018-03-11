package bls12

import (
	"bytes"
	"testing"
)

func TestVectorG2Compressed(t *testing.T) {
	var (
		data = readFile(t, "testdata/g2_compressed_valid_test_vectors.dat")
		ep2  = new(G2).SetZero()
		a    = new(G2)
		one  = new(G2).SetOne()
		d    = data
	)
	for i := 0; i < 1000; i++ {
		t.Logf("%d <- %x", i, d[:G2Size])
		_, err := a.Unmarshal(d[:G2Size])
		if err != nil {
			t.Errorf("%d: failed decoding: %v", i, err)
		}
		if !ep2.Equal(a) {
			t.Errorf("%d: different point", i)
		}
		buf := ep2.Marshal()
		t.Logf("%d -> %x", i, buf)
		if !bytes.Equal(buf, d[:G2Size]) {
			t.Errorf("%d: different encoding", i)
		}
		d = d[G2Size:]
		ep2.Add(one)
	}
}
func TestVectorG2Uncompressed(t *testing.T) {
	var (
		data = readFile(t, "testdata/g2_uncompressed_valid_test_vectors.dat")
		ep2  = new(G2).SetZero()
		a    = new(G2)
		one  = new(G2).SetOne()
		d    = data
	)
	for i := 0; i < 1000; i++ {
		t.Logf("%d <- %x", i, d[:G2UncompressedSize])
		_, err := a.Unmarshal(d[:G2UncompressedSize])
		if err != nil {
			t.Errorf("%d: failed decoding: %v", i, err)
		}
		if !ep2.Equal(a) {
			t.Errorf("%d: different point", i)
		}
		buf := ep2.MarshalUncompressed()
		t.Logf("%d -> %x", i, buf)
		if !bytes.Equal(buf, d[:G2UncompressedSize]) {
			t.Errorf("%d: different encoding", i)
		}
		d = d[G2UncompressedSize:]
		ep2.Add(one)
	}
}
