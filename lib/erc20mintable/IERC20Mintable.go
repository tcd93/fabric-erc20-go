package erc20mintable

import (
	"math/big"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/*MintableTokenInterface consists of Mint, this method should be very restricted*/
type MintableTokenInterface interface {
	Mint(stub shim.ChaincodeStubInterface,
		args []string,
		getOwner func(shim.ChaincodeStubInterface) (string, error),
		getBalanceOf func(shim.ChaincodeStubInterface, []string) (*big.Int, error),
		getTotalSupply func(shim.ChaincodeStubInterface) (*big.Int, error),
	) error
}
