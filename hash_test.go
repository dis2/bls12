package bls12

import (
	"testing"
	"crypto/rand"
)

func TestMapFouqueFq(te *testing.T) {
	for i := 1; i < 10000; i++ {
		var t,x,y Fq
		var buf [48]byte
		rand.Read(buf[:])
		buf[0] &= 0xf
		t.Unmarshal(buf[:])
		FouqueMapXtoY(&t,&x,&y)
		// try match parity of y with t
		if t.Limbs[0]&1 != y.Limbs[0]&1 {
			y.Neg(&y)
			if t.Limbs[0]&1 != y.Limbs[0]&1 {
				te.Fatal("odd/even parity invariant failed")
			}
		}

	}
}

func TestMapFouqueFq2(te *testing.T) {
	for i := 1; i < 10000; i++ {
		var t,x,y Fq2
		var buf [96]byte
		rand.Read(buf[:])
		buf[0] &= 0xf
		buf[48] &= 0xf
		t.Unmarshal(buf[:])
		FouqueMapXtoY(&t,&x,&y)
		// try match parity of y with t
		if t.C[0].Limbs[0]&1 != y.C[0].Limbs[0]&1 {
			y.Neg(&y)
			if t.C[0].Limbs[0]&1 != y.C[0].Limbs[0]&1 {
				te.Fatal("odd/even parity invariant failed")
			}
		}
	}
}

