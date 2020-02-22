package helpers

import (
	"fmt"
	"math/big"
	"strconv"
)

/*CheckArgsLength compares length of string with `expectedLength`*/
func CheckArgsLength(args []string, expectedLength int) error {
	if len(args) != expectedLength {
		return fmt.Errorf("invalid number of arguments. Expected %v, got %v", expectedLength, len(args))
	}
	return nil
}

/*CheckGreaterThanZero parses a string value to number and check if it's >= 0*/
func CheckGreaterThanZero(value string) error {
	n, err := strconv.ParseFloat(value, 64)
	if n < 0 {
		return fmt.Errorf("parsed version of %v, should be >= 0", value)
	}
	return err
}

/*CheckBalance checks if sender's balance is > 0*/
func CheckBalance(balance *big.Int, mspID string) error {
	if balance.Cmp(big.NewInt(0)) == -1 {
		return fmt.Errorf("Balance of sender %v is %v", mspID, balance)
	}
	return nil
}

/*CheckTotalSupply checks if totalSupplyAmount is valid*/
func CheckTotalSupply(amount *big.Int) error {
	if amount.Cmp(big.NewInt(0)) == -1 {
		return fmt.Errorf("Total supply is < 0")
	}
	return nil
}

/*IsSmallerOrEqual returns `nil` if a is <= b*/
func IsSmallerOrEqual(a *big.Int, b *big.Int) error {
	if a.Cmp(b) == 1 {
		return fmt.Errorf("%v should be <= to %v", a, b)
	}
	return nil
}

/*CheckApproved checks if approved amount is > 0*/
func CheckApproved(approved *big.Int, key string) error {
	if approved.Cmp(big.NewInt(0)) == -1 {
		return fmt.Errorf("Approved amount for %v is %v", key, approved)
	}
	return nil
}

/*CheckCallerIsOwner compares two strings*/
func CheckCallerIsOwner(caller string, owner string) error {
	if caller != owner {
		return fmt.Errorf("Function only accessible to token owner: %v", owner)
	}
	return nil
}
