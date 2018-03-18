package bls12

func reduce(a *Fq) {
	var b Fq

	var carry uint64
	for i, pi := range Q {
		ai := a[i]
		bi := ai - pi - carry
		b[i] = bi
		carry = (pi&^ai | (pi|^ai)&bi) >> 63
	}

	carry = -carry
	ncarry := ^carry
	for i := 0; i < 6; i++ {
		a[i] = (a[i] & carry) | (b[i] & ncarry)
	}
	return
}

func (aa *Fq) Mul(bb, cc *Fq) *Fq {
	if bb == nil {
		bb = aa
	}
	var abuf[12]uint64
	var carry, carry2 uint64

	for i := 0; i < 6; i++ {
		carry = 0
		b := bb[i]
		bh, bl := b>>32, b&mask32
		for j := 0; j < 6; j++ {
			c := cc[j]
			a := abuf[j+i]

			ch, cl := c>>32, c&mask32
			ah, al := a>>32, a&mask32

			w := bl * cl
			x := bh * cl
			y := bl * ch

			r0 := (w & mask32) + al + (carry & mask32)

			z := bh * ch

			r1 := (r0 >> 32) + (w >> 32) + (x & mask32) + (y & mask32) + ah + (carry >> 32)
			r2 := (r1 >> 32) + (x >> 32) + (y >> 32) + (z & mask32)
			carry = (((r2 >> 32) + (z >> 32)) << 32) | (r2 & mask32)
			abuf[i+j] = (r1 << 32) | (r0 & mask32)
		}
		abuf[i+6] = carry
	}

	for i := 0; i < 6; i++ {
		b := qInv64 * abuf[i]
		carry = 0
		bh, bl := b>>32, b&mask32
		for j := 0; j < 6; j++ {
			c := Q[j]
			a := abuf[i+j]

			ch, cl := c>>32, c&mask32
			ah, al := a>>32, a&mask32

			w := bl * cl
			x := bh * cl
			y := bl * ch

			r0 := (w & mask32) + al + (carry & mask32)

			z := bh * ch

			r1 := (r0 >> 32) + (w >> 32) + (x & mask32) + (y & mask32) + ah + (carry >> 32)
			r2 := (r1 >> 32) + (x >> 32) + (y >> 32) + (z & mask32)
			carry = (((r2 >> 32) + (z >> 32)) << 32) | (r2 & mask32)

			if j > 0 {
				abuf[i+j] = (r1 << 32) | (r0 & mask32)
			}
		}
		a := abuf[i+6]
		l := (a & mask32) + (carry2 & mask32) + (carry & mask32)
		h := (l >> 32) + (a >> 32) + (carry2 >> 32) + (carry >> 32)
		carry2 = h >> 32
		v := (h << 32) | (l & mask32)
		abuf[i+6] = v
	}
	aa[0] = abuf[6]
	aa[1] = abuf[7]
	aa[2] = abuf[8]
	aa[3] = abuf[9]
	aa[4] = abuf[10]
	aa[5] = abuf[11]

	reduce(aa)
	return aa
}

func (a *Fq) Square(b *Fq) *Fq {
	return a.Mul(b, b)
}

func (a *Fq) Add(b, c *Fq) *Fq {
	if b == nil {
		b = a
	}
	var carry uint64
	for i, ai := range a {
		bi := c[i]
		ci := ai + bi + carry
		a[i] = ci
		carry = (ai&bi | (ai|bi)&^ci) >> 63
	}
	reduce(a)
	return a
}

func (a *Fq) Sub(b, c *Fq) *Fq {
	if b == nil {
		b = a
	}
	var t Fq
	a.Add(b,t.Neg(c))
	reduce(a)
	return a
}

func (a *Fq) Neg(b *Fq) *Fq {
	var carry uint64
	for i, pi := range Q {
		ai := b[i]
		ci := pi - ai - carry
		a[i] = ci
		carry = (ai&^pi | (ai|^pi)&ci) >> 63
	}
	return a
}

