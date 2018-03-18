package bls12
import "math/big"

type Fq [6]uint64

func parseInt(s string) (i *big.Int) {
	i = new(big.Int)
	i.SetString(s, 10)
	return
}

func parseFq(s string) (a Fq) {
	a.SetWords(parseInt(s).Bits())
	return
}

// Get big.Int words representing (mont) the element.
func (a *Fq) SetWords(b []big.Word) {
	n := uint(len(b))
	if uint64(^uint(0))>>63 == 1 {
		for i:=uint(0); i < n; i++ {
			a[i] = uint64(b[i])
		}
	} else {
		for i:=uint(0); i < n; i++ {
			a[i/2] |= uint64(b[i]) << (32*(i%2))
		}
	}
}

// Set big.Int words representing (mont) the element.
func (a *Fq) GetWords() []big.Word {
	if uint64(^uint(0))>>63 == 1 {
		var b[6]big.Word
		for i:=0; i < 6; i++ {
			b[i] = big.Word(a[i])
		}
		return b[:]
	} else {
		var b[12]big.Word
		for i:=uint(0); i < 12; i++ {
			b[i] = big.Word(a[i]) >> (32*(i%2))
		}
		return b[:]
	}
}

func (a *Fq) RawInt() (res *big.Int) {
	res = new(big.Int)
	res.SetBits(a.GetWords())
	return
}

func (e *Fq) FromInt(i *big.Int) *Fq {
	e.SetWords(i.Bits())
	return e.Mul(e, &R2)
}

func (e *Fq) ToInt() (res *big.Int) {
	res = new(big.Int)
	var t Fq
	res.SetBits(t.Mul(e, &Fq{1}).GetWords())
	return
}

