package bls12

// #include "relic_bn.h"
import "C"
import "math/big"

// Represents a scalar.
type Scalar struct {
	st C.bn_st
}

// Convert to scalar big.Int
func (s *Scalar) FromInt(n *big.Int) *Scalar {
	return s.ScalaromBytes(n.Bytes())
}

// Convert to scalar raw bytes
func (s *Scalar) FromBytes(s []byte) *Scalar {
	C.bn_read_bin(&s.st, (*C.uint8_t)(&s[0]), C.int(len(s)))
	return s
}

