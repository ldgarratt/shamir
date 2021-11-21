package main

import (
    "testing"
    "math/big"
)

func TestFormat(t *testing.T) {
    p := polynomial{[]*big.Int{big.NewInt(2), big.NewInt(4), big.NewInt(3), big.NewInt(0), big.NewInt(2)}}
    result := p.format()
    if result != "2x^4 + 3x^2 + 4x + 2" {
        t.Errorf("expecting 2x^4 + 3x^2 + 4x + 2, got %s", result)
    }

    p = polynomial{[]*big.Int{big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(7)}}
    result = p.format()
    if result != "7x^4 + 1" {
        t.Errorf("expecting 7x^4 + 1, got %s", result)
    }

    p = polynomial{[]*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(5), big.NewInt(0), big.NewInt(0)}}
    result = p.format()
    if result != "5x^2 + x + 1" {
        t.Errorf("expecting 5x^2 + x + 1, got %s", result)
    }

    p = polynomial{[]*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(5), big.NewInt(0), big.NewInt(0)}}
    result = p.format()
    if result != "5x^2 + x" {
        t.Errorf("expecting 5x^2 + x, got %s", result)
    }

    p = polynomial{[]*big.Int{big.NewInt(2), big.NewInt(6)}}
    result = p.format()
    if result != "6x + 2" {
        t.Errorf("expecting  6x + 2, got %s", result)
    }
}

func TestEvaluatePolynomial(t *testing.T) {
    p := polynomial{[]*big.Int{big.NewInt(2), big.NewInt(4), big.NewInt(3), big.NewInt(0), big.NewInt(2)}}
    x := big.NewInt(3)
    modulus := big.NewInt(17)
    result := evaluatePolynomial(x, modulus, p)
    if result.Cmp(big.NewInt(16)) != 0 {
        t.Errorf("Expecting 3, got: %s", result.String())
    }

    p = polynomial{[]*big.Int{big.NewInt(0), big.NewInt(12), big.NewInt(-9), big.NewInt(0), big.NewInt(0)}}
    x = big.NewInt(4)
    modulus = big.NewInt(17)
    result = evaluatePolynomial(x, modulus, p)
    if result.Cmp(big.NewInt(6)) != 0 {
        t.Errorf("Expecting 6, got: %s", result.String())
    }

    p = polynomial{[]*big.Int{big.NewInt(-3), big.NewInt(12), big.NewInt(-9), big.NewInt(0), big.NewInt(1)}}
    x = big.NewInt(-5)
    modulus = big.NewInt(20)
    result = evaluatePolynomial(x, modulus, p)
    if result.Cmp(big.NewInt(17)) != 0 {
        t.Errorf("Expecting 17, got: %s", result.String())
    }
}

func Test_shamirSplitSecretwithFixedPolynomial(t *testing.T) {
    // Taken from https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing
    secret := big.NewInt(1234)
    modulus := big.NewInt(1613)
    n := 6
    threshold := 3
    poly := polynomial{[]*big.Int{secret, big.NewInt(166), big.NewInt(94)}}

    result := _shamirSplitSecretwithFixedPolynomial(secret, modulus, poly, n, threshold)
    expected := []*big.Int{big.NewInt(1494), big.NewInt(329), big.NewInt(965), big.NewInt(176), big.NewInt(1188), big.NewInt(775)}

    if len(expected) != len(result) {
        t.Error("Expected share length is %i, result length is %i", len(expected), len(result))
    }

    for i := 0; i < len(expected); i++ {
        if result[i].Cmp(expected[i]) != 0 {
            t.Errorf("Expecting %s, got: %s", expected[i], result[i])
        }
    }
}

func TestLagrange(t *testing.T) {
    // Taken from https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing
    modulus := big.NewInt(1399)
    points := map[int]big.Int{
        2 : *big.NewInt(1942),
        4 : *big.NewInt(3402),
        5 : *big.NewInt(4414),
    }
    result := lagrange(points, modulus)
    expected := big.NewInt(1234)

    if result.Cmp(expected) != 0 {
        t.Errorf("Expecting %s, got: %s", expected, result)
    }
}
