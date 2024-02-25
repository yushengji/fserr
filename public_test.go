package fserr

import (
	"errors"
	"git.sxidc.com/service-supports/fslog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type TestPublicSuite struct {
	suite.Suite
	outerErr,
	newErr, newFmtErr, stackErr,
	wrapErr, wrapFmtErr, wrapOuterErr, wrapNilErr, wrapEmptyErr,
	codeErr, codeOptionErr, codeOuterErr, codeNilErr error
	errBasicMsg string
}

func (s *TestPublicSuite) SetupTest() {
	s.outerErr = errors.New("outer error")
	s.errBasicMsg = "basic error"
	NewOK(ErrBasic, s.errBasicMsg)

	// new
	s.newErr = New("new error")
	s.newFmtErr = New("new %s", "error")

	// wrap
	s.wrapErr = Wrap(s.newErr, "wrap error")
	s.wrapFmtErr = Wrap(s.newErr, "wrap %s", "error")
	s.wrapOuterErr = Wrap(s.outerErr, "wrap error")
	s.wrapNilErr = Wrap(nil, "wrap error")
	s.wrapEmptyErr = Wrap(s.newErr, " ")

	// code
	s.codeErr = WithCode(s.newErr, ErrBasic)
	s.codeOptionErr = WithCode(s.newErr, ErrBasic, WithMessage("cover message"))
	s.codeOuterErr = WithCode(s.outerErr, ErrBasic)
	s.codeNilErr = WithCode(nil, ErrBasic)

	s.stackErr = WithStack(errors.New("stack error"))
}

func (s *TestPublicSuite) TestNew() {
	s.Equal("new error", s.newErr.Error())
	s.Equal("new error", s.newFmtErr.Error())
}

func (s *TestPublicSuite) TestWrap() {
	s.Equal("wrap error", s.wrapErr.Error())
	s.Equal("wrap error", s.wrapFmtErr.Error())
	s.Equal("wrap error", s.wrapOuterErr.Error())
	s.Equal("new error", s.wrapEmptyErr.Error())
	s.Nil(s.wrapNilErr)
}

func (s *TestPublicSuite) TestWithCode() {
	s.Equal(s.errBasicMsg, s.codeErr.Error())
	s.Equal("cover message", s.codeOptionErr.Error())
	s.Equal("basic error", s.codeOuterErr.Error())
	s.Equal("basic error", s.codeNilErr.Error())
}

func (s *TestPublicSuite) TestUnWrap() {
	s.Equal(s.newErr, UnWrap(s.newErr))

	s.Equal(s.newErr, UnWrap(s.wrapErr))
	s.Equal(s.outerErr, UnWrap(s.wrapOuterErr))
	s.Nil(UnWrap(s.wrapNilErr))

	s.Equal(s.newErr, UnWrap(s.codeErr))
	s.Equal(s.outerErr, UnWrap(s.codeOuterErr))
	s.Equal(nil, UnWrap(s.codeNilErr))
}

func (s *TestPublicSuite) TestIs() {
	s.True(Is(s.newErr, s.newErr))
	s.True(Is(s.wrapErr, s.newErr))
	s.True(Is(s.codeErr, s.newErr))

	s.False(Is(s.newErr, s.outerErr))
	s.False(Is(s.wrapErr, s.outerErr))
	s.False(Is(s.codeErr, s.outerErr))
}

func (s *TestPublicSuite) TestAs() {
	var originErr *fundamental
	s.True(As(s.wrapErr, &originErr))
	s.Equal(UnWrap(s.wrapErr), originErr)

	var codeErr *withCode
	s.False(As(s.wrapErr, &codeErr))
}

func (s *TestPublicSuite) TestOuterMsg() {
	s.Equal("new error", outerMsg(s.newErr))
	s.Equal("wrap error", outerMsg(s.wrapErr))
	s.Equal(s.errBasicMsg, outerMsg(s.codeErr))
	s.Equal("outer error", outerMsg(s.outerErr))
	s.Equal("", outerMsg(nil))
}

func (s *TestPublicSuite) TestParseCode() {
	code := ParseCode(s.codeErr)
	s.Equal(ErrBasic, code.BusinessCode)
	s.Equal(http.StatusOK, code.HttpCode)
	s.Equal(s.errBasicMsg, code.Msg)
	s.Equal(s.newErr, code.cause)

	err := ParseCode(Wrap(s.codeErr, "wrap error"))
	s.Equal(ErrBasic, err.BusinessCode)
	s.Equal(http.StatusOK, err.HttpCode)
	s.Equal("wrap error", err.Msg)
	s.Equal(s.newErr, err.cause)
}

func (s *TestPublicSuite) TestIsCode() {
	s.True(IsCode(s.codeErr, ErrBasic))
	s.False(IsCode(s.codeErr, ErrDb))
	s.False(IsCode(s.wrapErr, ErrBasic))
}

func (s *TestPublicSuite) TestSetServiceCode() {
	SetAppCode(3)
	NewOK(2001, "ok")
	err := WithCode(nil, 2001)
	s.Equal(32001, ParseCode(err).BusinessCode)
}

func (s *TestPublicSuite) TestWithStack() {
	s.Equal("stack error", s.stackErr.Error())
	fslog.Error(s.stackErr)
}

func TestPublic(t *testing.T) {
	suite.Run(t, &TestPublicSuite{})
}

func TestServiceCode(t *testing.T) {
	SetAppCode(1)
	NewInternalError(ErrBasic, "basic error")
	code := ParseCode(WithCode(nil, ErrBasic))
	assert.Equal(t, "basic error", code.Msg)
	assert.Equal(t, 10001, code.BusinessCode)
	assert.Equal(t, http.StatusInternalServerError, code.HttpCode)

	SetAppCode(2)
	NewBadRequest(ErrDb, "db error")
	code = ParseCode(WithCode(nil, ErrDb))
	assert.Equal(t, "db error", code.Msg)
	assert.Equal(t, 20002, code.BusinessCode)
	assert.Equal(t, http.StatusBadRequest, code.HttpCode)

	SetAppCode(3)
	NewInternalError(ErrParam, "param error")
	code = ParseCode(WithCode(nil, ErrParam))
	assert.Equal(t, "param error", code.Msg)
	assert.Equal(t, 30003, code.BusinessCode)
	assert.Equal(t, http.StatusInternalServerError, code.HttpCode)

	SetAppCode(4)
	NewConflict(ErrRetry, "retry error")
	code = ParseCode(WithCode(nil, ErrRetry))
	assert.Equal(t, 40004, code.BusinessCode)
	assert.Equal(t, http.StatusConflict, code.HttpCode)
	assert.Equal(t, "retry error", code.Msg)
}
