package main

import (
	. "erc20/helpers"
	"erc20/lib/erc20basic"
	"erc20/lib/erc20burnable"
	"erc20/lib/erc20detailed"
	"erc20/lib/erc20mintable"
	"erc20/lib/erc20ownable"
	"erc20/lib/erc20pausable"
	"fmt"
	"math"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

//InitialMintAmount * 10^(token decimals) is the initial `total supply` of tokens
const InitialMintAmount float64 = 1000000000

var logger = shim.NewLogger("token-logger")

/*SampleToken is a simple ERC20 Token example. Refer to https://eips.ethereum.org/EIPS/eip-20 for documentations.*/
type SampleToken struct {
	erc20basic.BasicTokenInterface
	erc20ownable.OwnableTokenInterface
	erc20detailed.DetailedTokenInterface
	erc20mintable.MintableTokenInterface
	erc20burnable.BurnableTokenInterface
	erc20pausable.PausableTokenInterface
}

// main function starts up the chaincode in the container during instantiate
func main() {
	//new instance of token that mostly implements standard library
	//erc20 basic type is extended with "memo" functionality
	sampleToken := &SampleToken{
		&CustomBasicToken{},
		&erc20ownable.Token{},
		&erc20detailed.Token{},
		&erc20mintable.Token{},
		&erc20burnable.Token{},
		&erc20pausable.Token{},
	}
	if err := shim.Start(sampleToken); err != nil {
		panic(err)
	}
}

//#region chain code implementation

/*Init chaincode for Token, this method is called when we instantiate or upgrade our token.
(https://hyperledger-fabric.readthedocs.io/en/release-1.4/chaincode4ade.html#initializing-the-chaincode).

Init takes in one argument as a JSON-formatted string for token configurations, specifies the token attributes.
Owner of the token is also initialized as the contract's invoker.

Examples: `{"name": "tokenName", "symbol": "tokenSymbol", "decimals": "18"}`*/
func (t *SampleToken) Init(stub shim.ChaincodeStubInterface) peer.Response {
	callerID, err := GetCallerID(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	// if this is not the first init call (chaincode upgrade)
	// then owner validation is needed
	if currentOwner, _ := t.GetOwner(stub); strings.TrimSpace(currentOwner) != "" {
		logger.Infof("Upgrading chaincode using %v...", callerID)
		if err := CheckCallerIsOwner(callerID, currentOwner); err != nil {
			return shim.Error(err.Error())
		}
	} else {
		logger.Infof("Init chaincode using %v...", callerID)
		// if this is first call, then initialize states
		args := stub.GetStringArgs()
		if err := CheckArgsLength(args, 1); err != nil {
			return shim.Error(err.Error())
		}

		coinConfig := JSONToMap(args[0])

		// checks if "decimals" is a string of number format
		n := StringToInt(coinConfig["decimals"].(string))

		err = stub.PutState("owner", []byte(callerID))
		if err != nil {
			return shim.Error(err.Error())
		}
		err = stub.PutState("name", []byte(coinConfig["name"].(string)))
		if err != nil {
			return shim.Error(err.Error())
		}
		err = stub.PutState("symbol", []byte(coinConfig["symbol"].(string)))
		if err != nil {
			return shim.Error(err.Error())
		}
		err = stub.PutState("decimals", []byte(coinConfig["decimals"].(string)))
		if err != nil {
			return shim.Error(err.Error())
		}
		//mint the initial total supply
		//https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/examples/SimpleToken.sol
		err = t.Mint(stub,
			[]string{callerID, FloatToString(InitialMintAmount * math.Pow10(n))},
			withOwnerIs(callerID),
			t.GetBalanceOf,
			t.GetTotalSupply,
		)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success(nil)
}

func withOwnerIs(owner string) func(shim.ChaincodeStubInterface) (string, error) {
	return func(shim.ChaincodeStubInterface) (string, error) {
		return owner, nil
	}
}

/*Invoke is called per transaction on the chaincode*/
func (t *SampleToken) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	methodName, params := stub.GetFunctionAndParameters()

	//some functions are locked when the token state is "paused"
	isPaused, err := t.IsPaused(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	if isPaused {
		switch methodName {
		case "Transfer", "TransferFrom", "UpdateApproval":
			return shim.Error("Calling " + methodName + " is not allowed when token is paused")
		}
	}

	switch methodName {
	case "GetBalanceOf":
		f, err := t.GetBalanceOf(stub, params)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(FloatToBuffer(f))
	case "GetTotalSupply":
		f, err := t.GetTotalSupply(stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(FloatToBuffer(f))
	case "GetAllowance":
		f, err := t.GetAllowance(stub, params)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(FloatToBuffer(f))
	case "GetOwner":
		s, err := t.GetOwner(stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte(s))
	case "TransferOwnership":
		err := t.TransferOwnership(stub, params, t.GetOwner)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	case "GetName":
		s, err := t.GetName(stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte(s))
	case "GetSymbol":
		s, err := t.GetSymbol(stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte(s))
	case "GetDecimals":
		f, err := t.GetDecimals(stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(FloatToBuffer(f))
	case "Mint":
		err := t.Mint(stub, params, t.GetOwner, t.GetBalanceOf, t.GetTotalSupply)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	case "Burn":
		err := t.Burn(stub, params, t.GetTotalSupply, t.GetBalanceOf)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	case "BurnFrom":
		err := t.BurnFrom(stub, params, t.GetAllowance, t.GetTotalSupply, t.GetBalanceOf)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	case "Transfer":
		err := t.Transfer(stub, params, t.GetBalanceOf)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	case "TransferFrom":
		err := t.TransferFrom(stub, params, t.GetBalanceOf, t.GetAllowance)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	case "UpdateApproval":
		err := t.UpdateApproval(stub, params)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	case "Pause":
		err := t.Pause(stub, t.GetOwner)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	case "Unpause":
		err := t.Unpause(stub, t.GetOwner)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	case "GetMemo":
		s, err := t.GetMemo(stub, params)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte(s))
	}

	return shim.Error("Input function is not defined in chaincode")
}

//#endregion chain code implementation

//#region custom non-standard ERC20 implementation (transaction memo)
var customLogger = shim.NewLogger("memo-logger")

/*CustomBasicToken adds "memo" feature on top of basic ERC20 implementation*/
type CustomBasicToken struct {
	parentToken *erc20basic.Token
}

/*GetBalanceOf reimplement erc20basic's GetBalanceOf method*/
func (t *CustomBasicToken) GetBalanceOf(stub shim.ChaincodeStubInterface, args []string) (float64, error) {
	return t.parentToken.GetBalanceOf(stub, args)
}

/*GetTotalSupply reimplement erc20basic's GetTotalSupply method*/
func (t *CustomBasicToken) GetTotalSupply(stub shim.ChaincodeStubInterface) (float64, error) {
	return t.parentToken.GetTotalSupply(stub)
}

/*GetAllowance reimplement erc20basic's GetAllowance method*/
func (t *CustomBasicToken) GetAllowance(stub shim.ChaincodeStubInterface, args []string) (float64, error) {
	return t.parentToken.GetAllowance(stub, args)
}

/*UpdateApproval reimplement erc20basic's UpdateApproval method*/
func (t *CustomBasicToken) UpdateApproval(stub shim.ChaincodeStubInterface, args []string) error {
	return t.parentToken.UpdateApproval(stub, args)
}

/*GetMemo is a customed non standard erc20 that return the last memo string attached with transaction.

* `args[0]` - the key ID of target client.*/
func (t *SampleToken) GetMemo(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if err := CheckArgsLength(args, 1); err != nil {
		return "", err
	}

	//expect only one element in iterator
	iterator, err := stub.GetStateByPartialCompositeKey("Memo", args)
	defer iterator.Close()

	if err != nil {
		customLogger.Errorf("[sample-token.GetMemo] error after GetStateByPartialCompositeKey: %v", err)
		return "", err
	}
	if iterator.HasNext() {
		customLogger.Infof("Getting last memo from composite key %v", args[0])
		queryResult, err := iterator.Next()
		return string(queryResult.GetValue()), err
	}
	customLogger.Warningf("Memo not found for ID %v", args[0])
	return "", fmt.Errorf("Memo not found for ID %v", args[0])
}

/*Transfer adds "memo" feature after erc20basic's Transfer method*/
func (t *CustomBasicToken) Transfer(stub shim.ChaincodeStubInterface, args []string, getBalanceOf func(shim.ChaincodeStubInterface, []string) (float64, error)) error {
	err := t.parentToken.Transfer(stub, args, getBalanceOf)
	if err != nil {
		return err
	}
	//assumes the 3rd element in `args` is the comment
	if len(args) == 3 {
		receiverID, _, comment := args[0], args[1], args[2]
		return setMemo(stub, receiverID, comment)
	}
	return nil
}

/*TransferFrom adds "memo" feature after erc20basic's TransferFrom method*/
func (t *CustomBasicToken) TransferFrom(stub shim.ChaincodeStubInterface,
	args []string,
	getBalanceOf func(shim.ChaincodeStubInterface, []string) (float64, error),
	getAllowance func(shim.ChaincodeStubInterface, []string) (float64, error),
) error {
	err := t.parentToken.TransferFrom(stub, args, getBalanceOf, getAllowance)
	if err != nil {
		return err
	}
	//assumes the 4th element in `args` is the comment
	if len(args) == 4 {
		_, receiverID, _, comment := args[0], args[1], args[2], args[3]
		return setMemo(stub, receiverID, comment)
	}
	return nil
}

//setMemo updates world-state with a composite key of objectType "Memo", attribute of `key` and value of `memo`
func setMemo(stub shim.ChaincodeStubInterface, key string, memo string) error {
	memoKey, err := stub.CreateCompositeKey("Memo", []string{key})
	if err != nil {
		return err
	}
	customLogger.Infof("setting %v to memo %v", memoKey, memo)
	return stub.PutState(memoKey, []byte(memo))
}

//#endregion custom non-standard ERC20 implementation (transaction memo)
