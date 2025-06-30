package internal

type ErrRestResponse struct {
	Code ErrorCode `json:"code"`
	Msg  string    `json:"msg"`
}
