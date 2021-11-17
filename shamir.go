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
const PRIME = 19

// TODO: maybe split this up so main just does the shamir stuff, while the
// polynomial manipulation is a separate package

// TODO: make polynomial slice of bigints, not ints.
// A polynomial is a slice of 64-bit numbers. E.g. 7x^2 + 5 is [5, 0, 7]
// These can be ints since it is field arithmetic
type polynomial struct {
	coefficients []int
}

func (p polynomial) format() string {
	degree := len(p.coefficients) - 1
	str := ""
	for i := degree; i >= 0; i-- {
		coeff := strconv.Itoa(p.coefficients[i])
		if coeff == "0" {
			continue
		}
        if coeff == "1" {
            coeff = ""
        }
		switch i {
		case 0:
			str += coeff
		case 1:
			str += coeff + "x" + " + "
		default:
			str += coeff + "x^" + strconv.Itoa(i) + " + "
		}
	}
	return str
}

// TODO: write a funcction to generate a random polynomial of degree m (m = k
// -1) for sss threshold
func generateRandomPolynomial(constant, degree, modulus int) polynomial {
    var coeffs []int
    coeffs = append(coeffs, constant)
	for i := 1; i < degree; i++ {
        num, err := rand.Int(rand.Reader, big.NewInt(PRIME))
//        if i == degree {
//            for ok := true; ok; ok != (num.String() == "0") {
//                num, err := rand.Int(rand.Reader, big.NewInt(27))
//            }
        if err != nil {
            panic(err)
        }
        coeffs = append(coeffs, int(num.Int64()))
    }
    var number int
    for {
        num, _ := rand.Int(rand.Reader, big.NewInt(PRIME))
        number = int(num.Int64())
        if number != 0 {
            break
        }
    }
    coeffs = append(coeffs, number)
    p := polynomial{coeffs}
    fmt.Println(p.format())
    return polynomial{coeffs}
}

//func (poly Polynomial)

// TODO: function to evaluate polynomial, write tests for it
//func evaluatePolynomial(x int, p polynomial) int {
//}

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

	letters := polynomial{[]int{2, 4, 3, 0, 2}}
	fmt.Println(len([]int{1, 3, 4, 4, 4}))
	fmt.Println(letters.format())


    new_letters := generateRandomPolynomial(11, 5, PRIME)
    fmt.Println(new_letters.format())
}
