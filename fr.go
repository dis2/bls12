package bls12

// #include "relic_bn.h"
import "C"
import "math/big"

type Fr struct {
	st C.bn_st
}

func (fr *Fr) FromInt(n *big.Int) *Fr {
	return fr.FromBytes(n.Bytes())
}

func (fr *Fr) FromBytes(s []byte) *Fr {
	C.bn_read_bin(&fr.st, (*C.uint8_t)(&s[0]), C.int(len(s)))
	return fr
}

