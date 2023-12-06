package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var aLimbs, bLimbs, expectedLimbs [5]int64
	for i := range aLimbs {
		aLimbs[i] = randomLimb()
		bLimbs[i] = randomLimb()
		expectedLimbs[i] = aLimbs[i] + bLimbs[i]
	}
	a := New(aLimbs)
	b := New(bLimbs)
	expected := New(expectedLimbs)

	result := Add(a, b)
	if result != expected {
		t.Errorf("Add(%v,%v) = %v, but we expect %v", a, b, result, expected)
	}
}

func randomLimb() int64 {
	return rand.Int63n(1 << 51)
}

func TestSub(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var aLimbs, bLimbs, expectedLimbs [5]int64
	for i := range aLimbs {
		aLimbs[i] = randomLimb()
		bLimbs[i] = randomLimb()
		expectedLimbs[i] = aLimbs[i] - bLimbs[i]
	}

	a := New(aLimbs)
	b := New(bLimbs)
	expected := New(expectedLimbs)

	result := Sub(a, b)
	if result != expected {
		t.Errorf("Sub(%v,%v) = %v, but we expect %v", a, b, result, expected)
	}
}

func TestReduce(t *testing.T) {
	var a = New([5]int64{0, 0, 0, 0, 0})
	var expected = New([5]int64{0, 0, 0, 0, 0})
	var result = Reduce(a)
	if result != expected {
		t.Errorf("Reduce(%d) gave %d, but we expected %d", a, result, expected)
	}

	// 4503599627370495 = 2^52-1
	a = New([5]int64{4503599627370495, 1, 1, 1, 1})
	expected = New([5]int64{2251799813685247, 2, 1, 1, 1})
	result = Reduce(a)

	if result != expected {
		t.Errorf("Reduce(%d) gave %d, but we expteced %d", a, result, expected)
	}

	// 4503599627370495 = 2^52-1
	// 115792089237316221134579693152543734506205109808361273083374000351476182941695
	// fully reduced = 25711008708143855826652935125142720709043916416343563053301797
	a = New([5]int64{4503599627370495, 4503599627370495, 4503599627370495, 4503599627370495, 4503599627370495})
	expected = New([5]int64{37, 1, 1, 1, 1})
	// this needs 2 rounds of reducing
	result = Reduce(a)
	result = Reduce(result)

	if result != expected {
		t.Errorf("Reduce(%d) gave %d, but we expteced %d", a, result, expected)
	}

	a = New([5]int64{9223372036854775807, 0, 0, 0, 0})
	expected = New([5]int64{2251799813685247, 4095, 0, 0, 0})
	result = Reduce(a)
	if result != expected {
		t.Errorf("Reduce(%d) gave %d, but we expected %d", a, result, expected)
	}
}

/*
Check test values with PariGP. Below tests are with testvalues smaller than p.

// Convert radix-51 representation to number
radix51_to_num(coeffs) = sum(i=1, length(coeffs), coeffs[i] * 2^(51*(i-1)));

// Convert number to radix-51 representation

	num_to_radix51(n) = {
			local(radix, coeffs, i, rem);
			radix = 2^51;
			coeffs = vector(5, i, 0);
			rem = n;

			for(i = 1, 5,
					coeffs[i] = rem % radix;
					rem = (rem - coeffs[i]) / radix;
			);

			coeffs;
	};

// Multiply two numbers in radix-51 representation

	mult_radix51(a, b) = {
		local(a_num, b_num, product);
		a_num = radix51_to_num(a);
		b_num = radix51_to_num(b);
		product = Mod(a_num * b_num, 2^255 - 19);
		num_to_radix51(lift(product));
	};
*/
func TestMult1(t *testing.T) {
	// 2^0 + 2^51 + 2^102 + 2^153 + 2^204
	// 25711008708143855826652935125142720709043916416343563053301761
	a := New([5]int64{1, 1, 1, 1, 1})
	// 2* 2^0 + 2* 2^51 + 2* 2^102 + 2* 2^153 + 2* 2^204
	// 51422017416287711653305870250285441418087832832687126106603522
	b := New([5]int64{2, 2, 2, 2, 2})
	// 257110087081438969313864850568238035249664095928860339165200538
	result := Mult(a, b)
	// 154* 2^0 + 118* 2^51 + 82* 2^102 + 46* 2^153 + 10* 2^204
	// = 257110087081438969313864850568238035249664095928860339165200538
	// (154, 118, 82, 46, 10)
	expected := New([5]int64{154, 118, 82, 46, 10})
	if result != expected {
		t.Errorf("Mul(%d, %d) gave %d, but we expected %d", a, b, result, expected)
	}
}

