package adapter

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/hanguyenk/ethrpc/adapter/types"
)

// EthClientAdapter ...
type EthClientAdapter interface {
	CallContract(ctx context.Context, msg types.CallMsg, blockNumber *big.Int) ([]byte, error)
	CallContractAtHash(ctx context.Context, msg types.CallMsg, blockHash common.Hash) ([]byte, error)
}
