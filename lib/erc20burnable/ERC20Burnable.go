package erc20burnable

import (
	. "erc20/helpers"
	"erc20/lib/erc20events"
	"fmt"
	"math/big"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("burn-logger")

/*Token burnable implementation of BurnableTokenInterface*/
type Token struct{}

/*Burn destroys an amount of tokens of the invoking identity, and total supply.

* `args[0]` - the amount that will be burnt*/
func (t *Token) Burn(stub shim.ChaincodeStubInterface,
	args []string,
	getTotalSupply func(stub shim.ChaincodeStubInterface) (*big.Int, error),
	getBalanceOf func(stub shim.ChaincodeStubInterface, args []string) (*big.Int, error),
) error {
	sValue := args[0]

	burneeID, err := GetCallerID(stub)
	if err != nil {
		return err
	}

	if err := CheckGreaterThanZero(sValue); err != nil {
		return err
	}

	burnAmount := StringToBigInt(sValue)

	logger.Infof("Burn: burning %v tokens from %v", burnAmount, burneeID)

	burnerBalance, err := getBalanceOf(stub, []string{burneeID})
	if err := CheckBalance(burnerBalance, burneeID); err != nil {
		return err
	}
	if err := IsSmallerOrEqual(burnAmount, burnerBalance); err != nil {
		return fmt.Errorf("burn amount should be less than balance of sender (%v): %v", burneeID, err)
	}

	err = stub.PutState(burneeID, []byte(Sub(burnerBalance, burnAmount).String()))
	if err != nil {
		return err
	}

	totalSupply, err := getTotalSupply(stub)
	if err != nil {
		return err
	}
	if err := IsSmallerOrEqual(burnAmount, totalSupply); err != nil {
		return fmt.Errorf("burn amount should be less than total supply (%v): %v", totalSupply, err)
	}

	err = stub.PutState("totalSupply", []byte(Sub(totalSupply, burnAmount).String()))
	if err != nil {
		return err
	}

	json := MalshalJSON(erc20events.Event{Origin: burneeID, Payload: erc20events.Payload{From: burneeID, To: "", Amount: burnAmount}})
	return stub.SetEvent(erc20events.TRANSFER, json)
}

/*BurnFrom burns a specific amount of tokens from the target identity and total supply,
the chaincode invoker must have sufficient allowance from burnee.

* `args[0]` - the ID of burnee.

* `args[1]` - the burn amount.*/
func (t *Token) BurnFrom(stub shim.ChaincodeStubInterface,
	args []string,
	getAllowance func(stub shim.ChaincodeStubInterface, args []string) (*big.Int, error),
	getTotalSupply func(stub shim.ChaincodeStubInterface) (*big.Int, error),
	getBalanceOf func(stub shim.ChaincodeStubInterface, args []string) (*big.Int, error),
) error {
	burneeID, sValue := args[0], args[1]

	if err := CheckGreaterThanZero(sValue); err != nil {
		return err
	}

	burnerID, err := GetCallerID(stub)
	if err != nil {
		return err
	}

	burnAmount := StringToBigInt(sValue)

	logger.Infof("BurnFrom: burning %v tokens from %v with %v...", burnAmount, burneeID, burnerID)

	balanceOfBurnee, err := getBalanceOf(stub, []string{burneeID})
	if err != nil {
		return err
	}
	approvedAmount, err := getAllowance(stub, []string{burneeID, burnerID})
	if err != nil {
		return err
	}
	if err := CheckBalance(balanceOfBurnee, burneeID); err != nil {
		return err
	}
	if err := CheckBalance(approvedAmount, burneeID+"-"+burnerID); err != nil {
		return err
	}
	if err := IsSmallerOrEqual(burnAmount, balanceOfBurnee); err != nil {
		return fmt.Errorf("burn amount should be less than balance of burnee (%v): %v", burneeID, err)
	}

	err = stub.PutState(burneeID, []byte(Sub(balanceOfBurnee, burnAmount).String()))
	if err != nil {
		return err
	}

	totalSupply, err := getTotalSupply(stub)
	if err != nil {
		return err
	}
	if err := IsSmallerOrEqual(burnAmount, totalSupply); err != nil {
		return fmt.Errorf("burn amount should be less than total supply (%v): %v", totalSupply, err)
	}

	err = stub.PutState("totalSupply", []byte(Sub(totalSupply, burnAmount).String()))
	if err != nil {
		return err
	}

	json := MalshalJSON(erc20events.Event{Origin: burnerID, Payload: erc20events.Payload{From: burneeID, To: "", Amount: burnAmount}})
	return stub.SetEvent(erc20events.TRANSFER, json)
}
