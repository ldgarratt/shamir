package main

import (
    "encoding/base64"
    "crypto/rand"
    "flag"
    "fmt"
    "math/big"
    "os"
    "strconv"
    "strings"
    "regexp"
    "unicode"
)

// 2^127 - 1.
const PRIME = "170141183460469231731687303715884105727"

// Split large secrets into smaller subsecrets to avoid wrapping around the
// prime modulus after encoding them as bigInts. Given we only accept ASCII
// secrets, splitting into chunks of size 15 is sufficent to never be larger
// than the prime.
const CHUNK_SIZE = 15

// A polynomial is a slice of big.Ints. polynomial[i] is the x^i coefficient.
// E.g. 7x^2 + 5 is [5, 0, 7].
type polynomial struct {
	coefficients []*big.Int
}

// Generates a random polynomial with specified constant, degree and modulus
// For Shamir Secret Sharing:
// constant = secret to split up
// degree of polynomial = threshold - 1
// modulus = PRIME for the galois field arithmetic
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

// Evaluates galois polynomial at x using Horner's method.
func evaluatePolynomial(x, modulus *big.Int, p polynomial) *big.Int {
    degree := len(p.coefficients) - 1
    result := new(big.Int).Set(p.coefficients[degree])
    for i := degree - 1; i >= 0; i-- {
        result.Mul(result, x)
        result.Add(result, p.coefficients[i])
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

// Calculates f(0) (mod m) given len(points) == threshhold
// Points are the secret shares (x1, y1), (x2, y2), etc. on the polynomial.
func lagrange(points map[int]big.Int, modulus *big.Int) *big.Int {
    result := big.NewInt(0)

    // This part is the outer sum of the Lagrange formula. At each iteration, it
    // adds the y * product term. The product is calculated in the inner loop.
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

// Reversibly encodes an arbitary string into a bigInt.
func stringToBigInt(s string) *big.Int {
    return new(big.Int).SetBytes([]byte(s))
}

// Reversibly decodes an arbitary bigInt into a string.
func bigIntToString(i *big.Int) string {
    data := i.Bytes()
    s := base64.StdEncoding.EncodeToString(data)
    res, err := base64.StdEncoding.DecodeString(s)
    if err != nil {
        panic("Decoding failed.")
    }
    return string(res)
}

func isASCII(s string) bool {
    for i := 0; i < len(s); i++ {
        if s[i] > unicode.MaxASCII {
            return false
        }
    }
    return true
}

func validSplitParameters(splitSecret *string, splitn, splitthreshold *int) bool {
    if *splitSecret == "" {
        fmt.Println("Empty secret.\nSee README.md for example usage.")
        return false
    }
    if !isASCII(*splitSecret) {
        fmt.Println("Secret must be ASCII.")
        return false
    }
    if *splitn < 1 {
        fmt.Println("Number of shares less than 1.\nSee README.md for example usage.")
        return false
    }
    if *splitthreshold < 2 {
        fmt.Println("Threshold less than 2.\nSee README.md for example usage.")
        return false
    }
    if *splitn < * splitthreshold {
        fmt.Println("Number of shares is less than the threshold.")
        return false
    }
    return true
}

// Splits a string into chunks, preserving order. Each substring will be size
// n, except possibly the last which might be smaller.
// E.g. for n = 3:
// Input: "Hello, World"
// Returns: ["Hel", "Lo,", " Wo", "rld", "!"]
// In Shamir Secret Sharing, we use this function to split a large secret into
// smaller ones because otherwise secrets would wrap around the modulus after
// encoding.
func splitStringIntoChunks(s string, n int) []string {
    res := []string{}
    // first chunks of length n. Spare chunk left over if not exactly multiple
    // of n
    for i := 0; i <= len(s) - n; i += n {
        res = append(res, s[i:i+n])
    }
    // Add the spare chunk if it exists.
    if len(s) % n != 0 {
        res = append(res, s[len(res)*n:])
    }
    return res
}

// Takes as input a slice of bigInt slices and pairwise zips them into a slice
// of strings concatenated with "+".
// In Shamir Secret Sharing, we use this when splitting up a secret into
// individual subsecrets. We perform SSS on each subsecret. The input is the
// slices of shares making up each subsecret. The output is rearranging it
// so subshares i are all together.
// E.g. for a secret which is split into 3 subsecrets and shared among 2 people:
// subsecret 1 shares: [23, 345]
// subsecret 2 shares: [100, 99]
// subsecret 3 shares: [19, 50]
// Input: [[23, 345], [100, 99], [19, 50]]
// Returns: [23+100+19, 345+99+50]
// This means you only need to give person 1 the share (1, 23+100+19) and
// person 2 the share (2, 345+99+50)
func pairwiseJoinSlices(a [][]*big.Int) []string {
    inner_slice_len := len(a[0])
    for _, element := range(a) {
        if len(element) != inner_slice_len {
            panic("Slices are not all the same length.")
        }
    }
    res := []string{}
    // Build the jth string
    for j := 0; j <= inner_slice_len  - 1; j++ {
        str := ""
        // To build the jth string, take the jth element from each inner slice and
        // add a "+" at the end for each, except the last.
        for i := 0; i < len(a) - 1; i++ {
            str += a[i][j].String() + "+"
        }
        str += a[len(a)-1][j].String()
        // Add the jth string to result
        res = append(res, str)
    }
    return res
}

func validCombineParameters(s []string) bool {
    if len(s) < 4 {
        fmt.Println("Must combine at least two shares.\nSee README.md for example usage.")
        return false
    }
    if len(s) % 2 != 0 {
        fmt.Println("Combine command takes an even number of arguments.\nSee README.md for example usage.")
        return false
    }

    var digitCheck = regexp.MustCompile(`^[0-9]+$`)
    num_subsecrets := len(strings.Split(s[1], "+"))

    // Check the shares (the even numbered parameters)
    for i := 1; i <= len(s); i+=2 {
        share := strings.Split(s[i], "+")
        if len(share) != num_subsecrets {
            fmt.Println("Each share must contain the same number of subsecrets (numbers separated by '+').\nSee README.md for example usage.")
            return false
        }
        for _,subsecret := range(share) {
            if !digitCheck.MatchString(subsecret) {
                fmt.Println("Shares must be of the form: 'int+int+int+..+int'.\nSee README.md for example usage.")
                return false
            }
        }
    }

    // Check the share number (the odd numbered parameters, i.e the x co-ordinates)
    for i := 0; i < len(s); i+=2 {
        _, err := strconv.Atoi(s[i])
        if err != nil {
            fmt.Println("Share numbers must be 64-bit ints.\nSee README.md for example usage.")
        }
    }
    return true
}

// Given an input like ./shamir combine 2 334343+23232 4 32312321+2312312, this
// will create a slice of maps like:
// [[2: 334343, 4: 32312321], [2: 23232, 4: 2312312]]
// Each map in the slice is itself a subsecret puzzle to solve with Lagrange.
func createSubsecretSliceMap(s []string) []map[int]big.Int {
    num_subsecrets := len(strings.Split(s[1], "+"))
    res := []map[int]big.Int{}
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
    PRIME, _ := new(big.Int).SetString(PRIME, 10)

    splitCmd := flag.NewFlagSet("split", flag.ExitOnError)
    splitSecret := splitCmd.String("secret", "", "Secret to split.")
    splitn := splitCmd.Int("n", 0, "Number of shares to split secret into.")
    splitthreshold := splitCmd.Int("t", 0, "Threshold needed to repiece together secret.")

    combineCmd := flag.NewFlagSet("combine", flag.ExitOnError)

    if len(os.Args) < 2 {
        fmt.Println("Expected 'split' or 'combine' subcommands.\nSee README.md for example usage.")
        os.Exit(1)
    }

    switch os.Args[1] {
        // TODO: make go routine split/combine each subsecret for concurrency
        case "split":
            splitCmd.Parse(os.Args[2:])
            if !validSplitParameters(splitSecret, splitn, splitthreshold) {
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
            if !validCombineParameters(tail) {
                os.Exit(1)
            }
            m := createSubsecretSliceMap(tail)
            res := ""
            for _, subsecretpuzzle := range(m) {
                subsecret := lagrange(subsecretpuzzle, PRIME)
                res += bigIntToString(subsecret)
            }
            fmt.Println(res)
        default:
            fmt.Println("Expected 'split' or 'combine' subcommands. See README.md for example usage.")
            os.Exit(1)
        }
}
