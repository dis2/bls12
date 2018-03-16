package bls12

import (
	"testing"
)

func TestMapFouqueFq(t *testing.T) {
	for i := 1; i < 10000; i += 2 {
		var t,x,y Fq
		t.SetInt64(int64(i) + 124124)
		FouqueMapXtoY(&t,&x,&y)
	}
}

func TestMapFouqueFq2(t *testing.T) {
	for i := 1; i < 10000; i += 2 {
		var t,x,y Fq2
		t.SetInt64(int64(i) + 124124)
		FouqueMapXtoY(&t,&x,&y)
	}
}
