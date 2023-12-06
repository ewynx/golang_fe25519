# Field arithmetic for p =2^255-19 using radix-2ˆ51

This repo contains a simple implementation of addition, subtraction and multiplication for field elements in a finite field modulo $2^{255}-19$. It follow the techniques discussed in [this presentation](https://cryptojedi.org/peter/data/pairing-20131122.pdf) by [Peter Schwabe](https://cryptojedi.org/peter/index.shtml).

In particular, it uses a radix- $2^{51}$, i.e every field element $a$ is represented as $(a_0,..,a_5)$ where $a = a_02^0+a_12^{51}+a_22^{102}+a_32^{153}+a_42^{204}$. For more details see [these slides](https://cryptojedi.org/peter/data/pairing-20131122.pdf).

The goal of this code was to get familiar with Golang. It is not intended to be a reference implementation and should be used solely for educational purposes.

## Build and run

Make sure to have [Golang](https://go.dev/doc/install) installed. 
The `main` function of this repo contains one of the arithmetic test cases, and can be easily adjusted to try out the functions. Run it as follows:

```
go run .
```

## Testing 
Run all tests (in `fp25519_test.go`):
```
go test
```

## Notes for arithmetic implementation

Notes made from [slides](https://cryptojedi.org/peter/data/pairing-20131122.pdf). These are added for reference. 

### 1. Representation (radix $2^{51}$)

Represent an integer as an object with 5 limbs of 64 bits. 
$(a_0,a_1,a_2,a_3,a_4)$ equals $\sum_{n=0}^{4} a_i2^{51\cdot i}$

This is in *reduced* form when all $a_i$ are in $[-(2^{52}-1),..,2^{52}-1]$ (notice this is the max value of 52 bits). 

### 2. Addition

When coefficients $a_i < 2^{63}-1$
```
res[0] = a[0] + b[0]
res[1] = a[1] + b[1]
res[2] = a[2] + b[2]
res[3] = a[3] + b[3]
res[4] = a[4] + b[4]
```
### 3. Subtraction

Use signed limbs and this can work fine. 

### 4. Carry & Reduce mod p

Carry for the first $4$ limbs, like this:
```
carry = a[0] >> 51
a[1] += carry
carry <<== 51
a[0] -= carry
```
Add the carry to the next limbs and make sure the original limb is reduces to 51 bits. 

Since $p = 2^{255}-19$ we carry and reduce last limb like this:

```
carry = a[4] >> 51
a[0] += 19*carry
carry <<== 51
a[4] -= carry
```

Because we reduced in the first few steps to 51 bits, and $19\cdot carry$ has an absolute value of at most 17 bits (since carry itself is max 12 bits), $a[0]+19\cdot carry$ is still in $[-(2^{52}-1),..,2^{52}-1]$.  

### 5. Multiplication

$A = \sum_{n=0}^{4} a_i2^{51\cdot i}$ and $B = \sum_{n=0}^{4} b_i2^{51\cdot i}$

```
r[0] = (int128) a[0]*b[0];
r[1] = (int128) a[0]*b[1] + (int128) a[1]*b[0];
r[2] = (int128) a[0]*b[2] + (int128) a[1]*b[1] + (int128) a[2]*b[0];
r[3] = (int128) a[0]*b[3] + (int128) a[1]*b[2] + 
(int128) a[2]*b[1] + (int128) a[3]*a[0];
r[4] = (int128) a[0]*b[4] + (int128) a[1]*b[3] + (int128) a[2]*b[2] + 
(int128) a[3]*b[1] + (int128) a[4]*a[0];
r[5] = (int128) a[1]*b[4] + (int128) a[2]*b[3] + 
(int128) a[3]*b[2] + (int128) a[4]*a[1];
r[6] = (int128) a[2]*b[4] + (int128) a[3]*b[3] + (int128) a[4]*b[2];
r[7] = (int128) a[3]*b[4] + (int128) a[4]*b[3];
r[8] = (int128) a[4]*b[4];
```

Multiplication gives $R = \sum_{n=0}^{8} r_i2^{51\cdot i}$ with $r_i$ up to 107 bits, which no longer fits in a 64 bit integer. Make sure to use a 128 bit integer. 

First, reduce from 9 coefficients to 5 with:

```
r0=r0+19*r5
r1=r1+19*r6
r2=r2+19*r7
r3=r3+19*r8
```

Then carry, as before. After round 1 we have signed 64-bit integers. Therefore, we need a second round of carries to obtain reduced coefficients. 
