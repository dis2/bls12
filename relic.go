// For generic platform.

// +build !amd64

package bls12

// #include "relic_core.h"
// #include "relic_fp.h"
// #cgo CFLAGS: -std=c99 -O2 -I. -DARCH=-1 -DWORD=32 -Irelic/include -Irelic/src -Irelic/include/low -Wno-unused -DALLOC=AUTO -Wno-discarded-qualifiers -Wno-incompatible-pointer-types
import "C"

const NLimbs = 12

type Limb = C.dig_t
