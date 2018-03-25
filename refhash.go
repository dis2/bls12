package bls12

import "golang.org/x/crypto/blake2b"
import "math/big"
var qint = new(big.Int).SetBytes(QBytes)

func hashRef(r *Fq, msg, suff []byte) {
	h, _ := blake2b.New512(nil)
	h.Write(msg)
	h.Write(suff)
	i := new(big.Int).SetBytes(h.Sum(nil))
	i = i.Mod(i, qint)
	b := i.Bytes()
	b = append(make([]byte, 48-len(b)), b...)
	r.Unmarshal(b)
}

// Hash to point
func (p *G1) HashToPoint(msg []byte) G {
	var np [2]G1

	for i := byte(0); i < 2; i++ {
		var t Fq
		// G1_i
		hashRef(&t, msg, []byte{0x47,0x31,0x5f,0x30+i})
		FouqueMapXtoY(&t, &np[i].X, &np[i].Y)
		np[i].Y.CopyParity(&t)
		np[i].SetNormalized()
	}
	*p = np[0]
	p.Add(&np[1])
	p.ScaleByCofactor()
	p.Normalize()

	return p
}

// Hash to point
func (p *G2) HashToPoint(msg []byte) G {
	var np [2]G2

	for i := byte(0); i < 2; i++ {
		var t Fq2
		// G2_i_c0
		hashRef(&t.C[0], msg, []byte{0x47,0x32,0x5f,0x30+i,0x5f,0x63,0x30})
		// G2_i_c1
		hashRef(&t.C[1], msg, []byte{0x47,0x32,0x5f,0x30+i,0x5f,0x63,0x31})
		FouqueMapXtoY(&t, &np[i].X, &np[i].Y)
		np[i].Y.CopyParity(&t)
		np[i].SetNormalized()
	}
	*p = np[0]
	p.Add(&np[1])
	p.ScaleByCofactor()
	p.Normalize()
	return p
}


