package erc20mintable

import (
	. "erc20/helpers"
	"erc20/lib/erc20events"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/*Token mintable implements MintableTokenInterface*/
type Token struct{}

/*Mint tokens, add to the `total supply` and `balance` of minter. This function call only be called by token owner.

* `args[0]` - the ID of minter.

* `args[1]` - the mint amount.

* `getOwner` - specifies the function of getting the current owner of token.

* `getBalanceOf` - specifies the function of getting the initial balance of minter.

* `getTotalSupply` - specifies the function of getting the current total supply of tokens.*/
func (t *Token) Mint(stub shim.ChaincodeStubInterface,
	args []string,
	getOwner func(shim.ChaincodeStubInterface) (string, error),
	getBalanceOf func(shim.ChaincodeStubInterface, []string) (float64, error),
	getTotalSupply func(shim.ChaincodeStubInterface) (float64, error),
) error {
	minterID, value := args[0], args[1]

	if err := CheckGreaterThanZero(value); err != nil {
		return err
	}
	transferAmount := StringToFloat(value)

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

	balanceMinter, err := getBalanceOf(stub, []string{minterID})
	if err != nil {
		return err
	}
	if err := CheckBalance(balanceMinter, minterID); err != nil {
		return err
	}

	totalSupplyAmount, err := getTotalSupply(stub)
	if err != nil {
		return err
	}
	if err := CheckTotalSupply(totalSupplyAmount); err != nil {
		return err
	}

	err = stub.PutState("totalSupply", FloatToBuffer(totalSupplyAmount+transferAmount))
	if err != nil {
		return err
	}
	err = stub.PutState(minterID, FloatToBuffer(balanceMinter+transferAmount))
	if err != nil {
		return err
	}

	json := MalshalJSON(erc20events.Event{Origin: callerID, Payload: erc20events.Payload{From: "", To: minterID, Amount: transferAmount}})
	return stub.SetEvent(erc20events.TRANSFER, json)
}
