package binance

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
	commontypes "github.com/shapeshift/bnb-chain-go-sdk/common/types"
	msgtypes "github.com/shapeshift/bnb-chain-go-sdk/types/msg"
	"github.com/shapeshift/unchained/pkg/cosmos"
)

func (c *HTTPClient) GetTxHistory(address string, cursor string, pageSize int) (*cosmos.TxHistoryResponse, error) {
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -7)

	history := &History{
		ctx:      c.ctx,
		cursor:   &Cursor{StartTime: startTime.UnixMilli(), EndTime: endTime.UnixMilli()},
		pageSize: pageSize,
		bc:       c.bc,
		encoding: c.HTTPClient.GetEncoding(),
	}

	if cursor != "" {
		if err := history.cursor.decode(cursor); err != nil {
			return nil, errors.Wrapf(err, "failed to decode cursor: %s", cursor)
		}
	}

	history.TxState = &TxState{
		hasMore:   true,
		startTime: history.cursor.StartTime,
		endTime:   history.cursor.EndTime,
	}

	txHistory, err := history.fetch()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get tx history for address: %s", address)
	}

	return txHistory, nil
}

func ParseMessages(msgs interface{}) []cosmos.Message {
	messages := []cosmos.Message{}

	coinToValue := func(c *commontypes.Coin) cosmos.Value {
		if c == nil {
			return cosmos.Value{}
		}

		return cosmos.Value{
			Amount: strconv.Itoa(int(c.Amount)),
			Denom:  c.Denom,
		}
	}

	for i, msg := range msgs.([]msgtypes.Msg) {
		switch v := msg.(type) {
		case msgtypes.SendMsg:
			addresses := []string{}
			for _, a := range v.GetInvolvedAddresses() {
				addresses = append(addresses, a.String())
			}

			message := cosmos.Message{
				Addresses: addresses,
				Index:     strconv.Itoa(i),
				Origin:    v.GetSigners()[0].String(),
				From:      v.Inputs[0].Address.String(),
				To:        v.Outputs[0].Address.String(),
				Value:     coinToValue(&v.Inputs[0].Coins[0]),
				Type:      v.Type(),
			}

			messages = append(messages, message)
		}
	}

	return messages
}
