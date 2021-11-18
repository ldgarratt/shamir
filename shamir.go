package main

import (
	"crypto/rand"
	"fmt"
    "log"
	"math/big"
    "strconv"
)

// TODO: currently this is a tiny prime and I'm crazily using crypto/rand to
// generate small bigints, converting them to int64 and then converting to int
// for the polynomial. Obiously my polynomial should really be bigints, but this
// is just for POC at the moment.
// TODO: should this be a big int? and also how does randomness work in that
// case?
const PRIME = 5

// TODO: maybe split this up so main just does the shamir stuff, while the
// polynomial manipulation is a separate package

// A polynomial is a slice of big.Ints. E.g. 7x^2 + 5 is [5, 0, 7]
type polynomial struct {
	coefficients []*big.Int
}

// Formats the polynomial so we can print e.g. [5, 0, 7] as 7x^2 + 5.
func (p polynomial) format() string {
	degree := len(p.coefficients) - 1
    if degree < 2 {
        // TODO: add the tedious cases just for fun
        log.Fatal("Polynomial too small\n")
    }
    final_term := p.coefficients[degree].String()

    // If its not actually an nth degree polynomial because the x^n term is 0,
    // then recursiely apply the method on the actual (n-1)th degree polynomial
    if final_term == "0" {
        var new_poly []*big.Int
        for i := 0; i <= degree - 2; i++ {
            new_poly = append(new_poly, p.coefficients[i])
        }
        p = polynomial{new_poly}
        return p.format()
    }

    str := ""
    // The const term is always non-zero, otherwise add logic for this case
    constant_term := p.coefficients[0].String()
    if constant_term != "0" {
        str = " + " + constant_term
    }

    // x term of polynomial
    x_term := p.coefficients[1].String()
    switch x_term {
        case "0":
        case "1":
            str = " + " + "x" + str
        default:
            str = " + " + x_term + "x" + str
    }

    // middle part of polynomial
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

// For SSS:
// constant = secret to split up
// degree (of polynomial) = shares to split up
// modulus = PRIME for the field arithmetic
//TODO: I should just have logic at the start saying degree has to ve >=
func generateRandomPolynomial(constant *big.Int, degree, modulus int) polynomial {
    var coeffs []*big.Int
    coeffs = append(coeffs, constant)
	for i := 1; i < degree; i++ {
        num, err := rand.Int(rand.Reader, big.NewInt(PRIME))
        if err != nil {
            panic(err)
        }
        coeffs = append(coeffs, num)
    }
    var number *big.Int
    for {
        num, err := rand.Int(rand.Reader, big.NewInt(PRIME))
        if err != nil {
            panic(err)
        }
        number = num
        if number.Cmp(big.NewInt(0)) != 0 {
            break
        }
    }
    coeffs = append(coeffs, number)
    p := polynomial{coeffs}
    fmt.Println(p.format())
    return polynomial{coeffs}
}

// Evaluates galois polynomial at x
func evaluatePolynomial(x, modulus *big.Int, p polynomial) *big.Int {
    degree := len(p.coefficients) - 1
    result := big.NewInt(0)
    term := big.NewInt(0)
	for e := degree; e >= 0; e-- {
        term.Exp(x, big.NewInt(int64(e)), nil)
        term.Mul(term, p.coefficients[e])
        result.Add(result, term)
    }
    return result.Mod(result, modulus)
}

// Shamir Secret Sharing splitting secret into n shares with threshold t to
// recover the secret.
// TODO: test this function.. perhaps a hidden version with a fixed polynomial for
// better end-to-end testing? Then the real version generates a random
// polynomial and just calls the fixed polynomial version
func shamirSplitSecret(secret *big.Int, n, t int) []*big.Int {
    var shares []*big.Int
    // TODO: make this a fixed polynomial for initial testing
    poly := generateRandomPolynomial(secret, t - 1, PRIME)
    fmt.Printf("The Shamir polynomial is: %s (mod %d)\n", poly.format(), PRIME)
    fmt.Printf("The individual shares are:\n")
    for x := 1; x <= t; x ++ {
        // TODO: my PRIME should itself be a big.Int already, really.
        // I really do not like casting of x crap we're doing here
        shares = append(shares, evaluatePolynomial(big.NewInt(int64(x)), big.NewInt(PRIME), poly))
        fmt.Printf("Person %d: share: (%d, %d)\n", x, x, shares[x - 1])
    }
    fmt.Println(shares)
    return shares
}

/*
func lagrange(points map[int]big.Int, modulus *big.Int) int {
// Calculates f(0) % p gien points when len(points) >= threshold

    result := big.NewInt(0)

    // Points are the secret shares (x1, y1), (x2, y2), etc on the polynomial.
    for x, y := range points {
        // For testing:

        fmt.Printf("Intermediate calculuation is %i * %i = %i", a, b, c)
    }

    term := big.NewInt(0)
	for e := degree; e >= 0; e-- {
        term.Exp(x, big.NewInt(int64(e)), nil)
        term.Mul(term, p.coefficients[e])
        result.Add(result, term)
    }
    return result.Mod(result, modulus)

}
*/ 

// TODO: make the function say "Welcome to Shamir's secret sharing scheme! What
// is your secret you wish to split?
// logic says it must be an int, if not, re-prompt the user.
// Next, it asks how many shares you want to split it into (n)
// again, must be an int, obviously.
// Next, ask the user for the threshold int.

// in future, hae a command line option too.
func main() {


    secret := big.NewInt(100)
    people := 5
    t := 3
    shares := shamirSplitSecret(secret, people, t)
    fmt.Println(shares)


    bigint := big.NewInt(123)
    fmt.Println(bigint)

	nBig, err := rand.Int(rand.Reader, big.NewInt(27))
	if err != nil {
		panic(err)
	}


	n := nBig.Int64()
	fmt.Printf("Here is a random %T in [0,27) : %d\n", n, n)

	letters := polynomial{[]*big.Int{big.NewInt(2), big.NewInt(4), big.NewInt(3), big.NewInt(0), big.NewInt(2)}}
	fmt.Println(len([]int{1, 3, 4, 4, 4}))
	fmt.Println(letters.format())


    m := make(map[int]big.Int)
	m[2] = *big.NewInt(1942)
    for k, v := range m {
        fmt.Println(k, "value is", v)
    }

    x := big.NewInt(2)
    modulus := big.NewInt(17)
    result := evaluatePolynomial(x, modulus, letters)
    fmt.Println(result)


    new_letters := generateRandomPolynomial(big.NewInt(11), 5, PRIME)
    fmt.Println(new_letters.format())
}
