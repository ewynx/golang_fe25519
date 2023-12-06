package main

import (
	"fmt"
	"math/big"
)

type BigInt256 struct {
	l [5]int64
}

func New(l [5]int64) BigInt256 {
	return BigInt256{l: l}
}

func (a BigInt256) String() string {
	return fmt.Sprintf("(%d, %d, %d, %d, %d)", a.l[0], a.l[1], a.l[2], a.l[3], a.l[4])
}

// Add separate limbs
// expects limbs max 2^63-1
// does not reduce
func Add(a, b BigInt256) BigInt256 {
	var res BigInt256
	for i := 0; i < 5; i++ {
		res.l[i] = a.l[i] + b.l[i]
	}
	return res
}

// Subtract separate limbs
// expects limbs max 2^63-1
// does not reduce
func Sub(a, b BigInt256) BigInt256 {
	var res BigInt256
	for i := 0; i < 5; i++ {
		res.l[i] = a.l[i] - b.l[i]
	}
	return res
}

func Reduce(a BigInt256) BigInt256 {
	var carry int64
	// Carry
	for i := 0; i < 4; i++ {
		carry = a.l[i] >> 51
		a.l[i+1] += carry
		carry <<= 51
		a.l[i] -= carry
	}
	carry = a.l[4] >> 51
	a.l[0] += 19 * carry
	carry <<= 51
	a.l[4] -= carry
	return a
}

