package binance

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/shapeshift/unchained/pkg/cosmos"
	"github.com/tendermint/tendermint/libs/bytes"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

type ResultBlock struct {
	BlockMeta struct {
		BlockId struct {
			Hash bytes.HexBytes `json:"hash"`
		} `json:"block_id"`
		Header struct {
			Height int64     `json:"height"`
			Time   time.Time `json:"time"`
		} `json:"header"`
	} `json:"block_meta"`
}

func (r *ResultBlock) GetBlockResponse() *cosmos.BlockResponse {
	return &cosmos.BlockResponse{
		Height:    int(r.BlockMeta.Header.Height),
		Hash:      r.BlockMeta.BlockId.Hash.String(),
		Timestamp: int(r.BlockMeta.Header.Time.Unix()),
	}
}

func (c *HTTPClient) GetBlock(height *int) (cosmos.Block, error) {
	var res *rpctypes.RPCResponse

	hs := ""
	if height != nil {
		hs = strconv.Itoa(*height)
	}

	_, err := c.RPC.R().SetResult(&res).SetQueryParam("height", hs).Get("/block")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get block: %s", hs)
	}

	if res.Error != nil {
		return nil, errors.Errorf("failed to get block: %s: %s", hs, res.Error.Error())
	}

	result := &ResultBlock{}
	if err := c.GetEncoding().Amino.UnmarshalJSON(res.Result, result); err != nil {
		return nil, errors.Wrapf(err, "failed to decode block: %v", res.Result)
	}

	return result, nil
}
