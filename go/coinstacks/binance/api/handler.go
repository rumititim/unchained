package api

import (
	"crypto/sha256"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	ws "github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/shapeshift/bnb-chain-go-sdk/client/rpc"
	txtypes "github.com/shapeshift/bnb-chain-go-sdk/types/tx"
	"github.com/shapeshift/unchained/coinstacks/binance"
	"github.com/shapeshift/unchained/pkg/api"
	"github.com/shapeshift/unchained/pkg/cosmos"
	"github.com/shapeshift/unchained/pkg/websocket"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

type Handler struct {
	*cosmos.Handler
	HTTPClient *binance.HTTPClient
	WSClient   *binance.WSClient
}

func (h *Handler) NewWebsocketConnection(conn *ws.Conn, manager *websocket.Manager) {
	c := websocket.NewConnection(conn, h.WSClient, manager)
	c.Start()
}

func (h *Handler) StartWebsocket() error {
	h.WSClient.TxHandler(func(tx types.EventDataTx, block *cosmos.BlockResponse) (interface{}, []string, error) {
		pTx, err := rpc.ParseTx(h.HTTPClient.GetEncoding().Amino.Amino, tx.Tx)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to handle tx: %v, in block: %v", tx.Tx, block.Height)
		}

		txid := fmt.Sprintf("%X", sha256.Sum256(tx.Tx))

		fees, err := h.HTTPClient.GetFees(&block.Height)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to get fees")
		}

		fee := int64(0)
		for _, msg := range pTx.GetMsgs() {
			fee += fees[msg.Type()]
		}

		t := cosmos.Tx{
			BaseTx: api.BaseTx{
				TxID:        txid,
				BlockHash:   &block.Hash,
				BlockHeight: block.Height,
				Timestamp:   block.Timestamp,
			},
			Fee:           cosmos.Value{Amount: strconv.Itoa(int(fee)), Denom: "BNB"},
			Confirmations: 1,
			GasUsed:       strconv.Itoa(int(tx.Result.GasUsed)),
			GasWanted:     strconv.Itoa(int(tx.Result.GasWanted)),
			Index:         int(tx.Index),
			Memo:          pTx.(txtypes.StdTx).Memo,
			Messages:      h.ParseMessages(pTx.GetMsgs(), nil),
		}

		seen := make(map[string]bool)
		addrs := []string{}

		for _, m := range t.Messages {
			if m.Addresses == nil {
				continue
			}

			// unique set of addresses
			for _, addr := range m.Addresses {
				if _, ok := seen[addr]; !ok {
					addrs = append(addrs, addr)
					seen[addr] = true
				}
			}
		}

		return t, addrs, nil
	})

	err := h.WSClient.Start()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (h *Handler) StopWebsocket() {
	h.WSClient.Stop()
}

// Contains info about the running coinstack
// swagger:model Info
type Info struct {
	// swagger:allOf
	cosmos.Info
}

func (h *Handler) GetInfo() (api.Info, error) {
	info, err := h.Handler.GetInfo()
	if err != nil {
		return nil, err
	}

	i := Info{Info: info.(cosmos.Info)}

	return i, nil
}

// Contains info about account details for an address
// swagger:model Account
type Account struct {
	// swagger:allOf
	cosmos.Account
}

func (h *Handler) GetAccount(pubkey string) (api.Account, error) {
	var res struct {
		Number   int64  `json:"account_number"`
		Address  string `json:"address"`
		Balances []struct {
			Symbol string `json:"symbol"`
			Free   string `json:"free"`
			Locked string `json:"locked"`
			Frozen string `json:"frozen"`
		} `json:"balances"`
		PublicKey []uint8 `json:"public_key"`
		Sequence  int64   `json:"sequence"`
		Flags     uint64  `json:"flags"`
	}

	_, err := h.HTTPClient.LCD.R().SetResult(&res).Get(fmt.Sprintf("/api/v1/account/%s", pubkey))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	balance := "0"
	assets := []cosmos.Value{}
	for _, b := range res.Balances {
		if b.Symbol == "BNB" {
			balance = b.Free
		} else {
			assets = append(assets, cosmos.Value{Amount: b.Free, Denom: b.Symbol})
		}
	}

	account := cosmos.Account{
		BaseAccount: api.BaseAccount{
			Balance:            balance,
			UnconfirmedBalance: "0",
			Pubkey:             res.Address,
		},
		AccountNumber: int(res.Number),
		Sequence:      int(res.Sequence),
		Assets:        assets,
	}

	return account, nil
}

