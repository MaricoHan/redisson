package service

import "gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"

const (
	NATIVE = "native"
	DDC    = "ddc721"
)

// AccountService accounts
type AccountService interface {
	Create(dto.CreateAccountP) (*dto.AccountRes, error)
	Show(dto.AccountsP) (*dto.AccountsRes, error)
	History(dto.AccountsP) (*dto.AccountOperationRecordRes, error) // 操作记录
}

// AccountBase accounts
type AccountBase struct {
	Module  string
	Service AccountService
}

// NFTClassService NFTClass
type NFTClassService interface {
	List(dto.NftClassesP) (*dto.NftClassesRes, error) // 列表
	Show(dto.NftClassesP) (*dto.NftClassRes, error)   // 详情
	Create(dto.CreateNftClassP) (*dto.TxRes, error)   // 创建
}

// NFTClassBase NFTClass
type NFTClassBase struct {
	Module  string
	Service NFTClassService
}

// NFTService NFT
type NFTService interface {
	List(dto.NftsP) (*dto.NftsRes, error)
	Create(dto.CreateNftsP) (*dto.TxRes, error)
	Show(dto.NftByNftIdP) (*dto.NftR, error)
	Update(dto.EditNftByNftIdP) (*dto.TxRes, error)
	Delete(dto.DeleteNftByNftIdP) (*dto.TxRes, error)
	History(dto.NftOperationHistoryByNftIdP) (*dto.BNftOperationHistoryByNftIdRes, error)
}

// NFTBase NFT
type NFTBase struct {
	Module  string
	Service NFTService
}

// TransferService Transfer
type TransferService interface {
	TransferNFTClass(dto.TransferNftClassByIDP) (*dto.TxRes, error) // 转让NFTClass
	TransferNFT(dto.TransferNftByNftIdP) (*dto.TxRes, error)        // 转让NFT
}

// TransferBase Transfer
type TransferBase struct {
	Module  string
	Service TransferService
}

// TXService TX
type TXService interface {
	Show(dto.TxResultByTxHashP) (*dto.TxResultByTxHashRes, error)
}

// TXBase TX
type TXBase struct {
	Module  string
	Service TXService
}
