package errcode

var (
	ScOK = NewServerError(0, "ok")

	ErrOther = NewServerError(500, "其他异常")
)
