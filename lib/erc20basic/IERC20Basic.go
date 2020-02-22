package erc20basic

import (
	"math/big"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/*BasicTokenInterface consists of basic ERC20 methods*/
type BasicTokenInterface interface {
	GetBalanceOf(stub shim.ChaincodeStubInterface, args []string) (*big.Int, error)

	GetTotalSupply(stub shim.ChaincodeStubInterface) (*big.Int, error)

	GetAllowance(stub shim.ChaincodeStubInterface, args []string) (*big.Int, error)

	Transfer(stub shim.ChaincodeStubInterface, args []string, getBalanceOf func(shim.ChaincodeStubInterface, []string) (*big.Int, error)) error

	TransferFrom(stub shim.ChaincodeStubInterface,
		args []string,
		getBalanceOf func(shim.ChaincodeStubInterface, []string) (*big.Int, error),
		getAllowance func(shim.ChaincodeStubInterface, []string) (*big.Int, error),
	) error

	UpdateApproval(stub shim.ChaincodeStubInterface, args []string) error
}
