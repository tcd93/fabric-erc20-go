package main_test

import (
	. "erc20"
	"erc20/lib/erc20burnable"
	"erc20/lib/erc20detailed"
	"erc20/lib/erc20mintable"
	"erc20/lib/erc20ownable"
	"erc20/lib/erc20pausable"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Account activation", func() {
	const (
		txID          = `test-activation-id`
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

	// var err error
	var mockStub *shim.MockStub = shim.NewMockStub("mockStubNormal", &sampleToken)
	// var initialTotalSupply *big.Int = Mul(big.NewInt(InitialMintAmount), Pow(10, StringToInt(tokenDecimals)))

	// ownerID := ownerOrg + "," + ownerIssuer + "," + ownerSubject
	// fromID := fromOrg + "," + issuer + "," + fromSubject
	// toID := toOrg + "," + issuer + "," + toSubject

	BeforeEach(func() {
		mockStub.MockTransactionStart(txID)
	})

	AfterEach(func() {
		mockStub.MockTransactionEnd(txID)
	})

	It("Should return error when getting balance of non-existent account", func() {
		balance, err := sampleToken.GetBalanceOf(mockStub, []string{"fake-account"})
		Expect(err.Error()).To(ContainSubstring("is not registered"))
		Expect(balance).To(BeNil())
	})

	It("Should activate the fake account", func() {
		err := sampleToken.Activate(mockStub, []string{"fake-account"}, sampleToken.GetBalanceOf)
		Expect(err).To(BeNil())
	})

	It("Should now be able to get balance of fake account", func() {
		balance, err := sampleToken.GetBalanceOf(mockStub, []string{"fake-account"})
		Expect(err).To(BeNil())
		Expect(balance.String()).To(Equal("0"))
	})
})
