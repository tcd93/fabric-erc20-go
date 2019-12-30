package erc20basic

import (
	. "erc20/helpers"
	"erc20/lib/erc20events"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("trans-logger")

/*Token basic implementation of BasicTokenInterface*/
type Token struct{}

/*GetBalanceOf sender by ID.

* `args[0]` - the ID of user.*/
func (t *Token) GetBalanceOf(stub shim.ChaincodeStubInterface, args []string) (float64, error) {
	tokenBalance, err := stub.GetState(args[0])
	logger.Infof("GetBalanceOf: getting balance of %v...", args[0])
	return BufferToFloat(DefaultToZeroIfEmpty(tokenBalance)), err
}

/*GetTotalSupply returns total number of tokens in existence*/
func (t *Token) GetTotalSupply(stub shim.ChaincodeStubInterface) (float64, error) {
	totalSupply, err := stub.GetState("totalSupply")
	return BufferToFloat(DefaultToZeroIfEmpty(totalSupply)), err
}

/*GetAllowance checks the amount of tokens that an owner allowed a spender to transfer in behalf of the owner to another receiver.

* `args[0]` - the ID of owner.

* `args[1]` - the ID of spender*/
func (t *Token) GetAllowance(stub shim.ChaincodeStubInterface, args []string) (float64, error) {
	ownerID, spenderID := args[0], args[1]

	allowance, err := stub.GetState(ownerID + "-" + spenderID)

	return BufferToFloat(DefaultToZeroIfEmpty(allowance)), err
}

/*Transfer token from current caller to a specified address.

* `args[0]` - the ID of receiver.

* `args[1]` - the transfer amount.

* `getBalanceOf` - specifies the function of getting the initial balances of token sender & receiver.*/
func (t *Token) Transfer(stub shim.ChaincodeStubInterface, args []string, getBalanceOf func(shim.ChaincodeStubInterface, []string) (float64, error)) error {
	receiverID, sValue := args[0], args[1]

	senderID, err := GetCallerID(stub)
	if err != nil {
		return err
	}

	if err := CheckGreaterThanZero(sValue); err != nil {
		return err
	}

	transferAmount := StringToFloat(sValue)

	logger.Infof("Tranfer: transferring %v tokens from %v to %v...", transferAmount, senderID, receiverID)

	balanceOfSender, err := getBalanceOf(stub, []string{senderID})
	if err != nil {
		return err
	}
	balanceOfReceiver, err := getBalanceOf(stub, []string{receiverID})
	if err != nil {
		return err
	}

	if err := CheckBalance(balanceOfSender, senderID); err != nil {
		return err
	}
	if err := CheckBalance(balanceOfReceiver, receiverID); err != nil {
		return err
	}
	if err := IsSmallerOrEqual(transferAmount, balanceOfSender); err != nil {
		return fmt.Errorf("transfer amount should be less than balance of sender (%v): %v", senderID, err)
	}

	err = stub.PutState(senderID, FloatToBuffer(balanceOfSender-transferAmount))
	if err != nil {
		return err
	}
	err = stub.PutState(receiverID, FloatToBuffer(balanceOfReceiver+transferAmount))
	if err != nil {
		return err
	}

	json := MalshalJSON(erc20events.Event{Origin: senderID, Payload: erc20events.Payload{From: senderID, To: receiverID, Amount: transferAmount}})
	return stub.SetEvent(erc20events.TRANSFER, json)
}

/*TransferFrom transfer tokens from token owner to receiver.

* `args[1]` - the ID of token owner.

* `args[1]` - the ID of receiver.

* `args[2]` - the transfer amount.

* `getBalanceOf` - defines the function of getting the initial balances of token owner & receiver.

* `getAllowance` - defines the function of getting the allowance that the current chaincode invoker can spend from the token owner, to transfer to the receiver.*/
func (t *Token) TransferFrom(stub shim.ChaincodeStubInterface,
	args []string,
	getBalanceOf func(shim.ChaincodeStubInterface, []string) (float64, error),
	getAllowance func(shim.ChaincodeStubInterface, []string) (float64, error),
) error {
	tokenOwnerID, receiverID, sValue := args[0], args[1], args[2]

	if err := CheckGreaterThanZero(sValue); err != nil {
		return err
	}

	spenderID, err := GetCallerID(stub)
	if err != nil {
		return err
	}

	transferAmount := StringToFloat(sValue)

	logger.Infof("TranferFrom: transferring %v tokens from %v to %v with %v...", transferAmount, tokenOwnerID, receiverID, spenderID)

	balanceOfTokenOwner, err := getBalanceOf(stub, []string{tokenOwnerID})
	if err != nil {
		return err
	}
	balanceOfReceiver, err := getBalanceOf(stub, []string{receiverID})
	if err != nil {
		return err
	}
	approvedAmount, err := getAllowance(stub, []string{tokenOwnerID, spenderID})
	if err != nil {
		return err
	}

	if err := CheckBalance(balanceOfTokenOwner, tokenOwnerID); err != nil {
		return err
	}
	if err := CheckBalance(balanceOfReceiver, receiverID); err != nil {
		return err
	}
	if err := CheckApproved(approvedAmount, tokenOwnerID+"-"+spenderID); err != nil {
		return err
	}
	if err := IsSmallerOrEqual(transferAmount, balanceOfTokenOwner); err != nil {
		return fmt.Errorf("transfer amount should be less than balance of token owner (%v): %v", tokenOwnerID, err)
	}
	if err := IsSmallerOrEqual(transferAmount, approvedAmount); err != nil {
		return fmt.Errorf("transfer amount should be less than approved spending amount of %v: %v", spenderID, err)
	}

	err = stub.PutState(tokenOwnerID, FloatToBuffer(balanceOfTokenOwner-transferAmount))
	if err != nil {
		return err
	}
	err = stub.PutState(tokenOwnerID+"-"+spenderID, FloatToBuffer(approvedAmount-transferAmount))
	if err != nil {
		return err
	}
	err = stub.PutState(receiverID, FloatToBuffer(balanceOfReceiver+transferAmount))
	if err != nil {
		return err
	}

	json := MalshalJSON(erc20events.Event{Origin: spenderID, Payload: erc20events.Payload{From: tokenOwnerID, To: receiverID, Amount: transferAmount}})
	return stub.SetEvent(erc20events.TRANSFER, json)
}

/*UpdateApproval approves the passed-in identity to spend/burn a maximum amount of tokens on behalf of the function caller.

* `args[0]` - the ID of approved user.

* `args[1]` - the maximum approved amount.*/
func (t *Token) UpdateApproval(stub shim.ChaincodeStubInterface, args []string) error {
	spenderID, newAllowance := args[0], args[1]

	if err := CheckGreaterThanZero(newAllowance); err != nil {
		return err
	}

	callerID, err := GetCallerID(stub)
	if err != nil {
		return err
	}

	err = stub.PutState(callerID+"-"+spenderID, []byte(newAllowance))
	if err != nil {
		return err
	}

	approvedAmount := StringToFloat(newAllowance)
	json := MalshalJSON(erc20events.Event{Origin: callerID, Payload: erc20events.Payload{From: callerID, To: spenderID, Amount: approvedAmount}})
	return stub.SetEvent(erc20events.APPROVAL, json)
}
