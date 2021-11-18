package main

import (
    "testing"
    "math/big"
    "fmt"
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
}

func TestEvaluatePolynomial(t *testing.T) {
    p := polynomial{[]*big.Int{big.NewInt(2), big.NewInt(4), big.NewInt(3), big.NewInt(0), big.NewInt(2)}}
    x := big.NewInt(3)
    modulus := big.NewInt(17)
    result := evaluatePolynomial(x, modulus, p)
    fmt.Println(result)
    if result.Cmp(big.NewInt(16)) != 0 {
        t.Errorf("Expecting 3, got: %s", result.String())
    }

    p = polynomial{[]*big.Int{big.NewInt(0), big.NewInt(12), big.NewInt(-9), big.NewInt(0), big.NewInt(0)}}
    x = big.NewInt(4)
    modulus = big.NewInt(17)
    result = evaluatePolynomial(x, modulus, p)
    fmt.Println(result)
    if result.Cmp(big.NewInt(6)) != 0 {
        t.Errorf("Expecting 6, got: %s", result.String())
    }

    p = polynomial{[]*big.Int{big.NewInt(-3), big.NewInt(12), big.NewInt(-9), big.NewInt(0), big.NewInt(1)}}
    x = big.NewInt(-5)
    modulus = big.NewInt(20)
    result = evaluatePolynomial(x, modulus, p)
    fmt.Println(result)
    if result.Cmp(big.NewInt(17)) != 0 {
        t.Errorf("Expecting 17, got: %s", result.String())
    }

}

