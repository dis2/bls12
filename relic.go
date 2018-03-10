package bls12

// #cgo CFLAGS: -std=c99  -Wall -O3 -funroll-loops -fomit-frame-pointer -finline-small-functions -march=native -mtune=native
// -DBN_PRECI=384 -DALLOC=DYNAMIC -DALIGN=8 -DARITH=EASY
// #include "relic_core.h"
// bn_t _bn_new() { bn_t t; bn_new(t); return t; };
// void _bn_free(bn_t t) { bn_free(t); };
import "C"
import (
	"math/big"
	"os"
)

var r *big.Int

func init() {
	C.core_init()
	C.ep_param_set_any_pairf()
	r = (&big.Int{}).SetBytes(ScalarOrder())
}

func checkError() {
	// nop for now
}

func ScalarOrder() []byte {
	var r C.bn_st
	C.ep2_curve_get_ord(&r)
	checkError()
	buf := make([]byte, 48)
	C.bn_write_bin((*C.uint8_t)(&buf[0]), C.int(len(buf)), &r)
	checkError()
	return buf
}

func IsScalar(s []byte) bool {
	bn := (&big.Int{}).SetBytes(s)
	return bn.Cmp(r) < 0
}
