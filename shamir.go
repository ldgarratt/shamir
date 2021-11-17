package main

import (
	"crypto/rand"
	"fmt"
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

// Formats the polynomial so we can print [5, 0, 7] as 7x^2 + 5.
// TODO: more edge cases of [0, 0, 0 , 2] and [1, 0, 0, 0], etc. Even though in
// SSS they won't ever be hit because the secret can't be 0 and the polynomial
// of degree n can't have x^n coefficent equal to 0.
// also I just hate how ugly this method is
func (p polynomial) format() string {
	degree := len(p.coefficients) - 1
    str := ""
    coeff := p.coefficients[degree].String()
    if !((coeff == "0") || (coeff == "1"))  {
        str = coeff + "x^" + strconv.Itoa(degree)
    } else if (coeff == "1") {
        str = "x"
    } else {
        str = ""
    }
	for i := degree - 1; i >= 0; i-- {
		coeff := p.coefficients[i].String()
		if coeff == "0" {
			continue
		}
        if ((coeff == "1") && i != 0) {
            coeff = " + " + ""
        }
		switch i {
		case 0:
            if coeff == "1" {
                str += " + " + coeff
            } else {
			str += " + " + coeff
            }
		case 1:
			str += " + " + coeff + "x"
		default:
			str += " + " + coeff + "x^" + strconv.Itoa(i)
		}
	}
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

//func (poly Polynomial)

// Evaluates polynomial
// Write some tests
func evaluatePolynomial(x, modulus *big.Int, p polynomial) *big.Int {
    degree := len(p.coefficients) - 1
    result := big.NewInt(0)
    term := big.NewInt(0)
	for e := degree; e >= 0; e-- {
        term.Exp(x, big.NewInt(int64(e)), nil)
        fmt.Println(term)
        term.Mul(term, p.coefficients[e])
        result.Add(result, term)
    }
    return result.Mod(result, modulus)
}

// TODO: function to give each person a point (x and evaluatePolynomial(x)

// TODO: make the function say "Welcome to Shamir's secret sharing scheme! What
// is your secret you wish to split?
// logic says it must be an int, if not, re-prompt the user.
// Next, it asks how many shares you want to split it into (n)
// again, must be an int, obviously.
// Next, ask the user for the threshold int.

// in future, hae a command line option too.
func main() {

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

    x := big.NewInt(2)
    modulus := big.NewInt(17)
    result := evaluatePolynomial(x, modulus, letters)
    fmt.Println(result)


    new_letters := generateRandomPolynomial(big.NewInt(11), 5, PRIME)
    fmt.Println(new_letters.format())
}
