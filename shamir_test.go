package main

import (
    "testing"
    "math/big"
    "os"
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

    result := _shamirSplitSecretWithFixedPolynomial(secret, modulus, poly, n, threshold)
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

func TestPairwiseJoinSlices(t *testing.T) {
    s1 := []*big.Int{big.NewInt(23), big.NewInt(345)}
    s2 := []*big.Int{big.NewInt(100), big.NewInt(99)}
    s3 := []*big.Int{big.NewInt(19), big.NewInt(50)}
    subsecret_shares := [][]*big.Int{s1, s2, s3}
    result := pairwiseJoinSlices(subsecret_shares)
    expected := []string{"23+100+19", "345+99+50"}

    if (result[0] != expected[0]) || (result[1] != expected[1])  {
        t.Errorf("Expecting %s, got: %s", expected, result)
    }

    s1 = []*big.Int{big.NewInt(321), big.NewInt(701183), big.NewInt(15263), big.NewInt(2574), big.NewInt(417)}
    s2 = []*big.Int{big.NewInt(117465), big.NewInt(599), big.NewInt(1207), big.NewInt(1752), big.NewInt(40624)}
    subsecret_shares = [][]*big.Int{s1, s2}

    result = pairwiseJoinSlices(subsecret_shares)
    expected = []string{"321+117465", "701183+588", "15263+1207", "2574+1752", "417+48624"}

    if (result[0] != expected[0])   {
        t.Errorf("Expecting %s, got: %s", expected, result)
    }

}


// TODO:
// Add end-to-end test
