package fserr

import (
	"git.sxidc.com/service-supports/fslog"
	"github.com/puzpuzpuz/xsync"
)

var codeMap *xsync.MapOf[int, ErrCode]

var defaultErrCode = ErrCode{
	HttpCode: 200,
}

func init() {
	codeMap = xsync.NewIntegerMapOf[int, ErrCode]()
}

func register(code ErrCode) {
	if _, ok := codeMap.Load(code.BusinessCode); ok {
		fslog.With("code", code.BusinessCode).
			With("message", code.Message).
			Warn("duplicate business code")
	}
	codeMap.Store(code.BusinessCode, code)
}

func getCode(business int) ErrCode {
	code, ok := codeMap.Load(business)
	if ok {
		return code
	}
	ret := defaultErrCode
	ret.BusinessCode = business
	return ret
}
