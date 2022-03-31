package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

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
func (e Error) Error() string {
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
func (res ResponseError) Error() string {
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
func (e TxError) Error() string {
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
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Error returns error info as string.
func (e JsonRPCError) Error() string {
	return fmt.Sprintf("statusCode: %d, message: \"%s\", data: \"%s\"", e.InternalError.Code,
		e.InternalError.Message, e.InternalError.Data)
}

// Error returns error info as string.
func (e JsonRPCInternalError) Error() string {
	return fmt.Sprintf("statusCode: %d, message: \"%s\", data: \"%s\"", e.Code,
		e.Message, e.Data)
}

// Error for queries without RPC/REST implementation
var ErrNotImplemented = errors.New("not implemented")

// Error indicating for universal decoding
var ErrIsRPCError = errors.New("rpc error")
var ErrMissing = errors.New("universal JSON decode missing logic")

////////////////////////////////////////////////////////////////
// Function to decrease boilerplate handling
////////////////////////////////////////////////////////////////

func processConnectionError(response *resty.Response, err error) error {
	if err != nil {
		return err
	}
	if response.IsError() {
		return NewResponseError(response)
	}
	return nil
}

// Universal JSON decoding
// valueStruct - MUST BE pointer to struct
// errorStruct - MUST BE pointer to struct or nil
// return error if valueStruct unmarshaling failed or valueStruct not ok
type validationFuncType func() (bool, bool)

func universalJSONDecode(source []byte, valueStruct interface{}, errorStruct interface{}, validator validationFuncType) error {
	var err1, err2 error
	err1 = json.Unmarshal(source, valueStruct)
	if errorStruct != nil {
		err2 = json.Unmarshal(source, errorStruct)
	}
	okValue, okError := validator()
	// all ok
	if okValue && err1 == nil {
		return nil
	}
	// error ok
	if okError {
		return ErrIsRPCError // indicate that error in errorStruct
	}

	// error during Unmarshaling (wrong JSON)
	if !okError && err1 != nil {
		return err1
	}

	if err2 != nil {
		return err2
	}
	log.Printf("INTERNAL: missing logic.")
	if valueStruct != nil {
		log.Printf("Value: %#v", valueStruct)
	}
	if errorStruct != nil {
		log.Printf("Error: %#v", errorStruct)
	}
	log.Printf("okValue=%v, okError=%v, err1=%v, err2=%v", okValue, okError, err1, err2)
	return ErrMissing //something missing in logic
}

func joinErrors(err1, err2 error) error {
	if err1 == ErrIsRPCError {
		return err2
	}
	return err1
}
