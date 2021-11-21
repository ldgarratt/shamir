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
)

// TODO: maybe split this up so main just does the shamir stuff, while the
// polynomial manipulation is a separate package

// A polynomial is a slice of big.Ints. polynomial[i] is the x^i coefficient.
// E.g. 7x^2 + 5 is [5, 0, 7].
type polynomial struct {
	coefficients []*big.Int
}

// Formats the polynomial so we can print e.g. [5, 0, 7] as 7x^2 + 5.
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
// For SSS:
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
func _shamirSplitSecretwithFixedPolynomial(secret, modulus *big.Int, poly polynomial, n, t int) []*big.Int {
    var shares []*big.Int
    fmt.Printf("The Shamir polynomial is: %s (mod %d)\n", poly.format(), modulus)
    fmt.Printf("The individual shares are:\n")
    for x := 1; x <= n; x ++ {
        shares = append(shares, evaluatePolynomial(big.NewInt(int64(x)), modulus, poly))
        fmt.Printf("Person %d: share: (%d, %d)\n", x, x, shares[x - 1])
    }
    fmt.Println(shares)
    return shares
}

// Shamir Secret Sharing splitting secret into n shares with threshold t to
// recover the secret.
func shamirSplitSecret(secret, modulus *big.Int, n, t int) []*big.Int {
    poly := generateRandomPolynomial(secret, modulus, t - 1)
    return _shamirSplitSecretwithFixedPolynomial(secret, modulus, poly, n, t)
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

// TODO: maybe this should just work on ints only.
// Write tests... write these better too probably
func stringToBigInt(s string) (*big.Int, error) {
    return new(big.Int).SetBytes([]byte(s))
}

// TODO: make cleaner, error handling
func bigIntToString(i *big.Int) (string) {
    data := i.Bytes()
    s := base64.StdEncoding.EncodeToString(data)
    q, _ := base64.StdEncoding.DecodeString(s)
    return string(q)
}

func main() {
    // 2**127 - 1
    PRIME := big.NewInt(2)
    PRIME.Exp(PRIME, big.NewInt(127), nil)
    PRIME.Sub(PRIME, big.NewInt(1))
    fmt.Println("The prime is: ", PRIME)

    splitCmd := flag.NewFlagSet("split", flag.ExitOnError)
    splitSecret := splitCmd.String("secret", "", "Secret to split")
    splitn := splitCmd.Int("n", 0, "Number of shares to split secret into")
    splitthreshold := splitCmd.Int("t", 0, "Threshold needed to repiece together secret")

    combineCmd := flag.NewFlagSet("combine", flag.ExitOnError)

    if len(os.Args) < 2 {
        fmt.Println("expected 'split' or 'combine' subcommands")
        os.Exit(1)
    }

    switch os.Args[1] {
        case "split":
            splitCmd.Parse(os.Args[2:])
            fmt.Println("subcommand 'split'")
            fmt.Println("split secret", *splitSecret )

            if *splitSecret == "" {
                fmt.Println("Empty secret")
                return
            }

            secret := stringToBigInt(*splitSecret)
            fmt.Println("Secret is: ", secret)
            shamirSplitSecret(secret, PRIME, *splitn, *splitthreshold)

        case "combine":
            combineCmd.Parse(os.Args[2:])
            fmt.Println("subcommand 'combine'")
            tail := combineCmd.Args()
            fmt.Println("  tail:", len(tail))
            points := make(map[int]big.Int)
            for i := 0; i < len(tail) - 1; i = i + 2 {
                x, _ := strconv.ParseInt(tail[i], 10, 64)
                y, _ := new(big.Int).SetString(tail[i+1], 10)
                points[int(x)] = *y
            }
            fmt.Println(points)
            secret := lagrange(points, PRIME)
            // TODO: Maybe just work on ints only?
            fmt.Println(bigIntToString(secret))
        default:
            fmt.Println("Expected 'split' or 'combine' subcommands")
            os.Exit(1)
        }
}
