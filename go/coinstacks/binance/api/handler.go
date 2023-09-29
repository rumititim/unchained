package api

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	ws "github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/shapeshift/unchained/coinstacks/binance"
	"github.com/shapeshift/unchained/pkg/api"
	"github.com/shapeshift/unchained/pkg/cosmos"
	"github.com/shapeshift/unchained/pkg/websocket"
)

//// custom amino unmarshaler to decode and convert to correct types
//unmarshal := func(bz []byte, result *coretypes.ResultEvent) error {
//	err := conf.Encoding.Amino.UnmarshalJSON(bz, result)
//	if err != nil {
//		return err
//	}

//	switch v := result.Data.(type) {
//	case commontypes.EventDataNewBlockHeader:
//		result.Data = &cosmos.BlockResponse{
//			Height:    int(v.Header.Height),
//			Hash:      v.Header.Hash().String(),
//			Timestamp: int(v.Header.Time.Unix()),
//		}
//	}

//	return nil
//}

type Handler struct {
	*cosmos.Handler
	HTTPClient *binance.HTTPClient
}

func (h *Handler) NewWebsocketConnection(conn *ws.Conn, manager *websocket.Manager) {
	c := websocket.NewConnection(conn, h.WSClient, manager)
	c.Start()
}

func (h *Handler) StartWebsocket() error {
	//h.WSClient.TxHandler(func(tx types.EventDataTx, block *cosmos.BlockResponse) (interface{}, []string, error) {
	//	pTx, err := rpc.ParseTx(h.HTTPClient.GetEncoding().Amino.Amino, tx.Tx)
	//	if err != nil {
	//		return nil, nil, errors.Wrapf(err, "failed to handle tx: %v, in block: %v", tx.Tx, block.Height)
	//	}

	//	txid := fmt.Sprintf("%X", sha256.Sum256(tx.Tx))

	//	t := Tx{
	//		BaseTx: api.BaseTx{
	//			TxID:        txid,
	//			BlockHash:   &block.Hash,
	//			BlockHeight: block.Height,
	//			Timestamp:   block.Timestamp,
	//		},
	//		Confirmations: 1,
	//		GasUsed:       strconv.Itoa(int(tx.Result.GasUsed)),
	//		GasWanted:     strconv.Itoa(int(tx.Result.GasWanted)),
	//		Index:         int(tx.Index),
	//		Memo:          pTx.(txtypes.StdTx).Memo,
	//		Messages:      h.ParseMessages(pTx.GetMsgs(), nil),
	//	}

	//	seen := make(map[string]bool)
	//	addrs := []string{}

	//	for _, m := range t.Messages {
	//		if m.Addresses == nil {
	//			continue
	//		}

	//		// unique set of addresses
	//		for _, addr := range m.Addresses {
	//			if _, ok := seen[addr]; !ok {
	//				addrs = append(addrs, addr)
	//				seen[addr] = true
	//			}
	//		}
	//	}

	//	return t, addrs, nil
	//})

	//err := h.WSClient.Start()
	//if err != nil {
	//	return errors.WithStack(err)
	//}

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

	balances := map[string]string{"BNB": "0"}
	for _, b := range res.Balances {
		balances[b.Symbol] = b.Free
	}

	account := cosmos.Account{
		BaseAccount: api.BaseAccount{
			Balance:            balances["BNB"],
			UnconfirmedBalance: "0",
			Pubkey:             res.Address,
		},
		AccountNumber: int(res.Number),
		Sequence:      int(res.Sequence),
		Assets:        []cosmos.Value{},
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

	fmt.Printf("%+v\n", tx)

	//t, err := h.formatTx(tx)
	//if err != nil {
	//	return nil, errors.Wrapf(err, "failed to format transaction: %s", tx.Hash.String())
	//}

	return nil, nil
}

func (h *Handler) SendTx(hex string) (string, error) {
	return h.HTTPClient.BroadcastTx(hex)
}

func (h Handler) EstimateGas(rawTx string) (string, error) {
	// no gas required for binance chain transactions
	return "0", nil
}

func (h *Handler) ParseMessages(msgs []sdk.Msg, events cosmos.EventsByMsgIndex) []cosmos.Message {
	return cosmos.ParseMessages(msgs, events)
}

func (h *Handler) ParseFee(tx signing.Tx, txid string, denom string) cosmos.Value {
	return cosmos.Fee(tx, txid, denom)
}

//func (h *Handler) formatTx(tx *coretypes.ResultTx) (*Tx, error) {
//	pTx, err := rpc.ParseTx(h.HTTPClient.GetEncoding().Amino.Amino, tx.Tx)
//	if err != nil {
//		return nil, errors.Wrapf(err, "failed to parse tx: %v", tx.Tx)
//	}
//
//	block, err := h.BlockService.GetBlock(int(tx.Height))
//	if err != nil {
//		return nil, errors.Wrapf(err, "failed to get block: %d", tx.Height)
//	}
//
//	t := &Tx{
//		BaseTx: api.BaseTx{
//			TxID:        tx.Hash.String(),
//			BlockHash:   &block.Hash,
//			BlockHeight: block.Height,
//			Timestamp:   block.Timestamp,
//		},
//		// TODO: reference fees from /api/v1/fees
//		Fee:           cosmos.Value{Amount: "0", Denom: "BNB"},
//		Confirmations: h.BlockService.Latest.Height - int(tx.Height) + 1,
//		GasUsed:       strconv.Itoa(int(tx.TxResult.GasUsed)),
//		GasWanted:     strconv.Itoa(int(tx.TxResult.GasWanted)),
//		Index:         int(tx.Index),
//		Memo:          pTx.(txtypes.StdTx).Memo,
//		Messages:      binance.ParseMessages(pTx.GetMsgs()),
//	}
//
//	return t, nil
//}
