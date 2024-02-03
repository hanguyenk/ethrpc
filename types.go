package ethrpc

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// CallMsg contains parameters for contract calls.
type CallMsg struct {
	From      common.Address  // the sender of the 'transaction'
	To        *common.Address // the destination contract (nil for contract creation)
	Gas       uint64          // if 0, the call executes with near-infinite gas
	GasPrice  *big.Int        // wei <-> gas exchange ratio
	GasFeeCap *big.Int        // EIP-1559 fee cap per gas.
	GasTipCap *big.Int        // EIP-1559 tip per gas.
	Value     *big.Int        // amount of wei sent along with the call
	Data      []byte          // input data, usually an ABI-encoded contract method invocation

	AccessList AccessList // EIP-2930 access list.
}

// AccessList is an EIP-2930 access list.
type AccessList []AccessTuple

// AccessTuple is the element type of an access list.
type AccessTuple struct {
	Address     common.Address
	StorageKeys []common.Hash
}

type MultiCallParam struct {
	Target   common.Address
	CallData []byte
}

type AggregateResult struct {
	BlockNumber *big.Int
	ReturnData  [][]byte
}

type TryAggregateResult struct {
	Success    bool
	ReturnData []byte
}

type TryAggregateResultList []TryAggregateResult

type TryBlockAndAggregateResult struct {
	BlockNumber *big.Int
	BlockHash   [32]byte
	ReturnData  []TryAggregateResult
}
