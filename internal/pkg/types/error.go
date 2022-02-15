package types

import "fmt"

const (
	RootCodeSpace = "NFTP-OPEN-API"
)

const (
	InternalFailed                 = "INTERNAL_FAILED"
	AuthenticationFailed           = "AUTHENTICATION_FAILED"
	ClientParamsError              = "CLIENT_PARAMS_ERROR"
	ConnectionChainFailed          = "CONNECTION_CHAIN_FAILED"
	FrequentRequestsNotSupports    = "FREQUENT_REQUESTS_NOT_SUPPORTS"
	NftclassStatusAbnormal         = "NFTCLASS_STATUS_ABNORMAL"
	NftStatusAbnormal              = "NFT_STATUS_ABNORMAL"
	NotFound                       = "NOT_FOUND"
	MaximumLimitExceeded           = "MAXIMUM_LIMIT_EXCEEDED"
	StructureSignTransactionFailed = "STRUCTURE_SIGN_TRANSACTION_FAILED"
	TxStatusSuccess                = "TX_STATUS_SUCCESS"
	TxStatusPending                = "TX_STATUS_PENDING"
	TxStatusUndo                   = "TX_STATUS_UNDO"
	StructureSendTransactionFailed = "STRUCTURE_SEND_TRANSACTION_FAILED"
)

var (
	ErrInternal        = Register(RootCodeSpace, InternalFailed, "internal")
	ErrAuthenticate    = Register(RootCodeSpace, AuthenticationFailed, "failed to authentication")
	ErrParams          = Register(RootCodeSpace, ClientParamsError, "failed to client params")
	ErrChainConn       = Register(RootCodeSpace, ConnectionChainFailed, "failed to connection chain")
	ErrIdempotent      = Register(RootCodeSpace, FrequentRequestsNotSupports, "failed to idempotent")
	ErrNftClassStatus  = Register(RootCodeSpace, NftclassStatusAbnormal, "the nft class status is invalid")
	ErrNftStatus       = Register(RootCodeSpace, NftStatusAbnormal, "the nft status is invalid")
	ErrNotFound        = Register(RootCodeSpace, NotFound, "resource not found")
	ErrLimit           = Register(RootCodeSpace, MaximumLimitExceeded, "maximum limit exceeded")
	ErrBuildAndSign    = Register(RootCodeSpace, StructureSignTransactionFailed, "failed to build and sign")
	ErrBuildAndSend    = Register(RootCodeSpace, StructureSendTransactionFailed, "failed to build and send")
	ErrTXStatusSuccess = Register(RootCodeSpace, TxStatusSuccess, "tx transaction success")
	ErrTXStatusPending = Register(RootCodeSpace, TxStatusPending, "tx transaction is in progress, please wait")
	ErrTXStatusUndo    = Register(RootCodeSpace, TxStatusUndo, "tx transaction not executed, please wait")
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
