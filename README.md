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
* **Transaction memo** - able to attach an 'memo' to a transaction with a extra parameter to `Transfer` or `TransferFrom` methods
* **Unregistered account check** - accounts that are not registered can not do transactions, register them first with `Activate` chaincode method
---
## Demo
Set up the network via development tool ([Hurley](https://github.com/worldsibu/hurley)) or manual set up via the [official document](https://hyperledger-fabric.readthedocs.io/en/release-1.4/dev-setup/devenv.html)

[Install & instantiate the chaincode](https://hyperledger-fabric.readthedocs.io/en/release-1.4/chaincode4noah.html#installing-chaincode) on the network

![](/readme/env.png) *chaincodes are placed inside a separated Docker container*

![](/readme/1.png) 
![](/readme/2.png)
![](/readme/3.png)*total supply of tokens are equal to the miner's (chaincode caller) account*

![](/readme/4.png)*try to transfer some tokens to another account --> error due to the target user is not 'activated'*

![](/readme/5.png)*'activate' said account*

![](/readme/6.png)*`transfer` success!*

![](/readme/7.png)

---

## Offline Unit Testing
*Manually* apply this patch for [ABAC testing with mockStub](https://gerrit.hyperledger.org/r/c/fabric/+/28744/2/core/chaincode/shim/mockstub.go#361) (_change line 69 & 361 like the patch in file `loyalty-token-hf\hlt\go_workspace\src\github.com\hyperledger\fabric\core\chaincode\shim\mockstub.go`_), after you've installed Go chaincode libraries of course

Get [Ginkgo test framework](https://onsi.github.io/ginkgo/)
```
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega/...
```
Execute `go test -v` to run the test suites without the need to start up the block chain

---

For details please read into `sample_token.go`
