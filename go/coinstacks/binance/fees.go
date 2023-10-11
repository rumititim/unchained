package binance

import (
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
	commontypes "github.com/shapeshift/bnb-chain-go-sdk/common/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

// Contains info about fees by message type
// swagger:model Fees
type Fees map[string]int64

func (c *HTTPClient) GetFees(height *int) (Fees, error) {
	res := &rpctypes.RPCResponse{}

	hs := "0"
	if height != nil {
		hs = strconv.Itoa(*height)
	}

	params := map[string]string{
		"path":   "\"/param/fees\"",
		"height": hs,
	}

	_, err := c.RPC.R().SetResult(res).SetError(res).SetQueryParams(params).Get("/abci_query")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get fees")
	}
	if res.Error != nil {
		return nil, errors.Wrap(errors.New(res.Error.Error()), "failed to get fees")
	}

	var result coretypes.ResultABCIQuery
	err = json.Unmarshal(res.Result, &result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal result")
	}

	var feeParams []commontypes.FeeParam
	err = c.GetEncoding().Amino.Amino.UnmarshalBinaryLengthPrefixed(result.Response.Value, &feeParams)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal fees")
	}

	fees := make(Fees)
	for _, fee := range feeParams {
		switch v := fee.(type) {
		case *commontypes.FixedFeeParams:
			fees[v.MsgType] = v.Fee
		case *commontypes.TransferFeeParam:
			fees[v.MsgType] = v.Fee
		}
	}

	return fees, nil
}
