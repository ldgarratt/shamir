package main

import (
    "testing"
    "fmt"
    "math/big"
)

func TestFormat(t *testing.T) {
    p := polynomial{[]*big.Int{big.NewInt(2), big.NewInt(4), big.NewInt(3), big.NewInt(0), big.NewInt(2)}}
    fmt.Println(p)
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

// TODO: write function to evaluate polynomial

/*
func TestProblem10(t *testing.T) {
	msg := []byte("YELLOW SUBMARINEYELLOW SUBMARINE")
	iv := make([]byte, 16)
	b, _ := aes.NewCipher([]byte("YELLOW SUBMARINE"))
	res := decryptCBC(encryptCBC(msg, b, iv), b, iv)
	if !bytes.Equal(res, msg) {
		t.Errorf("%q", res)
	}

	msg = decodeBase64(t, string(readFile(t, "10.txt")))
	t.Logf("%s", decryptCBC(msg, b, iv))
}
*/
