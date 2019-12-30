package testutils

import (
	"encoding/base64"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
)

const (
	AdminCert = `-----BEGIN CERTIFICATE-----
MIICEDCCAbagAwIBAgIQUFdWsXAm1ui0HGC0pZKyCzAKBggqhkjOPQQDAjBYMQsw
CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy
YW5jaXNjbzENMAsGA1UEChMET3JnMTENMAsGA1UEAxMET3JnMTAeFw0xODA4MjEw
ODI1MzNaFw0yODA4MTgwODI1MzNaMGYxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpD
YWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMRQwEgYDVQQKEwtPcmcx
LWNoaWxkMTEUMBIGA1UEAxMLT3JnMS1jaGlsZDEwWTATBgcqhkjOPQIBBggqhkjO
PQMBBwNCAAQqi4WZz5mp1irSrcZ8VpgYR+3j9KvydJunAFjstUzy3KOwqmi6T8L8
06hBmgIbnrG3xKFHslRJHI4+utk0lQtUo1QwUjAOBgNVHQ8BAf8EBAMCAaYwDwYD
VR0lBAgwBgYEVR0lADAPBgNVHRMBAf8EBTADAQH/MA0GA1UdDgQGBAQBAgMEMA8G
A1UdIwQIMAaABAECAwQwCgYIKoZIzj0EAwIDSAAwRQIgff2sx62yXS8WZHLyvgDI
lDGkUfm1fwzxg6qiZ6q3oHICIQDj+6Id8GJj7VQeQrJxO6v7Z/AyQPW1sRW8OjxO
VdJ1IQ==
-----END CERTIFICATE-----
`
	Client1Cert = `-----BEGIN CERTIFICATE-----
MIICHzCCAcagAwIBAgIQMspIqOt+/arNiGqJTnRmhTAKBggqhkjOPQQDAjBmMQsw
CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy
YW5jaXNjbzEUMBIGA1UEChMLT3JnMS1jaGlsZDExFDASBgNVBAMTC09yZzEtY2hp
bGQxMB4XDTE4MDgyMTA4MjUzM1oXDTI4MDgxODA4MjUzM1owdjELMAkGA1UEBhMC
VVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBGcmFuY2lzY28x
HDAaBgNVBAoTE09yZzEtY2hpbGQxLWNsaWVudDExHDAaBgNVBAMTE09yZzEtY2hp
bGQxLWNsaWVudDEwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAT9AveGhSGrUZEo
wiZuBUv9Ry7FEGvLwUkd3wQKpDD7LIZLLaQcOOccLBT9OnzJ1VD+0XRo7JbbVgqT
Zr7hmoefo0YwRDAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUHAwIw
DAYDVR0TAQH/BAIwADAPBgNVHSMECDAGgAQBAgMEMAoGCCqGSM49BAMCA0cAMEQC
IEYv+/eH0PtBhFvkOYZDh5SEvF2Q033CwJsKBqBhRFlTAiB7kpAISt/x4jJw6UEr
RS0fb772EnAaBEDXa7uwdIkgTw==
-----END CERTIFICATE-----
	
`
	Client2Cert = `-----BEGIN CERTIFICATE-----
MIICIDCCAcagAwIBAgIQLMR83E59jMzXwRK0DD06PDAKBggqhkjOPQQDAjBmMQsw
CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy
YW5jaXNjbzEUMBIGA1UEChMLT3JnMS1jaGlsZDExFDASBgNVBAMTC09yZzEtY2hp
bGQxMB4XDTE4MDgyMTA4MjUzM1oXDTI4MDgxODA4MjUzM1owdjELMAkGA1UEBhMC
VVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBGcmFuY2lzY28x
HDAaBgNVBAoTE09yZzEtY2hpbGQxLWNsaWVudDIxHDAaBgNVBAMTE09yZzEtY2hp
bGQxLWNsaWVudDIwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAATLs6+t/PTll+kk
Kh49SROhPz5x6mgukBbDglmcYZbFrW6wpnT5cjJKJ9wFP5cPBATZp+5OdisR8O3d
syMH2Jo5o0YwRDAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUHAwIw
DAYDVR0TAQH/BAIwADAPBgNVHSMECDAGgAQBAgMEMAoGCCqGSM49BAMCA0gAMEUC
IBz5+OtbtAJ/JBUs1MmB94YeMGcY5As6ww24fMAHVfuKAiEAoeLFHzqW18ICh/ug
b0ZeQh2/OLCPy1ECwPBEhuDa+f0=
-----END CERTIFICATE-----
`

	IdemixCred = `CiAGTPr9iYBj7gqHeU90RgXXklBhgIKbtVtUzMDnhJ9bCRIggb1mWLzNEO1ac+PfDNSpbCC02eY5OqRooSH1mnhlftEaQgoMaWRlbWl4TVNQSUQyEhBvcmcxLmRlcGFydG1lbnQxGiCEOl1qWZ+TRhDD3X6Cl59utnDnW8oio0y2vuJ1xHNMMiIOCgxpZGVtaXhNU1BJRDIq5gYKRAog1MgGWhcd/jgQpbpdO17LktSoelSsQJKfAmFhspaM6+QSICMcxQLs4JRPeSbyWG81KNepmLIi8C1AOyrgytJYMmKIEkQKICCwfD7vfGKqzGFEa7H7cBbR81kImbXcECJnDbj4QQNUEiB64Bi0jRahh18QuZzqnw6sksn8GBCi2sVsrjtLTKsvMRpECiBD9miof2CyCjHOr6s/JAiALRzjdogv0xQHEyqNAfIDABIgx56y7lUllGc0XtYsFIdq7CulDE55Re5xT1wvzRNhITUiIPSaozvr294lNGF3Wy5Yd7wlNW/IZBpBcXda/dQfGci9KiBbO2o2bWD1P4HOMfI//ebo8WrwTNgmPfmlqNBxzKuQMjIg5AwmQGEYnKN/pOVDFMjm/3a9hJDv9R2svI42aVBms0M6IHMSFIZ8j/yZH5nHtCkwpQMCuBFmI6krD2CfTjCiOUfoQiDO7cyRnCt9uEGIhQsBiwnSEXH+G9Il9qvfkUrAiZlbrkogv6dmb1xijfB3gsyVWxgfKlRNRtf78dMwjSf76jEnSrBSIJTkD7lSBwBepMFROxYneTHuG6JcSZpdoeOGqFl0drJWUiC8ndC2y9LsFJLKs2ddFqsFW7kNg+vROXuSLQdglSBffVog/eDzc90wTBEZu2T6LhWEbcP5oZ5TYdE/o+cOUfPgV4RiRAogBkz6/YmAY+4Kh3lPdEYF15JQYYCCm7VbVMzA54SfWwkSIIG9Zli8zRDtWnPj3wzUqWwgtNnmOTqkaKEh9Zp4ZX7RaiCXNeWrQz2UPkuAEZrt++TP/DbmAFF7cBQlYkb81jrn/nKIAQog/gwzULTJbCAoVg9XfCiROs4cU5oSv4Q80iYWtonAnvsSIE6mYFdzisBU21rhxjfYE7kk3Xjih9A1idJp7TSjfmorGiBwIEbnxUKjs3Z3DXUSTj5R78skdY1hWEjpCbSBvtwn/yIgBVTjvNOIwpBC7qZJKX6yn4tMvoCCGpiz4BKBEUqtBJt6ZzBlAjEAoBaHzX1HjvrnPMDXajqcLeHR5//AIIGDDcGQ+4GNqJu9Wawlw6Zs58Nnkpmh29ivAjBJNHeGNvX9sQb9lyzLAtCa5Il4xKNGGpGZ+uhQAjtNpRAZLtv2hgSqJAy0X6HwNXeAAQGKAQA=`
)

