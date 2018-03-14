package bls12

var (
	// X = -((2**63) + (2**62) + (2**60) + (2**57) + (2**48) + (2**16))
	X = hexConst("d201000000010000")
	// R = (X**4) - (X**2) + 1 is the ~256 bit base prime field, this is the order of G1 and G2
	R     = ScalarConst("73EDA753299D7D483339D80809A1D80553BDA402FFFE5BFEFFFFFFFF00000001")
	Order = R.ToInt()


	// Q = (((X - 1) ** 2) * ((X**4) - (X**2) + 1) // 3) + X is the ~384bit extended prime field
	Q = QConst("1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab")
	QMinus1 = QConst("1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaaa")
	// (Q-1)/2, used for legendre
	QMinus1Half = QConst("0d0088f51cbff34d258dd3db21a5d66bb23ba5c279c2895fb39869507b587b120f55ffff58a9ffffdcff7fffffffd555")
	// (Q+1)/4, used for tonelli-shanks
	QPlus1Quarter = ScalarConst("0680447a8e5ff9a692c6e9ed90d2eb35d91dd2e13ce144afd9cc34a83dac3d8907aaffffac54ffffee7fbfffffffeaab")

	// sqrt(-3), swenc const0
	QSqrtMinus3 = QConst("be32ce5fbeed9ca374d38c0ed41eefd5bb675277cdf12d11bc2fb026c41400045c03fffffffdfffd")

	// (sqrt(-3)-1) / 2, swenc const1
	QSqrtMinus3Minus1Half = QConst("5f19672fdf76ce51ba69c6076a0f77eaddb3a93be6f89688de17d813620a00022e01fffffffefffe") // SWENC_CONST1
	One = QConst("01")
	Zero = QConst("00")
	Four = QConst("04")
	Five = QConst("05")

	// G1 cofactor, ((X-1)**2) // 3
	G1_h = ScalarConst("396C8C005555E1568C00AAAB0000AAAB")
	// G2 cofactor,  ((X**8) - (4 * (X**7)) + (5 * (X**6)) - (4 * (X**4)) + (6 * (X**3)) - (4 * (X**2)) - (4*X) + 13) // 9
	G2_h = hexConst("5d543a95414e7f1091d50792876a202cd91de4547085abaa68a205b2e5a7ddfa628f1cb4d9e82ef21537e293a6691ae1616ec6e786f0c70cf1c38e31c7238e5")

	// G1 and G2 generators, see https://github.com/ebfull/pairing/tree/master/src/bls12_381#generators
	G1_x = QConst("17F1D3A73197D7942695638C4FA9AC0FC3688C4F9774B905A14E3A3F171BAC586C55E83FF97A1AEFFB3AF00ADB22C6BB")
	G1_y = QConst("08B3F481E3AAA0F1A09E30ED741D8AE4FCF5E095D5D00AF600DB18CB2C04B3EDD03CC744A2888AE40CAA232946C5E7E1")
	// G2 points are of the form c0 + c1 * u
	G2_x_c0 = QConst("024AA2B2F08F0A91260805272DC51051C6E47AD4FA403B02B4510B647AE3D1770BAC0326A805BBEFD48056C8C121BDB8")
	G2_x_c1 = QConst("13E02B6052719F607DACD3A088274F65596BD0D09920B61AB5DA61BBDC7F5049334CF11213945D57E5AC7D055D042B7E")
	G2_y_c0 = QConst("0CE5D527727D6E118CC9CDC6DA2E351AADFD9BAA8CBDD3A76D429A695160D12C923AC9CC3BACA289E193548608B82801")
	G2_y_c1 = QConst("0606C4A02EA734CC32ACD2B02BC28B99CB3E287E85A763AF267492AB572E99AB3F370D275CEC1DA1AAA9075FF05F79BE")
)
