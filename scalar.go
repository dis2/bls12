package bls12

// #include "relic_bn.h"
import "C"
import "math/big"
import "fmt"

const (
	ScalarSize = 32
)

// Represents a scalar.
type Scalar struct {
	st C.bn_st
}

// Convert to scalar from big.Int
func (s *Scalar) FromInt(n *big.Int) *Scalar {
	buf := n.Bytes()
	if len(buf) > ScalarSize {
		return nil
	}
	s.Unmarshal(buf)
	return s
}

// Convert scalar to big.Int
func (s *Scalar) ToInt() (n *big.Int) {
	return new(big.Int).SetBytes(s.Marshal())
}

func (s *Scalar) String() string {
	return fmt.Sprintf("bls12.Scalar(%x)", s.Marshal())
}

// Unmarshal scalar from a byte buffer.
func (s *Scalar) Unmarshal(buf []byte) []byte {
	nb := len(buf)
	if nb > ScalarSize {
		nb = ScalarSize
	}
	C.bn_read_bin(&s.st, (*C.uint8_t)(&buf[0]), C.int(nb))
	return buf[nb:]
}

// Marshal scalar to byte buffer.
func (s *Scalar) Marshal() []byte {
	var buf [ScalarSize]byte
	C.bn_write_bin((*C.uint8_t)(&buf[0]), C.int(len(buf)), &s.st)
	return buf[:]
}
