package serror

type StatusCode int

const (
	ScOk StatusCode = 0
)

func (s StatusCode) Success() bool {
	return s == ScOk
}

func (s StatusCode) Failed() bool {
	return !s.Success()
}

func (s StatusCode) ToInt32() int32 {
	return int32(s)
}

func NewStatusCode(status int) StatusCode {
	return StatusCode(status)
}
