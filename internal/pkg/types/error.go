package types

import "fmt"

const RootCodeSpace = "nftp-open-api"

var (
	ErrInternal            = Register(RootCodeSpace, "1", "internal")
	ErrAuthenticate        = Register(RootCodeSpace, "2", "failed to authentication ")
	ErrParams              = Register(RootCodeSpace, "3", "failed to client params")
	ErrMysqlConn           = Register(RootCodeSpace, "4", "failed to connection mysql")
	ErrRedisConn           = Register(RootCodeSpace, "5", "failed to connection redis")
	ErrChainConn           = Register(RootCodeSpace, "6", "failed to connection chain")
	ErrAccountCreate       = Register(RootCodeSpace, "7", "failed to create chain account")
	ErrGetAccountDetails   = Register(RootCodeSpace, "8", "failed to get chain account details")
	ErrNftClassCreate      = Register(RootCodeSpace, "9", "failed to create nft class")
	ErrNftClassesGet       = Register(RootCodeSpace, "10", "failed to get nft class list")
	ErrNftClassDetailsGet  = Register(RootCodeSpace, "11", "failed to get nft class details")
	ErrNftCreate           = Register(RootCodeSpace, "12", "failed to create nft class")
	ErrNftGet              = Register(RootCodeSpace, "13", "failed to get nft list")
	ErrNftDetailsGet       = Register(RootCodeSpace, "14", "failed to get nft details")
	ErrNftOptionHistoryGet = Register(RootCodeSpace, "15", "failed to get nft option history")
	ErrNftEdit             = Register(RootCodeSpace, "16", "failed to edit nft")
	ErrNftBatchEdit        = Register(RootCodeSpace, "17", "failed to batch edit nft")
	ErrNftBurn             = Register(RootCodeSpace, "18", "failed to burn nft")
	ErrNftBatchBurn        = Register(RootCodeSpace, "19", "failed to batch burn nft")
	ErrTxResult            = Register(RootCodeSpace, "20", "failed to get tx result")
	ErrIdempotent          = Register(RootCodeSpace, "21", "failed to idempotent")
	ErrNftClassTransfer    = Register(RootCodeSpace, "22", "failed to transfer nft class")
	ErrNftParams           = Register(RootCodeSpace, "23", "failed to get nft class")
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
