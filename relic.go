package bls12

// #cgo CFLAGS: -std=c99 -O2 -I. -DARCH=X64 -Irelic/include -Irelic/src -Irelic/include/low -Wno-unused -DALLOC=AUTO -Wno-discarded-qualifiers
// #include "relic_core.h"
import "C"
import (
	"math/big"
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
