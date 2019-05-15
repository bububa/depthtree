package client

import "fmt"

type ErrorCode = int

const (
	BADREQUEST_ERROR            ErrorCode = 400
	INTERNAL_ERROR              ErrorCode = 500
	NOTFOUND_ERROR              ErrorCode = 404
	UNAUTHORIZED_ERROR          ErrorCode = 401
	FEATURE_NOT_AVAILABLE_ERROR ErrorCode = 402
)

type Error struct {
	Code ErrorCode `json:"code,omitempty"`
	Msg  string    `json:"message,omitempty"`
}

func (this Error) Error() string {
	return fmt.Sprintf("CODE:%d, MSG:%s", this.Code, this.Msg)
}
