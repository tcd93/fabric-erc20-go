package erc20detailed

import (
	. "erc20/helpers"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/*Token detailed implementation of DetailedTokenInterface*/
type Token struct{}

/*GetName returns the name of the token*/
func (t *Token) GetName(stub shim.ChaincodeStubInterface) (string, error) {
	tokenNameBytes, err := stub.GetState("name")
	return string(tokenNameBytes), err
}

/*GetSymbol returns the symbol of the token*/
func (t *Token) GetSymbol(stub shim.ChaincodeStubInterface) (string, error) {
	tokenSymbolBytes, err := stub.GetState("symbol")
	return string(tokenSymbolBytes), err
}

/*GetDecimals returns the "decimals" value of the token, default to 0 if not set in the initial stage.

The decimals are only for visualization purposes.
All the operations are done using the smallest and indivisible token unit,
just as on Ethereum all the operations are done in wei.*/
func (t *Token) GetDecimals(stub shim.ChaincodeStubInterface) (float64, error) {
	tokenDecimalsBytes, err := stub.GetState("decimals")
	return BufferToFloat(DefaultToZeroIfEmpty(tokenDecimalsBytes)), err
}
