package erc20burnable

import (
	"math/big"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/*BurnableTokenInterface consists of Burn & BurnFrom*/
type BurnableTokenInterface interface {
	Burn(stub shim.ChaincodeStubInterface,
		args []string,
		getTotalSupply func(stub shim.ChaincodeStubInterface) (*big.Int, error),
		getBalanceOf func(stub shim.ChaincodeStubInterface, args []string) (*big.Int, error),
	) error

	BurnFrom(stub shim.ChaincodeStubInterface,
		args []string,
		getAllowance func(stub shim.ChaincodeStubInterface, args []string) (*big.Int, error),
		getTotalSupply func(stub shim.ChaincodeStubInterface) (*big.Int, error),
		getBalanceOf func(stub shim.ChaincodeStubInterface, args []string) (*big.Int, error),
	) error
}
