package erc20events

import "math/big"

/*enums for event names*/
const (
	TRANSFER = "transfer"
	APPROVAL = "approval"
)

/*Payload of the event*/
type Payload struct {
	From   string   `json:"from"`
	To     string   `json:"to"`
	Amount *big.Int `json:"amount"`
}

/*Event object to emit to clients, will be sent as JSON format*/
type Event struct {
	Origin  string  `json:"origin"` /*transaction invoker's ID*/
	Payload Payload `json:"payload"`
}
