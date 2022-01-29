package types

import "fmt"

const RootCodeSpace = "nftp-open-api"

var (
	// 增删改查统一错误码
	ErrCreate   = Register(RootCodeSpace, "1001", "CREATE_FAILED")
	ErrQuery    = Register(RootCodeSpace, "1002", "QUERY_FAILED")
	ErrEdit     = Register(RootCodeSpace, "1003", "EDIT_FAILED")
	ErrBurn     = Register(RootCodeSpace, "1004", "BURN_FAILED")
	ErrTransfer = Register(RootCodeSpace, "1005", "TRANSFER_FAILED")

	ErrInternal         = Register(RootCodeSpace, "1", "INTERNAL_FAILED")
	ErrAuthenticate     = Register(RootCodeSpace, "2", "AUTHENTICATION_FAILED")
	ErrParams           = Register(RootCodeSpace, "3", "CLIENT_PARAMS_ERROR")
	ErrChainConn        = Register(RootCodeSpace, "4", "CONNECTION_CHAIN_FAILED")
	ErrIdempotent       = Register(RootCodeSpace, "5", "FREQUENT_REQUESTS_NOT_SUPPORTS")
	ErrNftClassNotFound = Register(RootCodeSpace, "6", "NFTCLASS_NOT_EXIST")
	ErrNftClassStatus   = Register(RootCodeSpace, "7", "NFTCLASS_STATUS_ABNORMAL")
	ErrNftNotFound      = Register(RootCodeSpace, "8", "NFT_NOT_EXIST")
	ErrNftStatus        = Register(RootCodeSpace, "9", "NFT_STATUS_ABNORMAL")
	ErrBatch            = Register(RootCodeSpace, "10", "BATCH_OPERATION_NFT_ERROR")
	ErrLimit            = Register(RootCodeSpace, "11", "MAXIMUM_LIMIT_50")
	ErrNotOwner         = Register(RootCodeSpace, "12", "NOT_OWNER_ACCOUNT")
	ErrNoPermission     = Register(RootCodeSpace, "13", "NOT_APP_ OF_ACCOUNT")
	ErrBuildAndSign     = Register(RootCodeSpace, "14", "STRUCTURE_SIGN_TRANSACTION_FAILED")

	ErrIndexFormat     = Register(RootCodeSpace, "15", "Index format is invalid, must be unsigned numeric type")
	ErrIndicesFormat   = Register(RootCodeSpace, "16", "Indices format is invalid, must be unsigned numeric type,such as:1,2,3,4...")
	ErrRepeated        = Register(RootCodeSpace, "17", "Please do not fill in duplicate NFT in the request parameters")
	ErrTXStatusSuccess = Register(RootCodeSpace, "18", "tx transaction success")
	ErrTXStatusPending = Register(RootCodeSpace, "19", "tx transaction is in progress, please wait")
	ErrTXStatusUndo    = Register(RootCodeSpace, "20", "tx transaction not executed, please wait")
)

var usedErrorCodes = map[string]*AppError{}

func getUsedErrorCodes(codeSpace string, code string) *AppError {
	return usedErrorCodes[appErrorID(codeSpace, code)]
}

func setUsedErrorCodes(err *AppError) {
	usedErrorCodes[appErrorID(err.codeSpace, err.code)] = err
}

func appErrorID(codeSpace string, code string) string {
	return fmt.Sprintf("%s:%d", codeSpace, code)
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
		panic(fmt.Sprintf("error with code %d is already registered: %q", code, e.desc))
	}

	err := NewAppError(codeSpace, code, description)
	setUsedErrorCodes(err)

	return err
}
