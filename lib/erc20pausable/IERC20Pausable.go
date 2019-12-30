package erc20pausable

import "github.com/hyperledger/fabric/core/chaincode/shim"

/*PausableTokenInterface consists of Pause & Unpause (methods should be restricted), IsPaused to check state*/
type PausableTokenInterface interface {
	IsPaused(stub shim.ChaincodeStubInterface) (bool, error)

	Pause(stub shim.ChaincodeStubInterface,
		getOwner func(shim.ChaincodeStubInterface) (string, error),
	) error

	Unpause(stub shim.ChaincodeStubInterface,
		getOwner func(shim.ChaincodeStubInterface) (string, error),
	) error
}
