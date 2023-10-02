package binance

import (
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
)

type AccAddress []byte
type ABCICodeType uint32
type CodeType uint16
type CodespaceType uint16

type KVPair struct {
	Key   []byte
	Value []byte
}

type KVPairs []KVPair
type Tags = KVPairs
type Event = abci.Event
type Events = []Event

func (e Events) AppendEvent(event Event) Events {
	return append(e, event)
}

// AppendEvents adds a slice of Event objects to an exist slice of Event objects.
func (e Events) AppendEvents(events Events) Events {
	return append(e, events...)
}

// ToABCIEvents converts a slice of Event objects to a slice of abci.Event
// objects.
func (e Events) ToABCIEvents() []abci.Event {
	res := make([]abci.Event, len(e))
	for i, ev := range e {
		res[i] = abci.Event{Type: ev.Type, Attributes: ev.Attributes}
	}

	return res
}

type Result struct {

	// Code is the response code, is stored back on the chain.
	Code ABCICodeType

	// Data is any data returned from the app.
	Data []byte

	// Log is just debug information. NOTE: nondeterministic.
	Log string

	// Tx fee amount and denom.
	FeeAmount int64
	FeeDenom  string

	// Tags are used for transaction indexing and pubsub.
	Tags   Tags
	Events []abci.Event
}

type Error interface {
	Error() string
	Stacktrace() Error
	Trace(offset int, format string, args ...interface{}) Error
	Data() interface{}

	// convenience
	TraceSDK(format string, args ...interface{}) Error

	// set codespace
	WithDefaultCodespace(CodespaceType) Error

	RawError() string // return raw error message
	Code() CodeType
	Codespace() CodespaceType
	ABCILog() string
	ABCICode() ABCICodeType
	Result() Result
	QueryResult() abci.ResponseQuery
}

type Msg interface {
	Route() string
	Type() string
	ValidateBasic() Error
	GetSignBytes() []byte
	GetSigners() []AccAddress
	GetInvolvedAddresses() []AccAddress
}

type Tx interface {
	GetMsgs() []Msg
}

type StdSignature struct {
	crypto.PubKey `json:"pub_key"` // optional
	Signature     []byte           `json:"signature"`
	AccountNumber int64            `json:"account_number"`
	Sequence      int64            `json:"sequence"`
}

type StdTx struct {
	Msgs       []Msg          `json:"msg"`
	Signatures []StdSignature `json:"signatures"`
	Memo       string         `json:"memo"`
	Source     int64          `json:"source"`
	Data       []byte         `json:"data"`
}

func (tx StdTx) GetMsgs() []Msg { return tx.Msgs }

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterInterface((*Tx)(nil), nil)
	cdc.RegisterInterface((*Msg)(nil), nil)
	cdc.RegisterConcrete(StdTx{}, "auth/StdTx", nil)
}

func ParseTx(cdc *amino.Codec, txBytes []byte) (Tx, error) {
	var parsedTx StdTx
	err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &parsedTx)

	if err != nil {
		return nil, err
	}

	return parsedTx, nil
}
