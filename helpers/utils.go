package helpers

import (
	"encoding/json"
	"math/big"
	"strconv"
)

/*DefaultToZeroIfEmpty return empty array of bytes (0) if the input is empty*/
func DefaultToZeroIfEmpty(value []byte) []byte {
	if len(value) == 0 {
		return []byte{48} //represents "0"
	}
	return value
}

/*BufferToBigInt convert array of bytes to big.Int*/
func BufferToBigInt(buffer []byte) *big.Int {
	n := &big.Int{}
	v, ok := n.SetString(string(buffer), 10)
	if !ok {
		panic("BufferToBigInt conversion failed")
	}
	return v
}

/*StringToInt converts type string to type int64*/
func StringToInt(str string) int64 {
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err.Error())
	}
	return n
}

/*StringToBigInt converts type string to type big.Int*/
func StringToBigInt(str string) *big.Int {
	n := &big.Int{}
	value, success := n.SetString(str, 10)
	if !success {
		panic("StringToBigInt conversion failed")
	}
	return value
}

/*Add two big Int numbers*/
func Add(a *big.Int, b *big.Int) *big.Int {
	n := &big.Int{}
	return n.Add(a, b)
}

/*Sub two big Int numbers*/
func Sub(a *big.Int, b *big.Int) *big.Int {
	n := &big.Int{}
	return n.Sub(a, b)
}

/*Mul two big Int numbers*/
func Mul(a *big.Int, b *big.Int) *big.Int {
	n := &big.Int{}
	return n.Mul(a, b)
}

/*Pow returns the a raised to the power of b in big Int type*/
func Pow(a int64, b int64) *big.Int {
	n := &big.Int{}
	return n.Exp(big.NewInt(a), big.NewInt(b), nil)
}

/*MalshalJSON returns the JSON encoding of `payload`*/
func MalshalJSON(payload interface{}) []byte {
	bytes, err := json.Marshal(payload)
	if err != nil {
		panic(err.Error())
	}
	return bytes
}

/*JSONToMap returns a map from the input JSON string*/
func JSONToMap(str string) map[string]interface{} {
	var m map[string]interface{} //convert the `coinConfigBytes` into a `map`
	if err := json.Unmarshal([]byte(str), &m); err != nil {
		panic(err.Error())
	}
	return m
}
