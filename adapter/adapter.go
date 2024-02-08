package adapter

import (
	"github.com/hanguyenk/ethrpc/adapter/ethereum"

	"github.com/hanguyenk/ethrpc"
)

func New(chainID uint, url string) (ethrpc.EthClientAdapter, error) {
	switch chainID {
	default:
		return ethereum.NewAdapter(url)
	}
}
