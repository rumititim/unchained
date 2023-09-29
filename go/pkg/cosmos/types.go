package cosmos

type ErrorResponse struct {
	Code   int           `json:"code"`
	Msg    string        `json:"message"`
	Detail []interface{} `json:"detail"`
}

type Pagination struct {
	NextKey *[]byte `json:"next_key,omitempty"`
	Total   uint64  `json:"total,string,omitempty"`
}
