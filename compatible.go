package fserr

// TransferFromError 将错误转为带栈帧的错误
// Deprecated: 请使用 New、 Wrap、 WithCode、 WithStack
func TransferFromError(err error) error {
	return WithStack(err)
}
