package erc20basic

import "github.com/hyperledger/fabric/core/chaincode/shim"

/*BasicTokenInterface consists of basic ERC20 methods*/
type BasicTokenInterface interface {
	GetBalanceOf(stub shim.ChaincodeStubInterface, args []string) (float64, error)

	GetTotalSupply(stub shim.ChaincodeStubInterface) (float64, error)

	GetAllowance(stub shim.ChaincodeStubInterface, args []string) (float64, error)

	Transfer(stub shim.ChaincodeStubInterface, args []string, getBalanceOf func(shim.ChaincodeStubInterface, []string) (float64, error)) error

	TransferFrom(stub shim.ChaincodeStubInterface,
		args []string,
		getBalanceOf func(shim.ChaincodeStubInterface, []string) (float64, error),
		getAllowance func(shim.ChaincodeStubInterface, []string) (float64, error),
	) error

	UpdateApproval(stub shim.ChaincodeStubInterface, args []string) error
}
