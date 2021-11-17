package main

import (
    "testing"
    "fmt"
)

// TODO: write this to test that e.g. the polynomial [5, 3, 0, 2] is printed
// correctly as 2x^3 + 3x + 5
func TestFormat(t *testing.T) {
	// t.Fatal("not implemented")
    p := polynomial{[]int{2, 4, 3, 0, 2}}
    fmt.Println(p)
    result := "x"
    result = p.format()
    if result != "2x^4 + 3x^2 + 4x + 2" {
        t.Errorf("expecting 2x^4 + 3x^2 + 4x + 2, got %s", result)
    }
}

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
