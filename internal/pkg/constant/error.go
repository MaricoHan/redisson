package constant

import (
	"fmt"
	"google.golang.org/grpc/codes"
)

const (
	RootCodeSpace                     = "NFTP-OPEN-API"
	MtCodeSpace                       = "MT"
	AuthCodeSpace                     = "AUTH"
	UpstreamInternalFailed codes.Code = 10086
)

const (
	// InternalFailed		error code
	InternalFailed                 = "INTERNAL_ERROR"
	UpstreamInternalFaileds        = "UPSTREAM_INTERNAL_ERROR"
	AuthenticationFailed           = "FORBIDDEN"
	ClientParamsError              = "PARAMS_ERROR"
	FrequentRequestsNotSupports    = "FREQUENT_REQUESTS_NOT_SUPPORTS"
	NftClassStatusAbnormal         = "NFT_CLASS_STATUS_ABNORMAL"
	NftStatusAbnormal              = "NFT_STATUS_ABNORMAL"
	NotFound                       = "NOT_FOUND"
	MaximumLimitExceeded           = "MAXIMUM_LIMIT_EXCEEDED"
	StructureSignTransactionFailed = "STRUCTURE_SIGN_TRANSACTION_FAILED"
	ModuleFailed                   = "MODULE_ERROR"
	AccountFailed                  = "ACCOUNT_ERROR"
	TimestampTimeout               = "TIMESTAMP_TIME"
	DuplicateRequest               = "DUPLICATE_REQUEST"
	UnSupported                    = "NOT_IMPLEMENTED"

	// ErrOffset		error msg handle
	ErrName         = "name is a required field"
	ErrClientParams = "client params error"

	// ErrSelfTransfer		error msg service
	ErrNftStatusMsg = "nft status is invalid"
	ErrOutOfGas     = "out of gas"
	ErrApikey       = "api_key is not exist"
	ErrOrderType    = "order_type is invalid"

	ErrInternalFailed = "internal error"

	ErrUpstreamEntity = "returned entity does not conform to the rule"
	ErrNotFound       = "not found"
	ErrInvalidValue   = "invalid %s value"

	ErrProjectOrUserNotFound      = "user or project not exists"
	ErrServiceRedirectUrlNotFound = "service redirect address is not configured"
	ErrUpstreamInternal           = "upstream service is abnormal"
	ErrUpstreamForbidden          = "customer has no right to access this api"
	// ErrUpstreamInternalEntity     = Register(AuthCodeSpace, UpstreamInternalFaileds, ErrUpstreamEntity)
	ErrAuthVerifyExists  = "invalid exists value"
	ErrAuthUserAddress   = "invalid address value"
	ErrAuthUserChainName = "invalid chain_name value"
)

var (
	ErrInternal             = Register(RootCodeSpace, InternalFailed, "internal")
	ErrAuthenticate         = Register(RootCodeSpace, AuthenticationFailed, "authentication failed")
	ErrParams               = Register(RootCodeSpace, ClientParamsError, ErrClientParams)
	ErrIdempotent           = Register(RootCodeSpace, FrequentRequestsNotSupports, "frequent requests not supports")
	ErrNftStatus            = Register(RootCodeSpace, NftStatusAbnormal, ErrNftStatusMsg)
	ErrTimestamp            = Register(RootCodeSpace, TimestampTimeout, "timestamp is timeout")
	ErrDuplicate            = Register(RootCodeSpace, DuplicateRequest, "duplicate request")
	ErrUnSupported          = Register(MtCodeSpace, UnSupported, "not implemented")
	ErrUnmanagedUnSupported = Register(RootCodeSpace, UnSupported, "not implemented")
)

var usedErrorCodes = map[string]*AppError{}

func getUsedErrorCodes(codeSpace string, code string) *AppError {
	return usedErrorCodes[appErrorID(codeSpace, code)]
}

func setUsedErrorCodes(err *AppError) {
	usedErrorCodes[appErrorID(err.codeSpace, err.code)] = err
}

func appErrorID(codeSpace string, code string) string {
	return fmt.Sprintf("%s:%s", codeSpace, code)
}

type IError interface {
	error
	Code() string
	CodeSpace() string
}

type AppError struct {
	codeSpace string
	code      string
	desc      string
}

func NewAppError(codeSpace string, code string, desc string) *AppError {
	return &AppError{codeSpace: codeSpace, code: code, desc: desc}
}

func (e AppError) Error() string {
	return e.desc
}

func (e AppError) Code() string {
	return e.code
}

func (e AppError) CodeSpace() string {
	return e.codeSpace
}

func Register(codeSpace string, code string, description string) *AppError {
	// if e := getUsedErrorCodes(codeSpace, code); e != nil {
	// 	panic(fmt.Sprintf("error with code %s is already registered: %q", code, e.desc))
	// }

	err := NewAppError(codeSpace, code, description)
	setUsedErrorCodes(err)

	return err
}
