package erc20burnable

import "github.com/hyperledger/fabric/core/chaincode/shim"

/*BurnableTokenInterface consists of Burn & BurnFrom*/
type BurnableTokenInterface interface {
	Burn(stub shim.ChaincodeStubInterface,
		args []string,
		getTotalSupply func(stub shim.ChaincodeStubInterface) (float64, error),
		getBalanceOf func(stub shim.ChaincodeStubInterface, args []string) (float64, error),
	) error

	BurnFrom(stub shim.ChaincodeStubInterface,
		args []string,
		getAllowance func(stub shim.ChaincodeStubInterface, args []string) (float64, error),
		getTotalSupply func(stub shim.ChaincodeStubInterface) (float64, error),
		getBalanceOf func(stub shim.ChaincodeStubInterface, args []string) (float64, error),
	) error
}