// Multiplies and reduces
// expects limbs max 2^52-1
func Mult(a BigInt256, b BigInt256) BigInt256 {
	var a_big, b_big [5]big.Int
	for i := 0; i < 5; i++ {
		a_big[i].SetInt64(int64(a.l[i]))
		b_big[i].SetInt64(int64(b.l[i]))
	}
	// Resulting coefficients have max 107 bits, thus won't fit in a u64
	var res [9]big.Int

	// (1) Set each limb separately.

	// r[0] = (int128) a[0]*a[0];
	res[0].Mul(&a_big[0], &b_big[0])

	// r[1] = (int128) a[0]*a[1] + (int128) a[1]*a[0];
	var temp [6]big.Int
	// temp[0] = a[0]*b[1]
	temp[0].Mul(&a_big[0], &b_big[1])
	// temp[1] = a[1]*b[0]
	temp[1].Mul(&a_big[1], &b_big[0])
	res[1].Add(&temp[0], &temp[1])

	// r[2] = (int128) a[0]*a[2] + (int128) a[1]*a[1] + (int128) a[2]*a[0];
	// temp[0] = a[0]*b[2]
	temp[0].Mul(&a_big[0], &b_big[2])
	// temp[1] = a[1]*b[1]
	temp[1].Mul(&a_big[1], &b_big[1])
	// temp[2] = a[2]*b[0]
	temp[2].Mul(&a_big[2], &b_big[0])
	// temp[3] = temp[0] + temp[1]
	temp[3].Add(&temp[0], &temp[1])
	res[2].Add(&temp[3], &temp[2])

	// r[3] = (int128) a[0]*a[3] + (int128) a[1]*a[2] + (int128) a[2]*a[1] + (int128) a[3]*a[0];
	// temp[0] = a[0]*b[3]
	temp[0].Mul(&a_big[0], &b_big[3])
	// temp[1] = a[1]*b[2]
	temp[1].Mul(&a_big[1], &b_big[2])
	// temp[2] = a[2]*b[1]
	temp[2].Mul(&a_big[2], &b_big[1])
	// temp[3] = a[3]*b[0]
	temp[3].Mul(&a_big[3], &b_big[0])
	// temp[4] = temp[0] + temp[1]
	// temp[0] = temp[2] + temp[3]
	temp[4].Add(&temp[0], &temp[1])
	temp[0].Add(&temp[2], &temp[3])
	res[3].Add(&temp[4], &temp[0])

	// r[4] = (int128) a[0]*a[4] + (int128) a[1]*a[3] + (int128) a[2]*a[2] + (int128) a[3]*a[1] + (int128) a[4]*a[0];
	// temp[0] = a[0]*b[4]
	temp[0].Mul(&a_big[0], &b_big[4])
	// temp[1] = a[1]*b[3]
	temp[1].Mul(&a_big[1], &b_big[3])
	// temp[2] = a[2]*b[2]
	temp[2].Mul(&a_big[2], &b_big[2])
	// temp[3] = a[3]*b[1]
	temp[3].Mul(&a_big[3], &b_big[1])
	// temp[4] = a[4]*b[0]
	temp[4].Mul(&a_big[4], &b_big[0])
	// temp[5] = temp[0] + temp[1]
	// then we can reuse temp 0 and 1
	// temp[0] = temp[2] + temp[3]
	temp[5].Add(&temp[0], &temp[1])
	temp[0].Add(&temp[2], &temp[3])
	temp[1].Add(&temp[5], &temp[0])
	res[4].Add(&temp[1], &temp[4])

	// r[5] = (int128) a[1]*a[4] + (int128) a[2]*a[3] + (int128) a[3]*a[2] + (int128) a[4]*a[1];
	// temp[0] = a[0]*b[3]
	temp[0].Mul(&a_big[1], &b_big[4])
	// temp[1] = a[1]*b[2]
	temp[1].Mul(&a_big[2], &b_big[3])
	// temp[2] = a[2]*b[1]
	temp[2].Mul(&a_big[3], &b_big[2])
	// temp[3] = a[3]*b[0]
	temp[3].Mul(&a_big[4], &b_big[1])
	// temp[4] = temp[0] + temp[1]
	// temp[0] = temp[2] + temp[3]
	temp[4].Add(&temp[0], &temp[1])
	temp[0].Add(&temp[2], &temp[3])
	res[5].Add(&temp[4], &temp[0])

	// r[6] = (int128) a[2]*b[4] + (int128) a[3]*b[3] + (int128) a[4]*b[2];
	temp[0].Mul(&a_big[2], &b_big[4])
	temp[1].Mul(&a_big[3], &b_big[3])
	temp[2].Mul(&a_big[4], &b_big[2])
	temp[3].Add(&temp[0], &temp[1])
	res[6].Add(&temp[3], &temp[2])

	// r[7] = (int128) a[3]*a[4] + (int128) a[4]*a[3];
	temp[0].Mul(&a_big[3], &b_big[4])
	temp[1].Mul(&a_big[4], &b_big[3])
	res[7].Add(&temp[0], &temp[1])

	// r[8] = (int128) a[4]*a[4];
	res[8].Mul(&a_big[4], &b_big[4])

	// (2) Reduce from 9 coefficients to 5
	// r0=r0+19*r5
	// r1=r1+19*r6
	// r2=r2+19*r7
	// r3=r3+19*r8
	temp[0].Mul(big.NewInt(19), &res[5])
	res[0].Add(&res[0], &temp[0])

	temp[0].Mul(big.NewInt(19), &res[6])
	res[1].Add(&res[1], &temp[0])

	temp[0].Mul(big.NewInt(19), &res[7])
	res[2].Add(&res[2], &temp[0])

	temp[0].Mul(big.NewInt(19), &res[8])
	res[3].Add(&res[3], &temp[0])

	// (3) Carry 1 round on the 5 coeffs to obtain 64 bits integers
	var carry big.Int
	for i := 0; i < 4; i++ {
		carry.Rsh(&res[i], 51)
		res[i+1].Add(&res[i+1], &carry)
		carry.Lsh(&carry, 51)
		res[i].Sub(&res[i], &carry)
	}
	carry.Rsh(&res[4], 51)
	temp[0].Mul(&carry, big.NewInt(19))
	res[0].Add(&res[0], &temp[0])
	carry.Lsh(&carry, 51)
	res[4].Sub(&res[4], &carry)

	// (4) The coefficients fit in 64 bits
	var res_64 [5]int64
	for i := 0; i < 5; i++ {
		res_64[i] = res[i].Int64()
	}
	var res_bigint BigInt256 = BigInt256{l: res_64}
	res_bigint = Reduce(res_bigint)

	return Reduce(res_bigint)
}
