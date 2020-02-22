package main_test

import (
	. "erc20"
	. "erc20/helpers"
	"erc20/lib/erc20burnable"
	"erc20/lib/erc20detailed"
	"erc20/lib/erc20mintable"
	"erc20/lib/erc20ownable"
	"erc20/lib/erc20pausable"
	. "erc20/testutils"
	"fmt"
	"math/big"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("token - normal cert without attrs", func() {
	const (
		txID          = `test-transaction-id`
		tokenName     = `sample token name`
		tokenSymbol   = `(y)(y)`
		tokenDecimals = `14`

		//attributes matches the certs in /testutils
		issuer      = `Org1-child1`
		fromOrg     = `clientOrg1MSP`
		fromSubject = `Org1-child1-client1`

		toOrg     = `clientOrg2MSP`
		toSubject = `Org1-child1-client2`

		ownerIssuer  = `Org1`
		ownerOrg     = `sampleOrgMSP`
		ownerSubject = `Org1-child1`
	)

	sampleToken := SampleToken{
		&CustomBasicToken{},
		&erc20ownable.Token{},
		&erc20detailed.Token{},
		&erc20mintable.Token{},
		&erc20burnable.Token{},
		&erc20pausable.Token{},
	}

	var err error
	var mockStub *shim.MockStub = shim.NewMockStub("mockStubNormal", &sampleToken)
	var initialTotalSupply *big.Int = Mul(big.NewInt(InitialMintAmount), Pow(10, StringToInt(tokenDecimals)))

	ownerID := ownerOrg + "," + ownerIssuer + "," + ownerSubject
	fromID := fromOrg + "," + issuer + "," + fromSubject
	toID := toOrg + "," + issuer + "," + toSubject

	BeforeSuite(func() {
		mockStub, err = SetCurrentCaller(mockStub, ownerOrg, AdminCert)
		Expect(err).To(BeNil())
		Expect(mockStub.MockInit(
			txID,
			[][]byte{[]byte(
				fmt.Sprintf(
					`{"name": "%s", "symbol": "%s", "decimals": "%s"}`,
					tokenName, tokenSymbol, tokenDecimals,
				))},
		).Message).To(BeEmpty()) //OK status

		Describe("The owner ID", func() {
			caller, err := GetCallerID(mockStub)
			Expect(err).To(BeNil())
			Expect(caller).To(Equal(ownerID))
		})
	})

	BeforeEach(func() {
		mockStub.MockTransactionStart(txID)
	})

	AfterEach(func() {
		mockStub.MockTransactionEnd(txID)
	})

	Describe("Initial token attributes...", func() {
		It("Allow everyone to get token symbol", func() {
			symbol, err := sampleToken.GetSymbol(mockStub)
			Expect(err).To(BeNil())
			Expect(symbol).To(Equal(tokenSymbol))
		})

		It("Allow everyone to get token name", func() {
			name, err := sampleToken.GetName(mockStub)
			Expect(err).To(BeNil())
			Expect(name).To(Equal(tokenName))
		})

		It("Allow everyone to get token total supply", func() {
			totalSupply, err := sampleToken.GetTotalSupply(mockStub)
			Expect(err).To(BeNil())
			Expect(totalSupply).To(Equal(initialTotalSupply))

			Describe("Token owner's balance is the Token's total supply", func() {
				ownerBalance, err := sampleToken.GetBalanceOf(mockStub, []string{ownerID})
				Expect(err).To(BeNil())
				Expect(totalSupply).To(Equal(ownerBalance))
			})
		})

		It("Allow everyone to get token decimals", func() {
			decimals, err := sampleToken.GetDecimals(mockStub)
			Expect(err).To(BeNil())
			Expect(decimals).To(Equal(tokenDecimals))
		})
	})

	Describe("Mocked Token transfer functionalities...", func() {
		When("Chaincode invoker is NOT current owner", func() {
			It("Change the invoker to `Client1Cert`", func() {
				_, err = SetCurrentCaller(mockStub, fromOrg, Client1Cert)
				Expect(err).To(BeNil())
			})

			It("Expects chaincode caller to be `fromID`", func() {
				caller, err := GetCallerID(mockStub)
				Expect(err).To(BeNil())
				Expect(caller).To(Equal(fromID))
			})

			It("Expects this to be the same token", func() {
				totalSupply, err := sampleToken.GetTotalSupply(mockStub)
				Expect(err).To(BeNil())
				Expect(totalSupply).To(Equal(initialTotalSupply))
			})

			It("Transfer tokens between clients using Transfer", func() {
				transferAmount := big.NewInt(2)
				initialBalance := big.NewInt(5)

				err = sampleToken.Transfer(mockStub,
					[]string{toID, transferAmount.String()},
					WithBalanceOf(initialBalance),
				)
				Expect(err).To(BeNil())

				balance, err := sampleToken.GetBalanceOf(mockStub, []string{fromID})
				Expect(err).To(BeNil())
				Expect(balance).To(Equal(Sub(initialBalance, transferAmount)))

				balance, err = sampleToken.GetBalanceOf(mockStub, []string{toID})
				Expect(err).To(BeNil())
				Expect(balance).To(Equal(Add(initialBalance, transferAmount)))
			})
		})

		When("Chaincode invoker IS current owner", func() {
			It("Change the invoker back to token owner", func() {
				_, err = SetCurrentCaller(mockStub, ownerOrg, AdminCert)
				Expect(err).To(BeNil())
			})

			It("Expects chaincode caller to be `ownerID`", func() {
				caller, err := GetCallerID(mockStub)
				Expect(err).To(BeNil())
				Expect(caller).To(Equal(ownerID))
			})

			It("Expects this to be the same token", func() {
				totalSupply, err := sampleToken.GetTotalSupply(mockStub)
				Expect(err).To(BeNil())
				Expect(totalSupply).To(Equal(initialTotalSupply))
			})

			It("Transfer tokens between clients using TransferFrom", func() {
				transferAmount := big.NewInt(30)
				initialBalance := big.NewInt(50)
				allowance := big.NewInt(40)

				err := sampleToken.TransferFrom(mockStub,
					[]string{fromID, toID, transferAmount.String()},
					WithBalanceOf(initialBalance),
					WithAllowanceOf(allowance),
				)
				Expect(err).To(BeNil())

				balance, err := sampleToken.GetBalanceOf(mockStub, []string{fromID})
				Expect(err).To(BeNil())
				Expect(balance).To(Equal(Sub(initialBalance, transferAmount)))

				balance, err = sampleToken.GetBalanceOf(mockStub, []string{toID})
				Expect(err).To(BeNil())
				Expect(balance).To(Equal(Add(initialBalance, transferAmount)))
			})

			It("Should fails when transfer more tokens than allowed", func() {
				//this transaction should not be OK, as `allowance` > `transferAmount`
				transferAmount := big.NewInt(60)
				initialBalance := big.NewInt(100)
				allowance := big.NewInt(40)

				initialBalanceMockFunc := WithBalanceOf(initialBalance)
				allowanceMockFunc := WithAllowanceOf(allowance)

				err := sampleToken.TransferFrom(mockStub,
					[]string{fromID, toID, transferAmount.String()},
					initialBalanceMockFunc,
					allowanceMockFunc,
				)
				Expect(err).NotTo(BeNil())
			})

			It("Should fails when transfer tokens with empty balance", func() {
				//this transaction should not be OK, as `transferAmount` > `initialBalance`
				transferAmount := big.NewInt(60)
				initialBalance := big.NewInt(0)
				allowance := big.NewInt(100)

				initialBalanceMockFunc := WithBalanceOf(initialBalance)
				allowanceMockFunc := WithAllowanceOf(allowance)

				err := sampleToken.TransferFrom(mockStub,
					[]string{fromID, toID, transferAmount.String()},
					initialBalanceMockFunc,
					allowanceMockFunc,
				)
				Expect(err).NotTo(BeNil())
			})
		})
	})

	Describe("Integrated Token transfer functionalities...", func() {
		It("Transfer tokens from owner to client", func() {
			transferAmount := big.NewInt(110)

			initialBalanceOfClient, err := sampleToken.GetBalanceOf(mockStub, []string{toID})
			// Expect(initialBalanceOfClient).To(Equal(80.00))
			Expect(err).To(BeNil())
			initialBalanceOfOwner, err := sampleToken.GetBalanceOf(mockStub, []string{ownerID})
			Expect(err).To(BeNil())

			err = sampleToken.Transfer(mockStub,
				[]string{toID, transferAmount.String()},
				sampleToken.GetBalanceOf,
			)
			Expect(err).To(BeNil())

			balance, err := sampleToken.GetBalanceOf(mockStub, []string{toID})
			Expect(err).To(BeNil())
			Expect(balance).To(Equal(Add(initialBalanceOfClient, transferAmount)))

			balance, err = sampleToken.GetBalanceOf(mockStub, []string{ownerID})
			Expect(err).To(BeNil())
			Expect(balance).To(Equal(Sub(initialBalanceOfOwner, transferAmount)))
		})

		Describe("Transfer tokens with Comment", func() {
			comment := "Transfer 10 token"

			It("Should transfer with memo attached", func() {
				transferAmount := big.NewInt(10)

				err = sampleToken.Transfer(mockStub,
					[]string{toID, transferAmount.String(), comment},
					sampleToken.GetBalanceOf,
				)
				Expect(err).To(BeNil())
			})

			It("Should get the memo of last transaction", func() {
				memo, err := sampleToken.GetMemo(mockStub, []string{toID})
				Expect(err).To(BeNil())
				Expect(memo).To(Equal(comment))
			})

			It("Should be empty memo for wrongful key", func() {
				memo, err := sampleToken.GetMemo(mockStub, []string{"notACompositeKey" + toID})
				Expect(err).NotTo(BeNil())
				Expect(memo).To(BeEmpty())
			})
		})
	})

	Describe("Integrated Token burn functionalities...", func() {
		Describe("Burn some tokens of owner by owner itself", func() {
			It("Should burn successfully", func() {
				burnAmount := big.NewInt(50)

				initialBalanceOfOwner, err := sampleToken.GetBalanceOf(mockStub, []string{ownerID})
				Expect(err).To(BeNil())
				initialTotalSupply, err := sampleToken.GetTotalSupply(mockStub)
				Expect(err).To(BeNil())

				err = sampleToken.Burn(mockStub,
					[]string{burnAmount.String()},
					sampleToken.GetTotalSupply,
					sampleToken.GetBalanceOf,
				)
				Expect(err).To(BeNil())

				balance, err := sampleToken.GetBalanceOf(mockStub, []string{ownerID})
				Expect(err).To(BeNil())
				Expect(balance).To(Equal(Sub(initialBalanceOfOwner, burnAmount)))

				currentTotalSupply, err := sampleToken.GetTotalSupply(mockStub)
				Expect(err).To(BeNil())
				Expect(currentTotalSupply).To(Equal(Sub(initialTotalSupply, burnAmount)))
			})
		})

		Describe("Burning some tokens of client by owner", func() {
			When("Allowance of owner is not updated", func() {
				It("Should fail", func() {
					burnAmount := big.NewInt(75)

					err = sampleToken.BurnFrom(mockStub,
						[]string{fromID, burnAmount.String()},
						sampleToken.GetAllowance,
						sampleToken.GetTotalSupply,
						sampleToken.GetBalanceOf,
					)
					Expect(err).NotTo(BeNil()) //owner does not have enough allowance from `fromID`
				})
			})

			When("Allowance of owner is added", func() {
				Context("By Client2", func() {
					It("Change the invoker to `Client2Cert`", func() {
						_, err = SetCurrentCaller(mockStub, toOrg, Client2Cert)
						Expect(err).To(BeNil())
					})

					It("Expects chaincode caller to be `toID`", func() {
						caller, err := GetCallerID(mockStub)
						Expect(err).To(BeNil())
						Expect(caller).To(Equal(toID))
					})

					It("Update allowance of owner by `toID`", func() {
						allowanceOfOwner := big.NewInt(100)
						err := sampleToken.UpdateApproval(mockStub, []string{ownerID, allowanceOfOwner.String()})
						Expect(err).To(BeNil())
					})
				})

				Context("By Owner", func() {
					It("Change the invoker back to token owner", func() {
						_, err = SetCurrentCaller(mockStub, ownerOrg, AdminCert)
						Expect(err).To(BeNil())
					})

					It("Expects chaincode caller to be `ownerID`", func() {
						caller, err := GetCallerID(mockStub)
						Expect(err).To(BeNil())
						Expect(caller).To(Equal(ownerID))
					})

					It("Should now be able to burn tokens of `toID` by owner", func() {
						burnAmount := big.NewInt(75)

						initialBalanceOfBurnee, err := sampleToken.GetBalanceOf(mockStub, []string{toID})
						Expect(err).To(BeNil())
						initialTotalSupply, err := sampleToken.GetTotalSupply(mockStub)
						Expect(err).To(BeNil())

						err = sampleToken.BurnFrom(mockStub,
							[]string{toID, burnAmount.String()},
							sampleToken.GetAllowance,
							sampleToken.GetTotalSupply,
							sampleToken.GetBalanceOf,
						)
						Expect(err).To(BeNil())

						balance, err := sampleToken.GetBalanceOf(mockStub, []string{toID})
						Expect(err).To(BeNil())
						Expect(balance).To(Equal(Sub(initialBalanceOfBurnee, burnAmount)))

						currentTotalSupply, err := sampleToken.GetTotalSupply(mockStub)
						Expect(err).To(BeNil())
						Expect(currentTotalSupply).To(Equal(Sub(initialTotalSupply, burnAmount)))
					})

					When("Burn again", func() {
						It("Should fail as `toID` has no more balance", func() {
							burnAmount := big.NewInt(126)

							err = sampleToken.BurnFrom(mockStub,
								[]string{toID, burnAmount.String()},
								sampleToken.GetAllowance,
								sampleToken.GetTotalSupply,
								sampleToken.GetBalanceOf,
							)
							Expect(err).NotTo(BeNil())
						})
					})
				})
			})
		})
	})

	Describe("Integrated Token pause functionalities...", func() {
		It("Pauses the token state", func() {
			err := sampleToken.Pause(mockStub, sampleToken.GetOwner)
			Expect(err).To(BeNil())
		})

		When("Token is paused", func() {
			It("Should not be able to transfer tokens", func() {
				transferAmount := big.NewInt(30)
				Expect(mockStub.MockInvoke(txID, [][]byte{[]byte("Transfer"), []byte(toID), []byte(transferAmount.String())}).Message).NotTo(BeEmpty())
			})
		})

		It("Unpauses the token state", func() {
			err := sampleToken.Unpause(mockStub, sampleToken.GetOwner)
			Expect(err).To(BeNil())
		})

		When("Token is unpaused", func() {
			It("Should now be able to transfer tokens", func() {
				transferAmount := big.NewInt(30)
				initialBalanceOfClient, err := sampleToken.GetBalanceOf(mockStub, []string{toID})

				Expect(mockStub.MockInvoke(txID, [][]byte{[]byte("Transfer"), []byte(toID), []byte(transferAmount.String())}).Message).To(BeEmpty())

				balance, err := sampleToken.GetBalanceOf(mockStub, []string{toID})
				Expect(err).To(BeNil())
				Expect(balance).To(Equal(Add(initialBalanceOfClient, transferAmount))) //125 + 30 = 155
			})
		})
	})
})
