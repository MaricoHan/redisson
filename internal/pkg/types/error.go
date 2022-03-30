package types

import "fmt"

const (
	RootCodeSpace = "NFTP-OPEN-API"
)

const (
	// InternalFailed		error code
	InternalFailed                 = "INTERNAL_ERROR"
	AuthenticationFailed           = "FORBIDDEN"
	ClientParamsError              = "PARAMS_ERROR"
	ConnectionChainFailed          = "CONNECTION_CHAIN_FAILED"
	FrequentRequestsNotSupports    = "FREQUENT_REQUESTS_NOT_SUPPORTS"
	NftClassStatusAbnormal         = "NFT_CLASS_STATUS_ABNORMAL"
	NftStatusAbnormal              = "NFT_STATUS_ABNORMAL"
	NotFound                       = "NOT_FOUND"
	MaximumLimitExceeded           = "MAXIMUM_LIMIT_EXCEEDED"
	StructureSignTransactionFailed = "STRUCTURE_SIGN_TRANSACTION_FAILED"
	TxStatusSuccess                = "TX_STATUS_SUCCESS"
	TxStatusPending                = "TX_STATUS_PENDING"
	TxStatusUndo                   = "TX_STATUS_UNDO"
	StructureSendTransactionFailed = "STRUCTURE_SEND_TRANSACTION_FAILED"
	ModuleFailed                   = "MODULE_ERROR"
	AccountFAiled                  = "ACCOUNT_ERROR"

	// ErrOffset		error msg handle
	ErrOffset         = "offset format error"
	ErrOffsetInt      = "offset cannot be less than 0"
	ErrLimitParam     = "limit format error"
	ErrLimitParamInt  = "limit must be between 1 and 50"
	ErrCountLen       = "count length error"
	ErrStartDate      = "startDate format error"
	ErrEndDate        = "endDate format error"
	ErrDate           = "endDate before startDate"
	ErrRecipient      = "recipient is a required field"
	ErrRecipientAddr  = "the recipient address does not meet the specification of the current chain"
	ErrRecipientLen   = "recipient length error"
	ErrNftId          = "nft_id format error"
	ErrNftIdLen       = "nft_id is a required field"
	ErrNftIdString    = "nft_id cannot be nil"
	ErrRecipients     = "recipients is a required field"
	ErrName           = "name is a required field"
	ErrNameLen        = "name length error"
	ErrSymbolLen      = "symbol length error"
	ErrDescriptionLen = "description length error"
	ErrUri            = "uri format error"
	ErrUriLen         = "uri length error"
	ErrUriHashLen     = "uriHash length error"
	ErrDataLen        = "data length error"
	ErrOwner          = "owner is a required field"
	ErrOwnerLen       = "owner length error"
	ErrSortBy         = "sortBy is invalid"
	ErrIndices        = "indices format error"
	ErrIndicesLen     = "indices is a required field"
	ErrOperation      = "operation is invalid"
	ErrModule         = "module is invalid"
	ErrAmountInt      = "amount must be between 1 and 100"
	ErrRepeat         = "index is repeat"
	ErrClientParams   = "client params error"
	ErrUriChain       = "uri cannot be modified"
	ErrAccountCount   = "number that can be created exceeds the limit"

	// ErrSelfTransfer		error msg service
	ErrSelfTransfer      = "recipient cannot be owner"
	ErrRecipientFound    = "recipient not found"
	ErrNftFound          = "nft not found"
	ErrNftStatusMsg      = "nft status is invalid"
	ErrNftClassStatusMsg = "nft class status is invalid"
	ErrOwnerFound        = "owner not found"
	ErrDIDAlreadyExists  = "Authority: Account already exists!"
	ErrOutOfGas          = "out of gas"
)

var (
	ErrInternal        = Register(RootCodeSpace, InternalFailed, "internal")
	ErrAuthenticate    = Register(RootCodeSpace, AuthenticationFailed, "authentication failed")
	ErrParams          = Register(RootCodeSpace, ClientParamsError, ErrClientParams)
	ErrChainConn       = Register(RootCodeSpace, ConnectionChainFailed, "connection chain failed")
	ErrIdempotent      = Register(RootCodeSpace, FrequentRequestsNotSupports, "frequent requests not supports")
	ErrNftClassStatus  = Register(RootCodeSpace, NftClassStatusAbnormal, ErrNftClassStatusMsg)
	ErrNftStatus       = Register(RootCodeSpace, NftStatusAbnormal, ErrNftStatusMsg)
	ErrNotFound        = Register(RootCodeSpace, NotFound, "resource not found")
	ErrLimit           = Register(RootCodeSpace, MaximumLimitExceeded, "maximum limit exceeded")
	ErrBuildAndSign    = Register(RootCodeSpace, StructureSignTransactionFailed, "build and sign transaction failed")
	ErrBuildAndSend    = Register(RootCodeSpace, StructureSendTransactionFailed, "build and send transaction failed")
	ErrTXStatusSuccess = Register(RootCodeSpace, TxStatusSuccess, "tx transaction success")
	ErrTXStatusPending = Register(RootCodeSpace, TxStatusPending, "tx transaction is in progress, please wait")
	ErrTXStatusUndo    = Register(RootCodeSpace, TxStatusUndo, "tx transaction not executed, please wait")
	ErrModules         = Register(RootCodeSpace, ModuleFailed, ErrModule)
	ErrAccount         = Register(RootCodeSpace, AccountFAiled, ErrAccountCount)
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
