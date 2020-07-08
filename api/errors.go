package api

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

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

// ResponseError wraps Resty respons and allows to generate error info.
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
