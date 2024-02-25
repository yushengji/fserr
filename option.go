package fserr

type Option func(*withCode)

// WithMessage 替换错误码默认的提示信息
func WithMessage(msg string) Option {
	return func(w *withCode) {
		w.Msg = msg
	}
}
