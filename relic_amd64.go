package bls12

// #include "relic_core.h"
// #include "relic_fp.h"
// #cgo CFLAGS: -std=c99 -O2 -I. -DARCH=X64 -DWORD=64 -Irelic/include  -Irelic/src -Irelic/include/low -Wno-unused -DALLOC=AUTO -Wno-discarded-qualifiers -Wno-incompatible-pointer-types
import "C"

const NLimbs = 6
type Limb = C.dig_t
