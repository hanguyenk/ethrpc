package client

import (
	"github.com/hanguyenk/ethrpc"
	"github.com/hanguyenk/ethrpc/client/ethereum"
)

func NewClient(chainID uint, url string) (ethrpc.ETHClient, error) {
	switch chainID {
	default:
		return ethereum.NewClient(url)
	}
}
