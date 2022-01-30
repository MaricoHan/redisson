package types

import "fmt"

const (
	RootCodeSpace = "NFTP-OPEN-API"
)

const (
	QueryFailed                    = "QUERY_FAILED"
	CreateFailed                   = "CREATE_FAILED"
	EditFailed                     = "EDIT_FAILED"
	TransferFailed                 = "TRANSFER_FAILED"
	InternalFailed                 = "INTERNAL_FAILED"
	AuthenticationFailed           = "AUTHENTICATION_FAILED"
	ClientParamsError              = "CLIENT_PARAMS_ERROR"
	ConnectionChainFailed          = "CONNECTION_CHAIN_FAILED"
	FrequentRequestsNotSupports    = "FREQUENT_REQUESTS_NOT_SUPPORTS"
	NftclassNotExist               = "NFTCLASS_NOT_EXIST"
	NftclassStatusAbnormal         = "NFTCLASS_STATUS_ABNORMAL"
	NftNotExist                    = "NFT_NOT_EXIST"
	NftStatusAbnormal              = "NFT_STATUS_ABNORMAL"
	MaximumLimitExceeded           = "MAXIMUM_LIMIT_EXCEEDED"
	NotOwnerAccount                = "NOT_OWNER_ACCOUNT"
	NotAppOfAccount                = "NOT_APP_OF_ACCOUNT"
	StructureSignTransactionFailed = "STRUCTURE_SIGN_TRANSACTION_FAILED"
	RepeatError                    = "REPEAT_ERROR"
	TxStatusSuccesss               = "TX_STATUS_SUCCESSS"
	TxStatusPending                = "TX_STATUS_PENDING"
	TxStatusUndo                   = "TX_STATUS_UNDO"
	StructureSendTransactionFailed = "STRUCTURE_SEND_TRANSACTION_FAILED"
)

var (
	ErrCreate   = Register(RootCodeSpace, CreateFailed, "failed to create")
	ErrQuery    = Register(RootCodeSpace, QueryFailed, "failed to query")
	ErrEdit     = Register(RootCodeSpace, EditFailed, "failed to edit")
	ErrBurn     = Register(RootCodeSpace, EditFailed, "failed to burn")
	ErrTransfer = Register(RootCodeSpace, TransferFailed, "failed to transfer")

	ErrInternal         = Register(RootCodeSpace, InternalFailed, "internal")
	ErrAuthenticate     = Register(RootCodeSpace, AuthenticationFailed, "failed to authentication")
	ErrParams           = Register(RootCodeSpace, ClientParamsError, "failed to client params")
	ErrChainConn        = Register(RootCodeSpace, ConnectionChainFailed, "failed to connection chain")
	ErrIdempotent       = Register(RootCodeSpace, FrequentRequestsNotSupports, "failed to idempotent")
	ErrNftClassNotFound = Register(RootCodeSpace, NftclassNotExist, "the NFT Class does not exist")
	ErrNftClassStatus   = Register(RootCodeSpace, NftclassStatusAbnormal, "the NFT Class status is invalid")
	ErrNftNotFound      = Register(RootCodeSpace, NftNotExist, "the NFT does not exist")
	ErrNftStatus        = Register(RootCodeSpace, NftStatusAbnormal, "the NFT status is invalid")
	ErrLimit            = Register(RootCodeSpace, MaximumLimitExceeded, "")
	ErrNotOwner         = Register(RootCodeSpace, NotOwnerAccount, "This account is not the owner account")
	ErrNoPermission     = Register(RootCodeSpace, NotAppOfAccount, "This account is not an in-app account")
	ErrBuildAndSign     = Register(RootCodeSpace, StructureSignTransactionFailed, "failed to build and sign")
	ErrRepeated         = Register(RootCodeSpace, RepeatError, "Please do not fill in duplicate NFT in the request parameters")
	ErrTXStatusSuccess  = Register(RootCodeSpace, TxStatusSuccesss, "tx transaction success")
	ErrTXStatusPending  = Register(RootCodeSpace, TxStatusPending, "tx transaction is in progress, please wait")
	ErrTXStatusUndo     = Register(RootCodeSpace, TxStatusUndo, "tx transaction not executed, please wait")
	ErrBuildAndSend     = Register(RootCodeSpace, StructureSendTransactionFailed, "failed to build and send")
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