func SetCurrentCaller(stub *shim.MockStub, mspID string, cert string) (*shim.MockStub, error) {
	sid := &msp.SerializedIdentity{Mspid: mspID, IdBytes: []byte(cert)}
	b, err := proto.Marshal(sid)
	if err != nil {
		return nil, err
	}
	stub.Creator = b
	return stub, nil
}

func SetIdemixMockStubWithAttrs(stub *shim.MockStub, mspID string) (*shim.MockStub, error) {
	idBytes, err := base64.StdEncoding.DecodeString(IdemixCred)
	if err != nil {
		return nil, err
	}
	sid := &msp.SerializedIdentity{Mspid: mspID, IdBytes: idBytes}
	b, err := proto.Marshal(sid)
	if err != nil {
		return nil, err
	}
	stub.Creator = b
	return stub, nil
}

func SetMockStubWithNilCreator(stub *shim.MockStub) (*shim.MockStub, error) {
	stub.Creator = nil
	return stub, nil
}

func SetMockStubWithFakeCreator(stub *shim.MockStub) (*shim.MockStub, error) {
	stub.Creator = []byte("foo")
	return stub, nil
}

//WithBalanceOf mocks the "default" balance of all clients
func WithBalanceOf(value float64) func(shim.ChaincodeStubInterface, []string) (float64, error) {
	return func(shim.ChaincodeStubInterface, []string) (float64, error) {
		return value, nil
	}
}

//WithAllowanceOf mocks the "default" allowances of all clients
func WithAllowanceOf(value float64) func(shim.ChaincodeStubInterface, []string) (float64, error) {
	return func(shim.ChaincodeStubInterface, []string) (float64, error) {
		return value, nil
	}
}
