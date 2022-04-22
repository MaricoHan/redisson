package constant

import "fmt"

const (
	RootCodeSpace = "NFTP-OPEN-API"
)

const (
	// InternalFailed		error code
	InternalFailed                 = "INTERNAL_ERROR"
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
	TimestampTimeout               = "Timestamp_Timeout"

	// ErrOffset		error msg handle
	ErrName         = "name is a required field"
	ErrClientParams = "client params error"

	// ErrSelfTransfer		error msg service
	ErrNftStatusMsg = "nft status is invalid"
	ErrOutOfGas     = "out of gas"
	ErrApikey       = "api_key is not exist"
	ErrOrderType    = "order_type is invalid"

	ErrInternalFailed = "internal error"
)

var (
	ErrInternal     = Register(RootCodeSpace, InternalFailed, "internal")
	ErrAuthenticate = Register(RootCodeSpace, AuthenticationFailed, "authentication failed")
	ErrParams       = Register(RootCodeSpace, ClientParamsError, ErrClientParams)
	ErrIdempotent   = Register(RootCodeSpace, FrequentRequestsNotSupports, "frequent requests not supports")
	ErrNftStatus    = Register(RootCodeSpace, NftStatusAbnormal, ErrNftStatusMsg)
	ErrTimestamp    = Register(RootCodeSpace, TimestampTimeout, "timestamp is timeout")
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
	if e := getUsedErrorCodes(codeSpace, code); e != nil {
		panic(fmt.Sprintf("error with code %s is already registered: %q", code, e.desc))
	}

	err := NewAppError(codeSpace, code, description)
	setUsedErrorCodes(err)

	return err
}
