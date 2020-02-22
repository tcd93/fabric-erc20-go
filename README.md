# Sample ERC20 token implementation with Hyperledger Fabric chaincode (Golang)

This is the Golang chaincode version of [Ethereum's ERC20 token standard](https://eips.ethereum.org/EIPS/eip-20) 
based [OpenZeppelin's implementation](https://github.com/OpenZeppelin/openzeppelin-contracts/blob/9b3710465583284b8c4c5d2245749246bb2e0094/contracts/token/ERC20/ERC20.sol)

_Basic features:_
* Basic Token
* Ownable Token
* Detailed Token
* Mintable Token
* Burnable Token
* Pausable Token

_Custom feature:_
* Transaction memo (an extension from Basic Token feature)
* Unregistered account check - accounts that are not registered via Fabric Node-sdk can not do transactions   
---

## Offline Unit Testing
*Manually* apply this patch for [ABAC testing with mockStub](https://gerrit.hyperledger.org/r/c/fabric/+/28744/2/core/chaincode/shim/mockstub.go#361) (_change line 69 & 361 like the patch in file `loyalty-token-hf\hlt\go_workspace\src\github.com\hyperledger\fabric\core\chaincode\shim\mockstub.go`_), after you've installed Go chaincode libraries of course

Get [Ginkgo test framework](https://onsi.github.io/ginkgo/)
```
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega/...
```
Execute `go test -v` to run the test suites without the need to start up the block chain

_note: if there are problems with dependencies during development with govendor, try renaming the "vendor" folder to "\_vendor", and rename it back to "vendor" when deploying chaincode._   
---

For details please read into `sample_token.go`