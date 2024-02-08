package avalanche

import (
	"context"
	"math/big"

	"github.com/ava-labs/coreth/core/types"
	"github.com/ava-labs/coreth/ethclient"
	"github.com/ava-labs/coreth/interfaces"
	"github.com/ethereum/go-ethereum/common"

	adaptertypes "github.com/hanguyenk/ethrpc/adapter/types"
)

type Adapter struct {
	client ethclient.Client
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
	avalancheCallMsg := a.convertToAvalancheCallMsg(msg)

	return a.client.CallContract(ctx, avalancheCallMsg, blockNumber)
}

func (a *Adapter) CallContractAtHash(ctx context.Context, msg adaptertypes.CallMsg, blockHash common.Hash) ([]byte, error) {
	avalancheCallMsg := a.convertToAvalancheCallMsg(msg)

	return a.client.CallContractAtHash(ctx, avalancheCallMsg, blockHash)
}

func (a *Adapter) convertToAvalancheCallMsg(originMsg adaptertypes.CallMsg) interfaces.CallMsg {
	return interfaces.CallMsg{
		From:       originMsg.From,
		To:         originMsg.To,
		Gas:        originMsg.Gas,
		GasPrice:   originMsg.GasPrice,
		GasFeeCap:  originMsg.GasFeeCap,
		GasTipCap:  originMsg.GasTipCap,
		Value:      originMsg.Value,
		Data:       originMsg.Data,
		AccessList: a.convertToAvalancheAccessList(originMsg.AccessList),
	}
}

func (a *Adapter) convertToAvalancheAccessList(originAccessList adaptertypes.AccessList) types.AccessList {
	accessList := make([]types.AccessTuple, 0, len(originAccessList))

	for _, originAccessTuple := range originAccessList {
		accessList = append(accessList, types.AccessTuple{
			Address:     originAccessTuple.Address,
			StorageKeys: originAccessTuple.StorageKeys,
		})
	}

	return accessList
}
