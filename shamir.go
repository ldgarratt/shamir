package main

import (
    "encoding/base64"
    "crypto/rand"
    "flag"
    "fmt"
    "log"
    "math/big"
    "os"
    "strconv"
    "strings"
    "unicode"
)

// 2**127 - 1.
const PRIME_STRING = "170141183460469231731687303715884105727"

// Split large secrets into smaller ones to avoid wrapping around the prime
// modulus during encoding. Given we only accept ASCII values, splitting into
// chunks of size 15 is sufficent to never be larger than our prime.
const CHUNK_SIZE = 15

// A polynomial is a slice of big.Ints. polynomial[i] is the x^i coefficient.
// E.g. 7x^2 + 5 is [5, 0, 7].
type polynomial struct {
	coefficients []*big.Int
}

// Formats the polynomial so we can print e.g. [5, 0, 7] as 7x^2 + 5.
// TODO: Was useful in debugging but not anymore.
func (p polynomial) format() string {
	degree := len(p.coefficients) - 1

    if degree <= -1 {
        log.Fatal("Empty polynomial\n")
    }

    if degree == 0 {
        return p.coefficients[0].String()
    }

    const_term := p.coefficients[0].String()
    x_term := p.coefficients[1].String()

    if degree == 1 {
        if x_term == "1" && const_term == "0" {
            return "x"
        }
        if x_term == "1" && const_term != "0" {
            return "x + " + const_term
        }
        if x_term == "0" {
            return const_term
        }
        return x_term + "x + " + const_term
    }

    final_term := p.coefficients[degree].String()
    // If its not actually an nth degree polynomial because the x^n term is 0,
    // then recursiely apply the method on the actual (n-1)th degree polynomial
    if final_term == "0" {
        var new_poly []*big.Int
        for i := 0; i <= degree - 1;  i++ {
            new_poly = append(new_poly, p.coefficients[i])
        }
        p = polynomial{new_poly}
        return p.format()
    }

    // From here we can assume the polynomial is at least degree 2.
    str := ""

    if const_term != "0" {
        str = " + " + const_term
    }

    // x term of polynomial
    switch x_term {
        case "0":
        case "1":
            str = " + " + "x" + str
        default:
            str = " + " + x_term + "x" + str
    }

    // middle part of polynomial (does nothing if degree 2)
    for i := 2; i < degree; i++ {
        coeff := p.coefficients[i].String()
        switch coeff {
        case "0":
            continue
        case "1":
            str = " + " + "x^" + strconv.Itoa(i) + str
        default:
            str = " + " + coeff + "x^" + strconv.Itoa(i) + str
        }
    }

    // Final term must be non-zero, otherwise it's not a polynomial of that
    // degree
    if final_term == "1" {
        final_term = ""
    }
    str = final_term + "x^" + strconv.Itoa(degree) + str

	return str
}

// Generates a random polynomial with specified constant, degree and modulus
// For Shamir Secret Sharing:
// constant = secret to split up
// degree of polynomial = # shares to split up
// modulus = prime for the galois field arithmetic
func generateRandomPolynomial(constant, modulus *big.Int, degree int) polynomial {

    // Start with the (pre-selected) constant term of the polynomial
    coefficients := []*big.Int{constant}

    // Randomly select all the middle terms
	for i := 1; i < degree; i++ {
        num, err := rand.Int(rand.Reader, modulus)
        if err != nil {
            panic(err)
        }
        coefficients = append(coefficients, num)
    }

    // Randomly select the final term and ensure it is non-zero so it is really
    // a polynomial of the nth degree.
    for {
        num, err := rand.Int(rand.Reader, modulus)
        if err != nil {
            panic(err)
        }

        if num.Cmp(big.NewInt(0)) != 0 {
            coefficients = append(coefficients, num)
            break
        }
    }
    return polynomial{coefficients}
}

