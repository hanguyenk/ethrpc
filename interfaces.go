package ethrpc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type RequestExecutor interface {
	Execute(request *Request) (*Response, error)
	GetMulticallContractAddress() common.Address
	GetMulticallABI() abi.ABI
}

type ETHClient interface {
	CallContract(ctx context.Context, msg CallMsg, blockNumber *big.Int) ([]byte, error)
	CallContractAtHash(ctx context.Context, msg CallMsg, blockHash common.Hash) ([]byte, error)
}
