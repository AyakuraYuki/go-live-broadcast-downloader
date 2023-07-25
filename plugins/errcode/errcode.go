package errcode

type ServerError struct {
	ErrCode int64
	ErrMsg  string
}

func NewServerError(errCode int64, errMsg string) ServerError {
	return ServerError{ErrCode: errCode, ErrMsg: errMsg}
}

func (e ServerError) Error() string {
	return e.ErrMsg
}
