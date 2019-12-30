package helpers

import (
	"encoding/json"
	"strconv"
)

/*DefaultToZeroIfEmpty return empty array of bytes (0) if the input is empty*/
func DefaultToZeroIfEmpty(value []byte) []byte {
	if len(value) == 0 {
		return []byte{48} //represents "0"
	}
	return value
}

/*BufferToFloat convert array of bytes to float64*/
func BufferToFloat(buffer []byte) float64 {
	value, err := strconv.ParseFloat(string(buffer), 64)
	if err != nil {
		panic(err)
	}
	return value
}

/*StringToInt converts type string to type int*/
func StringToInt(str string) int {
	n, err := strconv.Atoi(str)
	if err != nil {
		panic(err.Error())
	}
	return n
}

/*StringToFloat converts type string to type float64*/
func StringToFloat(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(err.Error())
	}
	return f
}

/*FloatToString converts a float64 number to string*/
func FloatToString(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

/*FloatToBuffer converts a float64 number to array of bytes*/
func FloatToBuffer(value float64) []byte {
	return []byte(FloatToString(value))
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
