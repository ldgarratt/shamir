# Shamir Secret Sharing

Splits a secret s into n parts with threshold t according to the SSS algorithm.


This was mainly to practice implementing something crypto in Go. I make no
claims about the actual security of this implementation :)

## Usage:


To split e.g. secret = "my secret" between n = 6 people with threshold t = 4:

./shamir split -secret="my secret" -t=4 -n=6

Example output:

Person 1: share: (1, 28485993219178821771714474316113792564)
Person 2: share: (2, 37142848113374314264863710222663272796)
Person 3: share: (3, 64607012802959991351641665784703859409)
Person 4: share: (4, 149514935408309364884819092075808491128)
Person 5: share: (5, 160361880589326714985479436453666000951)
Person 6: share: (6, 135784296466385553506393449991849327603)

To combine (using 4 shares):

./shamir combine 1 28485993219178821771714474316113792564 3 64607012802959991351641665784703859409 4 149514935408309364884819092075808491128 5 160361880589326714985479436453666000951

Output:

"my secret"
