// For generic platform.

// +build !amd64

package bls12

// #cgo CFLAGS: -std=c99 -O2 -I. -DARCH=-1 -DWORD=32 -Irelic/include -Irelic/src -Irelic/include/low -Wno-unused -DALLOC=AUTO -Wno-discarded-qualifiers
import "C"

