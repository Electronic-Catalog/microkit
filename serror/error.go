package serror

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

type SError interface {
	error
	fmt.Stringer

	GetMessage() string
	WithMessageKey(key string) SError

	GetMessageKey() string
	WithMessage(msg string) SError

	GetStatusCode() StatusCode
	WithStatusCode(statusCode StatusCode) SError

	GetTrace() string
	WithTrace(trace string) SError

	GetPayload() []byte
	WithPayload(payload []byte) SError

	GetParams() []interface{}
	WithParams(params []interface{}) SError

	GetDevTrace() string
	StringDetail() string
	StringDetailSimple() string

	ReloadDevTrace(skip int) SError
	Clone() SError
}

func New() SError {

	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	s := sError{DevTrace: fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function)}

	return s
}

func Is(e error) (SError, bool) {
	if se, ok := e.(SError); ok {
		return se, true
	}

	return nil, false
}

func Equals(ea error, eb error) bool {
	a, oka := Is(ea)
	if !oka {
		return false
	}

	b, okb := Is(eb)
	if !okb {
		return false
	}

	ae, ok := a.(sError)
	if !ok {
		return false
	}

	be, ok := b.(sError)
	if !ok {
		return false
	}

	if strings.Compare(ae.GetMessageKey(), be.GetMessageKey()) != 0 {
		return false
	}

	return true
}

type sError struct {
	StatusCode      StatusCode    `json:"status_code"`
	PayloadValue    []byte        `json:"payload_value"`
	TextFormatValue []string      `json:"text_format_value"`
	Message         string        `json:"message"`
	MessageKey      string        `json:"message_key"`
	Trace           string        `json:"trace"`
	DevTrace        string        `json:"dev_trace"`
	Params          []interface{} `json:"params"`
}

func (e sError) Clone() SError {

	se := New().WithMessageKey(e.GetMessageKey()).
		WithMessage(e.GetMessage()).
		WithTrace(e.GetTrace()).
		WithStatusCode(e.GetStatusCode()).
		WithPayload(e.GetPayload()).
		WithParams(e.GetParams())

	return se
}

func (e sError) ReloadDevTrace(skip int) SError {

	pc := make([]uintptr, 15)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	e.DevTrace = fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function)

	return e

}

func (e sError) StringDetailSimple() string {

	s := e.StringDetail()
	s = strings.ReplaceAll(s, `"`, "")
	s = strings.ReplaceAll(s, `\n`, "")
	return s

}

func (e sError) GetParams() []interface{} {

	return e.Params
}

func (e sError) WithParams(params []interface{}) SError {

	e.Params = params
	return e
}

func (e sError) GetMessage() string {
	return e.Message
}

func (e sError) GetMessageKey() string {

	return e.MessageKey
}

func (e sError) GetTrace() string {

	return e.Trace
}

func (e sError) GetDevTrace() string {
	return e.DevTrace
}

func (e sError) GetStatusCode() StatusCode {

	return e.StatusCode
}

func (e sError) GetPayload() []byte {
	return e.PayloadValue
}

func (e sError) WithMessage(msg string) SError {

	e.Message = msg
	return e
}

func (e sError) WithStatusCode(statusCode StatusCode) SError {

	e.StatusCode = statusCode

	return e
}

func (e sError) WithMessageKey(key string) SError {
	e.MessageKey = key
	return e
}

func (e sError) WithTrace(trace string) SError {
	e.Trace = trace
	return e
}

func (e sError) WithPayload(payload []byte) SError {
	e.PayloadValue = payload
	return e
}

func (e sError) Error() string {
	return e.String()
}

func (e sError) String() string {
	str := fmt.Sprintf("error --> status_code: %d , message: %v, trace: %s", e.StatusCode, e.Message, e.Trace)
	return str
}

func (e sError) StringDetail() string {
	marshal, err := json.Marshal(e)
	if err != nil {
		return e.String()
	}

	return string(marshal)
}
