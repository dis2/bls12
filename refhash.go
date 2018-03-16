package bls12

import "golang.org/x/crypto/blake2b"
import "encoding/binary"
import "bytes"

// Random input sampling used by Rust HTP
// roughly Fq::hash in rust
func (e *Fq) HashRef(key, nonce []byte) bool {
	// There is no bound, so we introduce arbitrary one
	for i := uint32(1); i < (1<<16); i++ {
		h, err := blake2b.New384(key)
		if err != nil {
			return false
		}
		h.Write(key)
		h.Write(nonce)

		// NativeEndian::write_u32(&mut count_u8, count); ??
		binary.Write(h, binary.LittleEndian, i)

		sample := h.Sum(nil)[:]
		// shave off 3 bits
		sample[0] &= 0x1f
		if bytes.Compare(sample, QBytes) < 0 {
			e.Unmarshal(sample)
			return true
		}
	}
	return false
}

func (p *G1) HashToPoint(key, non []byte) bool {
	var np[2] G1
	var t[2] Fq
	var msg [64]byte
	var nonce [33]byte

	// Rust silently truncates keys to 64 bytes, so do we.
	copy(msg[:],key)

	if non != nil {
		copy(nonce[:], non)
	}

	for i := 0; i < 2; i++ {
		nonce[32] = byte(-i)
		if !t[i].HashRef(msg[:], nonce[:]) {
			return false
		}
		FouqueMapXtoY(&t[i],&np[i].X,&np[i].Y)
		if !t[i].IsResidue() {
			np[i].Y.Neg(nil)
		}
		np[i].SetNormalized()
		np[i].ScaleByCofactor()
	}
	*p = np[0]
	p.Add(&np[1])
	return true
}

func (p *G2) HashToPoint(key, non []byte) bool {
	var np[2] G2
	var t[2] Fq2
	var msg [64]byte
	var nonce [33]byte

	copy(msg[:], key)
	if non != nil {
		copy(nonce[:], non)
	}

	for i := 0; i < 2; i++ {
		// let t1 = Self::Base::hash(seed, nonce); ->>
		//   let half = nonce.len() / 2;
		//   Fq2 { c0: Fq::hash(k, &nonce[ .. half]), c1: Fq::hash(k, &nonce[half ..]) }
		nonce[32] = byte(-i)
		if !t[i].C[0].HashRef(msg[:], nonce[:len(nonce)/2]) {
			return false
		}
		if !t[i].C[1].HashRef(msg[:], nonce[len(nonce)/2:]) {
			return false
		}
		FouqueMapXtoY(&t[i],&np[i].X,&np[i].Y)
		if !t[i].IsResidue() {
			np[i].Y.Neg(nil)
		}
		np[i].SetNormalized()
		np[i].ScaleByCofactor()
	}
	*p = np[0]
	p.Add(&np[1])
	return true
}


