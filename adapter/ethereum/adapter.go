package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	adaptertypes "github.com/hanguyenk/ethrpc/adapter/types"
)

type Adapter struct {
	client *ethclient.Client
}

func NewAdapter(url string) (*Adapter, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		client: client,
	}, nil
}

func (a *Adapter) CallContract(ctx context.Context, msg adaptertypes.CallMsg, blockNumber *big.Int) ([]byte, error) {
	ethereumCallMsg := a.convertToEthereumCallMsg(msg)

	return a.client.CallContract(ctx, ethereumCallMsg, blockNumber)
}

func (a *Adapter) CallContractAtHash(ctx context.Context, msg adaptertypes.CallMsg, blockHash common.Hash) ([]byte, error) {
	ethereumCallMsg := a.convertToEthereumCallMsg(msg)

	return a.client.CallContractAtHash(ctx, ethereumCallMsg, blockHash)
}

func (a *Adapter) convertToEthereumCallMsg(originMsg adaptertypes.CallMsg) ethereum.CallMsg {
	return ethereum.CallMsg{
		From:       originMsg.From,
		To:         originMsg.To,
		Gas:        originMsg.Gas,
		GasPrice:   originMsg.GasPrice,
		GasFeeCap:  originMsg.GasFeeCap,
		GasTipCap:  originMsg.GasTipCap,
		Value:      originMsg.Value,
		Data:       originMsg.Data,
		AccessList: a.convertToEthereumAccessList(originMsg.AccessList),
	}
}

func (a *Adapter) convertToEthereumAccessList(originAccessList adaptertypes.AccessList) types.AccessList {
	accessList := make([]types.AccessTuple, 0, len(originAccessList))

	for _, originAccessTuple := range originAccessList {
		accessList = append(accessList, types.AccessTuple{
			Address:     originAccessTuple.Address,
			StorageKeys: originAccessTuple.StorageKeys,
		})
	}

	return accessList
}
