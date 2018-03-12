Bilinear pairing BLS12-381. 

## RTFM

* https://github.com/ebfull/pairing/tree/master/src/bls12_381#bls12-381 for specification.
* https://github.com/dis2/bls12 for API documentation.
* https://github.com/dis2/blsig for example applicatoin.

## Performance

(2.6ghz Haswell)

```
BenchmarkBaseMultG1              10000            204899 ns/op
BenchmarkMultG1                   3000            543081 ns/op
BenchmarkMultG2                   3000            593233 ns/op
BenchmarkPair                     1000           1930860 ns/op
```

Ie about ~500 pairings/s. This is not far off from BN256.

## Building gotchas

This is wrapper for RELIC. It is self-contained in that it links to amalgam
of c files from git submodule of relic, no libraries involved. Meaning import
should just straight work.
