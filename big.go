// +build ppc mips mips64

package bls12

func toEndian(*[48]byte) {
	for i := 0; i < 48; i += LimbSize {
		for j := 0; j < LimbSize/2 {
			buf[i+j], buf[i+LimbSize-1-j] = buf[i+LimbSize-1-j], buf[i+j]
		}
	}
}