func TestMult2(t *testing.T) {
	// 33645039376242384704430134258477280882292434460591223735468431356427724771976
	a := New([5]int64{1621387689972360, 0, 0, 0, 0})
	// 40366082900617430235565328797369117188960636571717375586789175107711011369464
	b := New([5]int64{1690142389023224, 0, 0, 0, 0})
	result := Mult(a, b)
	// = 2740376063862730982079958088640
	// [2099909120228288, 1216971440892824, 0, 0, 0]
	expected := New([5]int64{2099909120228288, 1216971440892824, 0, 0, 0})
	if result != expected {
		t.Errorf("Mul(%d, %d) gave %d, but we expected %d", a, b, result, expected)
	}
}

func TestMult3(t *testing.T) {
	// 9242459468621171906146226982473256942747135744698040831562376
	a := New([5]int64{1621387689972360, 922524701973052, 1829966140650555, 809465266247700, 0})
	// 22276686632890467149569150358968737131025025870088354267640312
	b := New([5]int64{1690142389023224, 1604293222359650, 2195352116801794, 1951017923057161, 0})
	result := Mult(a, b)
	// = 56377071098014759228856083824018241963500121865567763703020913802992311544142
	// [1022010010052942, 336831545554675, 2202455391182773, 820524535937719, 2192721092275060]
	expected := New([5]int64{1022010010052942, 336831545554675, 2202455391182773, 820524535937719, 2192721092275060})
	if result != expected {
		t.Errorf("Mul(%d, %d) gave %d, but we expected %d", a, b, result, expected)
	}
}

func TestMult4(t *testing.T) {
	// 57896044618658097711785492504343953926634992332820282019728792003956564819948
	a := New([5]int64{2251799813685228, 2251799813685247, 2251799813685247, 2251799813685247, 2251799813685247})
	// 22276686632890467149569150358968737131025025870088354267640312
	b := New([5]int64{1690142389023224, 1604293222359650, 2195352116801794, 1951017923057161, 0})
	result := Mult(a, b)
	// = 57896044618658075435098859613876804357484633364083150994702921915602297179637
	// [561657424662005, 647506591325597, 56447696883453, 300781890628086, 2251799813685247]
	expected := New([5]int64{561657424662005, 647506591325597, 56447696883453, 300781890628086, 2251799813685247})
	if result != expected {
		t.Errorf("Mul(%d, %d) gave %d, but we expected %d", a, b, result, expected)
	}
}

func TestMult5(t *testing.T) {
	// 57896044618658097711785492504343953926634992332820282019728792003956564819948
	a := New([5]int64{2251799813685228, 2251799813685247, 2251799813685247, 2251799813685247, 2251799813685247})
	// 57896044618658094277463417250966694824391873842955772689506761567625467638264
	b := New([5]int64{1690142389023224, 1604293222359650, 2195352116801794, 1951017923057161, 2251799813685247})
	result := Mult(a, b)
	// = 3434322075253377259102243118489864509330222030436331097181685
	// [561657424662005, 647506591325597, 56447696883453, 300781890628086, 0]
	expected := New([5]int64{561657424662005, 647506591325597, 56447696883453, 300781890628086, 0})
	if result != expected {
		t.Errorf("Mul(%d, %d) gave %d, but we expected %d", a, b, result, expected)
	}
}
