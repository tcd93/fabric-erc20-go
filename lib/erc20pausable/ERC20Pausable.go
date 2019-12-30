package erc20pausable

import (
	. "erc20/helpers"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/*Token pausable implements PausableTokenInterface*/
type Token struct{}

/*IsPaused get the "isPaused" state of token*/
func (t *Token) IsPaused(stub shim.ChaincodeStubInterface) (bool, error) {
	isPaused, err := stub.GetState("isPaused")
	if err != nil {
		return true, err
	}

	if string(isPaused) == "" {
		return false, nil
	}
	return strconv.ParseBool(string(isPaused))
}

/*Pause freezes the transfer/approve functions of the token, callable by token owner*/
func (t *Token) Pause(stub shim.ChaincodeStubInterface,
	getOwner func(shim.ChaincodeStubInterface) (string, error),
) error {
	callerID, err := GetCallerID(stub)
	if err != nil {
		return err
	}

	tokenOwnerID, err := getOwner(stub)
	if err != nil {
		return err
	}

	if err := CheckCallerIsOwner(callerID, tokenOwnerID); err != nil {
		return err
	}

	return stub.PutState("isPaused", []byte("true"))
}

/*Unpause un-freezes the transfer/approve functions of the token, callable by token owner*/
func (t *Token) Unpause(stub shim.ChaincodeStubInterface,
	getOwner func(shim.ChaincodeStubInterface) (string, error),
) error {
	callerID, err := GetCallerID(stub)
	if err != nil {
		return err
	}

	tokenOwnerID, err := getOwner(stub)
	if err != nil {
		return err
	}

	if err := CheckCallerIsOwner(callerID, tokenOwnerID); err != nil {
		return err
	}

	return stub.PutState("isPaused", []byte("false"))
}
