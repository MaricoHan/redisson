package types

import "fmt"

const RootCodeSpace = "nftp-open-api"

var (
	ErrInternal               = Register(RootCodeSpace, "1", "internal")
	ErrAuthenticate           = Register(RootCodeSpace, "2", "failed to authentication ")
	ErrParams                 = Register(RootCodeSpace, "3", "failed to client params")
	ErrMysqlConn              = Register(RootCodeSpace, "4", "failed to connection mysql")
	ErrRedisConn              = Register(RootCodeSpace, "5", "failed to connection redis")
	ErrChainConn              = Register(RootCodeSpace, "6", "failed to connection chain")
	ErrAccountCreate          = Register(RootCodeSpace, "7", "failed to create chain account")
	ErrGetAccountDetails      = Register(RootCodeSpace, "8", "failed to get chain account details")
	ErrNftClassCreate         = Register(RootCodeSpace, "9", "failed to create nft class")
	ErrNftClassesGet          = Register(RootCodeSpace, "10", "failed to get nft class list")
	ErrNftClassDetailsGet     = Register(RootCodeSpace, "11", "failed to get nft class details")
	ErrNftCreate              = Register(RootCodeSpace, "12", "failed to create nft class")
	ErrNftGet                 = Register(RootCodeSpace, "13", "failed to get nft list")
	ErrNftDetailsGet          = Register(RootCodeSpace, "14", "failed to get nft details")
	ErrNftOptionHistoryGet    = Register(RootCodeSpace, "15", "failed to get nft option history")
	ErrNftEdit                = Register(RootCodeSpace, "16", "failed to edit nft")
	ErrNftBatchEdit           = Register(RootCodeSpace, "17", "failed to batch edit nft")
	ErrNftBurn                = Register(RootCodeSpace, "18", "failed to burn nft")
	ErrNftBatchBurn           = Register(RootCodeSpace, "19", "failed to batch burn nft")
	ErrTxResult               = Register(RootCodeSpace, "20", "failed to get tx result")
	ErrIdempotent             = Register(RootCodeSpace, "21", "failed to idempotent")
	ErrNftTooMany             = Register(RootCodeSpace, "22", "The maximum number of NFTs to edit is 50")
	ErrNftMissing             = Register(RootCodeSpace, "23", "Cannot find the NFT")
	ErrNftBurnPend            = Register(RootCodeSpace, "24", "The platform has received the destruction request and put it on the chain, but it has not been packaged and confirmed")
	ErrNftClassTransfer       = Register(RootCodeSpace, "25", "failed to transfer nft class")
	ErrBuildAndSign           = Register(RootCodeSpace, "26", "failed to build and sign")
	ErrNftTransfer            = Register(RootCodeSpace, "27", "failed to transfer nft")
	ErrNftBatchTransfer       = Register(RootCodeSpace, "28", "failed to batch transfer nft")
	ErrNotOwner               = Register(RootCodeSpace, "29", "sender is not the nft‘s owner")
	ErrNoPermission           = Register(RootCodeSpace, "30", "sender is not the one of the app‘s accounts")
	ErrNftClassesSet          = Register(RootCodeSpace, "31", "failed to set nft class")
	ErrIndexFormat            = Register(RootCodeSpace, "32", "Index format is invalid, must be unsigned numeric type")
	ErrTxMsgInsert            = Register(RootCodeSpace, "33", "failed to insert ttx")
	ErrTxMsgGet               = Register(RootCodeSpace, "34", "failed to get ttx")
	ErrGetTx                  = Register(RootCodeSpace, "35", "failed to get tx by hash")
	ErrGetNftOperationDetails = Register(RootCodeSpace, "36", "failed to get nft operation details")
	ErrNftStatus              = Register(RootCodeSpace, "37", "One of these NFTs does not exist or its status is not active")
	ErrIndicesFormat          = Register(RootCodeSpace, "38", "Indices format is invalid, must be unsigned numeric type,such as:1,2,3,4...")
	ErrNftClassStatus         = Register(RootCodeSpace, "39", "One of these NFT Class does not exist or its status is not active")
	ErrClassStatus            = Register(RootCodeSpace, "40", "nftClass does not exist or its status is not active")
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

func UpdateDescription(codeSpace string, code string, description string) *AppError {
	e := getUsedErrorCodes(codeSpace, code)
	if e == nil {
		panic(fmt.Sprintf("error with code %d No registered: %q", code, e.desc))
	}
	e.desc = description
	return e
}
