// +build cgo

package bls12

// #include "relic_core.h"
// #include "relic_ep.h"
// void _fp_neg(fp_t r, const fp_t p) { fp_neg(r, p); }
// void _bn_new(bn_t bn) { bn_new(bn); }
import "C"

var initDone = false

func initPending() {
	if initDone {
		return
	}
	initDone = true
	C.core_init()
	C.ep_param_set_any_pairf()
}


