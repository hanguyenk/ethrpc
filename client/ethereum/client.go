package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/hanguyenk/ethrpc"
)

type Client struct {
	client *ethclient.Client
}

func NewClient(url string) (*Client, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) CallContract(ctx context.Context, msg ethrpc.CallMsg, blockNumber *big.Int) ([]byte, error) {
	ethereumCallMsg := c.convertToEthereumCallMsg(msg)

	return c.client.CallContract(ctx, ethereumCallMsg, blockNumber)
}

func (c *Client) CallContractAtHash(ctx context.Context, msg ethrpc.CallMsg, blockHash common.Hash) ([]byte, error) {
	ethereumCallMsg := c.convertToEthereumCallMsg(msg)

	return c.client.CallContractAtHash(ctx, ethereumCallMsg, blockHash)
}

func (c *Client) convertToEthereumCallMsg(originMsg ethrpc.CallMsg) ethereum.CallMsg {
	return ethereum.CallMsg{
		From:       originMsg.From,
		To:         originMsg.To,
		Gas:        originMsg.Gas,
		GasPrice:   originMsg.GasPrice,
		GasFeeCap:  originMsg.GasFeeCap,
		GasTipCap:  originMsg.GasTipCap,
		Value:      originMsg.Value,
		Data:       originMsg.Data,
		AccessList: c.convertToEthereumAccessList(originMsg.AccessList),
	}
}

func (c *Client) convertToEthereumAccessList(originAccessList ethrpc.AccessList) types.AccessList {
	accessList := make([]types.AccessTuple, 0, len(originAccessList))

	for _, originAccessTuple := range originAccessList {
		accessList = append(accessList, types.AccessTuple{
			Address:     originAccessTuple.Address,
			StorageKeys: originAccessTuple.StorageKeys,
		})
	}

	return accessList
}