// Evaluates galois polynomial at x
func evaluatePolynomial(x, modulus *big.Int, p polynomial) *big.Int {
    degree := len(p.coefficients) - 1
    result := big.NewInt(0)
    term := big.NewInt(0)
	for n := 0; n <= degree; n++ {
        // x^n mod m
        term.Exp(x, big.NewInt(int64(n)), modulus)

        // a_n * x^n
        term.Mul(term, p.coefficients[n])

        // Add together all the terms
        result.Add(result, term)
    }
    return result.Mod(result, modulus)
}

// Used for testing and in the call to shamirSplitSecret. Not secure to call
// directly unless the polynomial is generated with generateRandomPolynomial.
func _shamirSplitSecretWithFixedPolynomial(secret, modulus *big.Int, poly polynomial, n, t int) []*big.Int {
    var shares []*big.Int
    for x := 1; x <= n; x ++ {
        shares = append(shares, evaluatePolynomial(big.NewInt(int64(x)), modulus, poly))
    }
    return shares
}

// Shamir Secret Sharing splitting secret into n shares with threshold t to
// recover the secret.
func shamirSplitSecret(secret, modulus *big.Int, n, t int) []*big.Int {
    poly := generateRandomPolynomial(secret, modulus, t - 1)
    return _shamirSplitSecretWithFixedPolynomial(secret, modulus, poly, n, t)
}

// Calculates f(0) (mod m) given len(points) == treshhold
// Points are the secret shares (x1, y1), (x2, y2), etc on the polynomial.
func lagrange(points map[int]big.Int, modulus *big.Int) *big.Int {
    result := big.NewInt(0)

    // This part is the outer sum of the Lagrange formula. At each iteration, it
    // adds the y * product term. The product is calculated in the other
    // inner loop.
    for x, y := range points {

        // Calculate the product to multiply against the y term.
        prod := big.NewInt(1)
        for m, _ := range points {
            if m == x {
                continue
            }
            d := new(big.Int).Sub(big.NewInt(int64(m)), big.NewInt(int64(x)))
            d.ModInverse(d, modulus)
            d.Mul(big.NewInt(int64(m)), d)
            prod.Mul(prod, d)
            prod.Mod(prod, modulus)
        }

        // Multiply the product by y and then add to the sum total.
        // Repeat until the sum is complete.
        term := new(big.Int).Mul(&y, prod)
        result = result.Add(result, term)
    }
    return result.Mod(result, modulus)
}

func stringToBigInt(s string) *big.Int {
    return new(big.Int).SetBytes([]byte(s))
}

func bigIntToString(i *big.Int) string {
    data := i.Bytes()
    s := base64.StdEncoding.EncodeToString(data)
    q, _ := base64.StdEncoding.DecodeString(s)
    return string(q)
}

func isASCII(s string) bool {
    for i := 0; i < len(s); i++ {
        if s[i] > unicode.MaxASCII {
            return false
        }
    }
    return true
}

func validParameters(splitSecret *string, splitn, splitthreshold *int) bool {
    if *splitSecret == "" {
        fmt.Println("Empty secret.\nSee README.md for example usage")
        return false
    }
    if !isASCII(*splitSecret) {
        fmt.Println("Secret must be ASCII.")
        return false
    }
    if *splitn < 1 {
        fmt.Println("Number of shares less than 1.\nSee README.md for example useage")
        return false
    }
    if *splitthreshold < 1 {
        fmt.Println("Threshold less than 1.\nSee README.md for example useage")
        return false
    }
    if *splitn < * splitthreshold {
        fmt.Println("Number of shares is less than the threshold.")
        return false
    }
    return true
}

// Splits a string into an array of smaller strings, where each element is a
// string of size <= i. They will all be size n, except the last chunk which
// might be smaller. No padding required.
// We use this function to split a large secret into smaller ones, because
// otherwise secrets would wrap around the modulus during encoding.
func splitStringIntoChunks(s string, n int) []string {
    res := []string{}
    // first chunks of length 5. Spare chunk left over if not exactly multiple
    // of n
    for i := 0; i <= len(s) - n; i += n {
        res = append(res, s[i:i+n])
    }
    // Add the spare part if it exists.
    if len(s) % n != 0 {
        res = append(res, s[len(res)*n:])
    }
    return res
}

