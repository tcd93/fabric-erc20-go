package erc20ownable

import "github.com/hyperledger/fabric/core/chaincode/shim"

/*OwnableTokenInterface consists of GetOwner & TransferOwnership*/
type OwnableTokenInterface interface {
	GetOwner(stub shim.ChaincodeStubInterface) (string, error)

	TransferOwnership(stub shim.ChaincodeStubInterface,
		args []string,
		getOwner func(shim.ChaincodeStubInterface) (string, error),
	) error
}
