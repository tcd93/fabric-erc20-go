package helpers

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
)

/*In Fabric, the cert's CN is the enrollmentID of user, so this can be used as a unique identifier within an MSP.
Each MspID must be unique within the system, so the combination of MspID and CNs produces an unique identifier
for very single identity that participates the Hyperledger blockchain system.

https://hyperledger-fabric-ca.readthedocs.io/en/release-1.4/users-guide.html#cafiles*/

/*GetCallerID returns the unique ID (inside ledger) of chaicode caller.

The ID is a key of: `[mspID],[SubjectCN],[IssuerCN]`; this ID is used as "key" in the world-state database to identity a token owner.

- `mspID`: the MspID of the calling org

- `IssuerCN`: the common name from the calling identity's certificate CA provider

- `SubjectCN`: the common name from the calling identity's x509 certificate*/
func GetCallerID(stub shim.ChaincodeStubInterface) (string, error) {
	callerCert, err := cid.GetX509Certificate(stub)
	if err != nil {
		return "", err
	}
	orgMspID, err := cid.GetMSPID(stub)
	if err != nil {
		return "", err
	}
	return orgMspID + "," + callerCert.Issuer.CommonName + "," + callerCert.Subject.CommonName, nil
}
