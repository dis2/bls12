// Fr scalar field

package bls12

// #include "relic_bn.h"
import "C"
import "math/big"
import "fmt"

const (
	ScalarSize = 32
	BigScalarSize = 32
)

// Represents a scalar.
type Scalar = C.bn_st

// Convert to scalar from big.Int
func (s *Scalar) FromInt(n *big.Int) *Scalar {
	buf := n.Bytes()
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


// Unmarshal scalar from a byte buffer. Only small 256bit scalars are
// to be marshalled.
func (s *Scalar) Unmarshal(buf []byte) {
	C.bn_read_bin(s, (*C.uint8_t)(&buf[0]), C.int(len(buf)))
}

// Marshal scalar to byte buffer.
func (s *Scalar) Marshal() []byte {
	var buf [ScalarSize]byte
	C.bn_write_bin((*C.uint8_t)(&buf[0]), C.int(len(buf)), s)
	return buf[:]
}

// up to 512bit scalar
func ScalarConst(s string) (bn Scalar) {
	if len(s)%2 != 0 {
		panic("bad const padding for "+s)
	}
	bi := new(big.Int)
	bi.SetString(s, 16)
	bn.FromInt(bi)
	return
}

