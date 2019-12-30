package erc20ownable

import (
	. "erc20/helpers"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/*Token ownable implements OwnableTokenInterface*/
type Token struct{}

/*GetOwner returns the owner ID of token*/
func (t *Token) GetOwner(stub shim.ChaincodeStubInterface) (string, error) {
	owner, err := stub.GetState("owner")
	return string(owner), err
}

/*TransferOwnership allows the current owner to transfer control of the contract to a new owner.

* `args[0]` - the ID of the new owner.

* `getOwner` - specifies the function of getting the current owner of token.*/
func (t *Token) TransferOwnership(stub shim.ChaincodeStubInterface,
	args []string,
	getOwner func(shim.ChaincodeStubInterface) (string, error),
) error {
	callerID, err := GetCallerID(stub)
	if err != nil {
		return err
	}

	tokenOwnerID, err := getOwner(stub)
	if err := CheckCallerIsOwner(callerID, tokenOwnerID); err != nil {
		return err
	}

	newOwnerID := args[0]

	err = stub.PutState("owner", []byte(newOwnerID))
	return err
}
