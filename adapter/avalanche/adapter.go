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

func (a *Adapter) SubscribeNewHead(ctx context.Context, headerChannel chan<- *adaptertypes.Header) (adaptertypes.Subscription, error) {
	originHeaderChannel := make(chan *types.Header)
	sub, err := a.client.SubscribeNewHead(ctx, originHeaderChannel)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(originHeaderChannel)

		for {
			select {
			case <-ctx.Done():
				return
			case originHeader := <-originHeaderChannel:
				headerChannel <- a.convertFromAvalancheHeader(originHeader)
			}
		}
	}()

	return sub, nil
}

func (a *Adapter) FilterLogs(ctx context.Context, query adaptertypes.FilterQuery) ([]adaptertypes.Log, error) {
	logs, err := a.client.FilterLogs(ctx, a.convertToAvalancheFilterQuery(query))
	if err != nil {
		return nil, err
	}

	return a.convertFromAvalancheLogs(logs), nil
}

func (a *Adapter) BlockNumber(ctx context.Context) (uint64, error) {
	return a.client.BlockNumber(ctx)
}

func (a *Adapter) HeaderByHash(ctx context.Context, hash common.Hash) (*adaptertypes.Header, error) {
	originHeader, err := a.client.HeaderByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	return a.convertFromAvalancheHeader(originHeader), nil
}

func (a *Adapter) HeaderByNumber(ctx context.Context, number *big.Int) (*adaptertypes.Header, error) {
	originHeader, err := a.client.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	return a.convertFromAvalancheHeader(originHeader), nil
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

func (a *Adapter) convertToAvalancheFilterQuery(originFilterQuery adaptertypes.FilterQuery) interfaces.FilterQuery {
	return interfaces.FilterQuery{
		BlockHash: originFilterQuery.BlockHash,
		FromBlock: originFilterQuery.FromBlock,
		ToBlock:   originFilterQuery.ToBlock,
		Addresses: originFilterQuery.Addresses,
		Topics:    originFilterQuery.Topics,
	}
}

func (a *Adapter) convertFromAvalancheHeader(originHeader *types.Header) *adaptertypes.Header {
	return &adaptertypes.Header{
		Hash:       originHeader.Hash(),
		ParentHash: originHeader.ParentHash,
		Number:     originHeader.Number,
		Time:       originHeader.Time,
	}
}

func (a *Adapter) convertFromAvalancheLogs(originLogs []types.Log) []adaptertypes.Log {
	logs := make([]adaptertypes.Log, 0, len(originLogs))
	for _, origin := range originLogs {
		topics := make([]common.Hash, 0, len(origin.Topics))

		for _, topic := range origin.Topics {
			topics = append(topics, topic)
		}

		logs = append(logs, adaptertypes.Log{
			Address:     origin.Address,
			Topics:      topics,
			Data:        origin.Data,
			BlockNumber: origin.BlockNumber,
			TxHash:      origin.TxHash,
			TxIndex:     origin.TxIndex,
			BlockHash:   origin.BlockHash,
			Index:       origin.Index,
			Removed:     origin.Removed,
		})
	}

	return logs
}
