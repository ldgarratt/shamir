package main

import (
    "testing"
    "math/big"
    "reflect"
)

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

    // Evaluate polynomial twice as different points to ensure it does not
    // change when messing around with pointers.
    p = polynomial{[]*big.Int{big.NewInt(1234), big.NewInt(166), big.NewInt(94)}}
    modulus = big.NewInt(1613)
    x = big.NewInt(0)
    result = evaluatePolynomial(x, modulus, p)
    if result.Cmp(big.NewInt(1234)) != 0 {
        t.Errorf("Expecting 1234, got: %s", result.String())
    }

    x = big.NewInt(1)
    result = evaluatePolynomial(x, modulus, p)
    if result.Cmp(big.NewInt(1494)) != 0 {
        t.Errorf("Expecting 1494, got: %s", result.String())
    }
}

func Test_shamirSplitSecretWithFixedPolynomial(t *testing.T) {
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
            t.Errorf("Expecting %s at x=%d, got: %s", expected[i], i+1, result[i])
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

func TestBigIntStringEncodingDecoding(t *testing.T) {
    str := "Hello, world!"
    enc := stringToBigInt(str)
    result := bigIntToString(enc)

    if result != str {
        t.Errorf("Expecting %s, got: %s", str, result)
    }
}

func TestIsASCII(t *testing.T) {
    str:= "Hello, World!"
    if !isASCII(str) {
        t.Errorf("Expected %s to be ASCII", str)
    }

    str = "🧡💛💚💙💜"
    if isASCII(str) {
        t.Errorf("Expected %s not to be ASCII", str)
    }
}

func TestSplitStringIntoChunks(t *testing.T) {
    str:= "Hello, World!"
    result := splitStringIntoChunks(str, 3)
    expected := []string{"Hel", "lo,", " Wo", "rld", "!"}
    for i := 0; i < len(result); i++ {
        if result[i] != expected[i] {
            t.Errorf("Expected %s, got %s", expected, result)
        }
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

    for i := 0; i < len(result); i++ {
        if (result[0] != expected[0])   {
            t.Errorf("Expecting %s, got: %s", expected, result)
        }
    }
}

func TestCreateSubsecretSliceMap(t *testing.T) {
    s := []string{"2", "334343", "4", "32312321"}
    result := createSubsecretSliceMap(s)
    m1 := map[int]big.Int{
        2 : *big.NewInt(334343),
        4 : *big.NewInt(32312321),
    }
    expected := []map[int]big.Int{m1}
    if len(expected) != len(result) {
        t.Error("Expected slice length is %i, result length is %i", len(expected), len(result))
    }
    for i, elem := range(expected) {
        if reflect.DeepEqual(expected[i], result[i]) == false {
            t.Error("Slices are not the same. Expected %i, got %i", elem, result[i])
        }
    }

    s = []string{"2", "334343+23232", "4", "32312321+2312312"}
    result = createSubsecretSliceMap(s)
    m2 := map[int]big.Int{
        2 : *big.NewInt(23232),
        4 : *big.NewInt(2312312),
    }
    expected = []map[int]big.Int{m1, m2}
    if len(expected) != len(result) {
        t.Error("Expected slice length is %i, result length is %i", len(expected), len(result))
    }
    for i, elem := range(expected) {
        if reflect.DeepEqual(expected[i], result[i]) == false {
            t.Error("Slices are not the same. Expected %i, got %i", elem, result[i])
        }
    }

    s = []string{"2", "334343+23232+0", "4", "32312321+2312312+234"}
    result = createSubsecretSliceMap(s)
    m3 := map[int]big.Int{
        2 : *big.NewInt(0),
        4 : *big.NewInt(234),
    }
    expected = []map[int]big.Int{m1, m2, m3}
    if len(expected) != len(result) {
        t.Error("Expected slice length is %i, result length is %i", len(expected), len(result))
    }
    for i, elem := range(expected) {
        if reflect.DeepEqual(expected[i], result[i]) == false {
            t.Error("Slices are not the same. Expected %i, got %i", elem, result[i])
        }
    }
}

