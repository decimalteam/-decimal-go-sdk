package api

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

////////////////////////////////////////////////////////////////
// Error - contains Decimal API error response fields.
////////////////////////////////////////////////////////////////

// Error contains Decimal API error response fields.
type Error struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Err        string `json:"error"`
}

// Error returns error info as string.
func (e *Error) Error() string {
	return fmt.Sprintf("statusCode: %d, message: \"%s\", data: \"%s\"", e.StatusCode, e.Message, e.Err)
}

////////////////////////////////////////////////////////////////
// ResponseError - wraps Resty response error.
////////////////////////////////////////////////////////////////

// ResponseError wraps Resty response error and allows to generate error info.
type ResponseError struct {
	*resty.Response
}

// NewResponseError creates new ResponseError object.
func NewResponseError(response *resty.Response) *ResponseError {
	return &ResponseError{Response: response}
}

// Error returns error info as JSON string.
func (res *ResponseError) Error() string {
	detailError := map[string]string{
		"statusCode": fmt.Sprintf("%d", res.StatusCode()),
		"status":     res.Status(),
		"time":       fmt.Sprintf("%f seconds", res.Time().Seconds()),
		"receivedAt": fmt.Sprintf("%v", res.ReceivedAt()),
		"headers":    fmt.Sprintf("%v", res.Header()),
		"body":       res.String(),
	}
	marshal, _ := json.Marshal(detailError)
	return string(marshal)
}

////////////////////////////////////////////////////////////////
// TxError - contains Decimal Node error response fields.
////////////////////////////////////////////////////////////////

// TxError contains Decimal Node error response fields.
type TxError struct {
	Height string `json:"height"`
	TxHash string `json:"txhash"`
	Code   int    `json:"code"`
	RawLog string `json:"raw_log"`
}

// Error returns error info as JSON string.
func (e *TxError) Error() string {
	return fmt.Sprintf("height: %s, txHash: %s, code: %d, raw_log: \"%s\"", e.Height, e.TxHash, e.Code, e.RawLog)
}

////////////////////////////////////////////////////////////////
// JsonRPCError - contains Decimal Node error response fields.
////////////////////////////////////////////////////////////////

// JsonRPCError contains API response.
type JsonRPCError struct {
	JSONRPC       string               `json:"jsonrpc"`
	ID            int64                `json:"id"`
	InternalError JsonRPCInternalError `json:"error"`
}

type JsonRPCInternalError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// Error returns error info as string.
func (e *JsonRPCError) Error() string {
	return fmt.Sprintf("statusCode: %d, message: \"%s\", data: \"%s\"", e.InternalError.Code,
		e.InternalError.Message, e.InternalError.Data)
}
