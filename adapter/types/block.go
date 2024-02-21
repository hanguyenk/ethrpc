package types

import "math/big"

type Block struct {
	Number      *big.Int `json:"number"`
	Hash        string   `json:"hash"`
	Timestamp   uint64   `json:"timestamp"`
	ParentHash  string   `json:"parentHash"`
	ReorgedHash string   `json:"reorgedHash"`
	Logs        []Log    `json:"logs"`
}