// Takes as input a slice of bigInt slices and pairwise combines them and
// returns a slice of strings.
// We use this when splitting up a secret into individual subsecrets, shamir
// splitting each one, then re-combining each share into a master share in a
// useful string format for later parsing.
// E.g. for a secret which is split into 3 subsecrets and shared among 2 people:
// subsecret 1 shares: [23, 345]
// subsecret 2 shares: [100, 99]
// subsecret 3 shares: [19, 50]
// Input: [[23, 345], [100, 99], [19, 50]]
// Returns: [23+100+19, 345+99+50]
func pairwiseJoinSlices(a [][]*big.Int) []string {
    inner_slice_len := len(a[0])
    for _, element := range(a) {
        if len(element) != inner_slice_len {
            panic("Slices not all the same length.")
        }
    }
    res := []string{}
    // Build the jth string
    // To build the jth string, take the jth element from each inner slice and
    // add a "+" at the end for each, except the last.
    for j := 0; j <= inner_slice_len  - 1; j++ {
        str := ""
        for i := 0; i < len(a) - 1; i++ {
            str += a[i][j].String() + "+"
        }
        str += a[len(a)-1][j].String()
        // Add the jth string to result
        res = append(res, str)
    }
    return res
}

// Given an input like ./shamir combine 2 334343+23232 4 32312321+2312312, this
// will create a slice of maps like:
// [[2: 334343, 4: 32312321], [2: 23232, 4: 2312312]]
// Each map in the slice is itself a subsecret puzzle to solve with Lagrange.
func createSubsecretSliceMap(s []string) []map[int]big.Int {
    num_subsecrets := len(strings.Split(s[1], "+"))
    res := make([]map[int]big.Int, num_subsecrets)
    for i := 0; i < num_subsecrets; i++ {
        subsecret_map := make(map[int]big.Int)
        for j := 0; j < len(s); j += 2 {
            x, _ := strconv.Atoi((s[j]))
            y := new(big.Int)
            y, _ = y.SetString(strings.Split(s[j+1], "+")[i], 10)
            subsecret_map[x] = *y
        }
        res = append(res, subsecret_map)
    }
    return res
}

func main() {
    PRIME, _ := new(big.Int).SetString(PRIME_STRING, 10)

    splitCmd := flag.NewFlagSet("split", flag.ExitOnError)
    splitSecret := splitCmd.String("secret", "", "Secret to split")
    splitn := splitCmd.Int("n", 0, "Number of shares to split secret into")
    splitthreshold := splitCmd.Int("t", 0, "Threshold needed to repiece together secret")

    combineCmd := flag.NewFlagSet("combine", flag.ExitOnError)

    if len(os.Args) < 2 {
        fmt.Println("Expected 'split' or 'combine' subcommands.\nSee README.md for example usage.")
        os.Exit(1)
    }

    switch os.Args[1] {
        case "split":
            splitCmd.Parse(os.Args[2:])
            if !validParameters(splitSecret, splitn, splitthreshold) {
                os.Exit(1)
            }
            fmt.Println("Secret to split:", *splitSecret )

            secret := splitStringIntoChunks(*splitSecret, CHUNK_SIZE)
            var result [][]*big.Int
            for _, subsecret := range secret {
                subsecret_int := stringToBigInt(subsecret)
                subsecret_shares := shamirSplitSecret(subsecret_int, PRIME, *splitn, *splitthreshold)
                result = append(result, subsecret_shares)
            }
            shares := pairwiseJoinSlices(result)
            for i, _ := range(shares) {
                fmt.Printf("Share %d: (%d, %s)\n", i+1, i+1, shares[i])
            }

        case "combine":
            combineCmd.Parse(os.Args[2:])
            tail := combineCmd.Args()
            m := createSubsecretSliceMap(tail)
            res := ""
            for subsecretpuzzle := range(m) {
                subsecret := lagrange(m[subsecretpuzzle], PRIME)
                res += bigIntToString(subsecret)
            }
            fmt.Println(res)
        default:
            fmt.Println("Expected 'split' or 'combine' subcommands")
            os.Exit(1)
        }
}
