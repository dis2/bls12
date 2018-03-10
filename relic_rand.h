// do not let relic use randomness
#define rand_bytes(a,b) { abort(); }
#define rand_seed(a,b) { }
#define rand_init()
#define rand_clean()