//func (h *Handler) GetTxHistory(pubkey string, cursor string, pageSize int) (api.TxHistory, error) {
//	res, err := h.HTTPClient.GetTxHistory(pubkey, cursor, pageSize)
//	if err != nil {
//		return nil, errors.Wrapf(err, "failed to get tx history")
//	}
//
//	txs := []Tx{}
//	for _, t := range res.Txs {
//		tx, err := h.formatTx(t)
//		if err != nil {
//			return nil, errors.Wrapf(err, "failed to format transaction: %s", t.Hash)
//		}
//
//		txs = append(txs, *tx)
//	}
//
//	txHistory := TxHistory{
//		BaseTxHistory: api.BaseTxHistory{
//			Pagination: api.Pagination{
//				Cursor: res.Cursor,
//			},
//			Pubkey: pubkey,
//		},
//		Txs: txs,
//	}
//
//	return txHistory, nil
//}

func (h *Handler) GetTx(txid string) (api.Tx, error) {
	tx, err := h.HTTPClient.GetTx(txid)
	if err != nil {
		return nil, err
	}

	t, err := h.formatTx(tx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to format transaction: %s", tx.Hash.String())
	}

	return t, nil
}

func (h *Handler) SendTx(hex string) (string, error) {
	return h.HTTPClient.BroadcastTx(hex)
}

func (h Handler) EstimateGas(rawTx string) (string, error) {
	// no gas required for binance chain transactions
	return "0", nil
}

func (h Handler) GetFees() (binance.Fees, error) {
	return h.HTTPClient.GetFees(nil)
}

func (h *Handler) ParseMessages(msgs interface{}, events cosmos.EventsByMsgIndex) []cosmos.Message {
	return binance.ParseMessages(msgs)
}

func (h *Handler) ParseFee(tx signing.Tx, txid string, denom string) cosmos.Value {
	return cosmos.Fee(tx, txid, denom)
}

func (h *Handler) formatTx(tx *coretypes.ResultTx) (*cosmos.Tx, error) {
	pTx, err := rpc.ParseTx(h.HTTPClient.GetEncoding().Amino.Amino, tx.Tx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse tx: %v", tx.Tx)
	}

	height := int(tx.Height)

	block, err := h.BlockService.GetBlock(height)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get block: %d", height)
	}

	fees, err := h.HTTPClient.GetFees(&height)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get fees")
	}

	fee := int64(0)
	for _, msg := range pTx.GetMsgs() {
		fee += fees[msg.Type()]
	}

	t := &cosmos.Tx{
		BaseTx: api.BaseTx{
			TxID:        tx.Hash.String(),
			BlockHash:   &block.Hash,
			BlockHeight: block.Height,
			Timestamp:   block.Timestamp,
		},
		Fee:           cosmos.Value{Amount: strconv.Itoa(int(fee)), Denom: "BNB"},
		Confirmations: h.BlockService.Latest.Height - height + 1,
		GasUsed:       strconv.Itoa(int(tx.TxResult.GasUsed)),
		GasWanted:     strconv.Itoa(int(tx.TxResult.GasWanted)),
		Index:         int(tx.Index),
		Memo:          pTx.(txtypes.StdTx).Memo,
		Messages:      h.ParseMessages(pTx.GetMsgs(), nil),
	}

	return t, nil
}
