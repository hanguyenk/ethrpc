package adapter

import (
	"github.com/hanguyenk/ethrpc/adapter/avalanche"
	"github.com/hanguyenk/ethrpc/adapter/ethereum"
)

var (
	ChainIDAvalancheCChain uint = 43114
)

func New(chainID uint, url string) (EthClientAdapter, error) {
	switch chainID {
	case ChainIDAvalancheCChain:
		return avalanche.NewAdapter(url)
	default:
		return ethereum.NewAdapter(url)
	}
}
