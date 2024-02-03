package ethrpc

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var zeroHash = common.Hash{}

type Adapter struct {
	client ETHClient

	options AdapterOptions

	requestMiddlewares  []RequestMiddleware
	responseMiddlewares []ResponseMiddleware
}

type AdapterOptions struct {
	MultiCallContractAddress common.Address
	MultiCallABI             abi.ABI
}

func NewAdapter(options ...func(*Adapter)) *Adapter {
	adapter := &Adapter{}

	for _, o := range options {
		o(adapter)
	}

	return adapter
}

func WithClient(client ETHClient) func(*Adapter) {
	return func(adapter *Adapter) {
		adapter.client = client
	}
}

func WithMulticall(multiCallContractAddress common.Address, multiCallABI abi.ABI) func(*Adapter) {
	return func(adapter *Adapter) {
		adapter.options.MultiCallContractAddress = multiCallContractAddress
		adapter.options.MultiCallABI = multiCallABI
	}
}

func WithRequestMiddlewares(middlewares ...RequestMiddleware) func(*Adapter) {
	return func(adapter *Adapter) {
		adapter.requestMiddlewares = middlewares
	}
}

func WithResponseMiddlewares(middlewares ...ResponseMiddleware) func(*Adapter) {
	return func(adapter *Adapter) {
		adapter.responseMiddlewares = middlewares
	}
}

func (a *Adapter) NewRequest() *Request {
	return &Request{
		executor: a,
	}
}

func (a *Adapter) Execute(req *Request) (*Response, error) {
	for _, f := range a.requestMiddlewares {
		if err := f(a, req); err != nil {
			return nil, err
		}
	}

	rawResponse, err := a.callContract(req)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		Request:     req,
		RawResponse: rawResponse,
	}

	for _, f := range a.responseMiddlewares {
		if err := f(a, resp); err != nil {
			return nil, err
		}
	}

	return resp, err
}

func (a *Adapter) GetMulticallContractAddress() common.Address {
	return a.options.MultiCallContractAddress
}

func (a *Adapter) GetMulticallABI() abi.ABI {
	return a.options.MultiCallABI
}

func (a *Adapter) callContract(req *Request) ([]byte, error) {
	if req.BlockHash != zeroHash {
		return a.client.CallContractAtHash(req.Context(), req.RawCallMsg, req.BlockHash)
	}

	return a.client.CallContract(req.Context(), req.RawCallMsg, req.BlockNumber)
}
