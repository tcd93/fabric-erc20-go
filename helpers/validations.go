package helpers

import (
	"fmt"
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
func CheckBalance(balance float64, mspID string) error {
	if balance < 0 {
		return fmt.Errorf("Balance of sender %v is %v", mspID, balance)
	}
	return nil
}

/*CheckTotalSupply checks if totalSupplyAmount is valid*/
func CheckTotalSupply(amount float64) error {
	if amount < 0 {
		return fmt.Errorf("Total supply is < 0")
	}
	return nil
}

/*IsSmallerOrEqual returns `nil` if a is <= b*/
func IsSmallerOrEqual(a float64, b float64) error {
	if a > b {
		return fmt.Errorf("%v should be <= to %v", a, b)
	}
	return nil
}

/*CheckApproved checks if approved amount is > 0*/
func CheckApproved(approved float64, key string) error {
	if approved <= 0 {
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
