package ethrpc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/hanguyenk/ethrpc/adapter/types"
)

var (
	ErrRequestParserIsNotRegistered = errors.New("requestParser is not registered")
	ErrWrongCallParams              = errors.New("wrong call params")
)

var (
	requestParserRegistry = map[RequestMethod]RequestParser{}
)

func RegisterRequestParser(method RequestMethod, parser RequestParser) {
	requestParserRegistry[method] = parser
}

func GetRequestParser(method RequestMethod) (RequestParser, error) {
	parser, ok := requestParserRegistry[method]
	if !ok {
		return nil, errors.Wrapf(ErrRequestParserIsNotRegistered, "method: [%s]", method)
	}

	return parser, nil
}

func init() {
	RegisterRequestParser(RequestMethodCall, RequestParserCall)
	RegisterRequestParser(RequestMethodAggregate, RequestParserAggregate)
	RegisterRequestParser(RequestMethodTryBlockAndAggregate, RequestParserTryBlockAndAggregate)
}

type RequestParser func(executor RequestExecutor, req *Request) error

// RequestParserCall parses raw call msg for RequestMethodCall
func RequestParserCall(_ RequestExecutor, req *Request) error {
	if len(req.Calls) != 1 {
		return ErrWrongCallParams
	}

	call := req.Calls[0]
	data, err := call.ABI.Pack(call.Method, call.Params...)
	if err != nil {
		return err
	}

	target := common.HexToAddress(call.Target)

	req.RawCallMsg = types.CallMsg{
		To:   &target,
		Data: data,
	}

	return nil

}

// RequestParserAggregate parse raw call msg for RequestMethodAggregate
func RequestParserAggregate(executor RequestExecutor, req *Request) error {
	multicallContractAddress := executor.GetMulticallContractAddress()
	multicallABI := executor.GetMulticallABI()

	// TODO: validate multicall contract address

	var multiCallParamList []MultiCallParam

	for _, c := range req.Calls {
		callData, err := c.ABI.Pack(c.Method, c.Params...)
		if err != nil {
			return err
		}

		multiCallParamList = append(
			multiCallParamList,
			MultiCallParam{
				Target:   common.HexToAddress(c.Target),
				CallData: callData,
			},
		)
	}

	data, err := multicallABI.Pack(string(RequestMethodAggregate), multiCallParamList)
	if err != nil {
		return err
	}

	msg := types.CallMsg{To: &multicallContractAddress, Data: data}
	req.RawCallMsg = msg

	return nil
}

// RequestParserTryBlockAndAggregate parse raw call msg for RequestMethodTryBlockAndAggregate
func RequestParserTryBlockAndAggregate(executor RequestExecutor, req *Request) error {
	multicallContractAddress := executor.GetMulticallContractAddress()
	multicallABI := executor.GetMulticallABI()

	// TODO: validate multicall contract address

	var multiCallParamList []MultiCallParam

	for _, call := range req.Calls {
		callData, err := call.ABI.Pack(call.Method, call.Params...)
		if err != nil {
			return err
		}

		multiCallParamList = append(
			multiCallParamList, MultiCallParam{
				Target:   common.HexToAddress(call.Target),
				CallData: callData,
			},
		)
	}

	data, err := multicallABI.Pack(string(RequestMethodTryBlockAndAggregate), req.RequireSuccess, multiCallParamList)
	if err != nil {
		return err
	}

	msg := types.CallMsg{To: &multicallContractAddress, Data: data}
	req.RawCallMsg = msg

	return nil
}
