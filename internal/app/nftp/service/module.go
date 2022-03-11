package service

import "gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"

const (
	NATIVE = "wenchangchain-native"
	DDC    = "wenchangchain-ddc"
)

// AccountService accounts
type AccountService interface {
	Create(params dto.CreateAccountP) (*dto.AccountRes, error)
	Show(params dto.AccountsP) (*dto.AccountsRes, error)
	History(params dto.AccountsP) (*dto.AccountOperationRecordRes, error) // 操作记录
}

// AccountBase accounts
type AccountBase struct {
	Module  string
	Service AccountService
}

// NFTClassService NFTClass
type NFTClassService interface {
	List(params dto.NftClassesP) (*dto.NftClassesRes, error) // 列表
	Show(params dto.NftClassesP) (*dto.NftClassRes, error)   // 详情
	Create(params dto.CreateNftClassP) (*dto.TxRes, error)   // 创建
}

// NFTClassBase NFTClass
type NFTClassBase struct {
	Module  string
	Service NFTClassService
}

// NFTService NFT
type NFTService interface {
	List(params dto.NftsP) (*dto.NftsRes, error)
	Create(params dto.CreateNftsP) (*dto.TxRes, error)
	Show(params dto.NftByNftIdP) (*dto.NftR, error)
	Update(params dto.EditNftByNftIdP) (*dto.TxRes, error)
	Delete(params dto.DeleteNftByNftIdP) (*dto.TxRes, error)
	History(params dto.NftOperationHistoryByNftIdP) (*dto.BNftOperationHistoryByNftIdRes, error)
}

// NFTBase NFT
type NFTBase struct {
	Module  string
	Service NFTService
}

// TransferService Transfer
type TransferService interface {
	TransferNFTClass(params dto.TransferNftClassByIDP) (*dto.TxRes, error) // 转让NFTClass
	TransferNFT(params dto.TransferNftByNftIdP) (*dto.TxRes, error)        // 转让NFT
}

// TransferBase Transfer
type TransferBase struct {
	Module  string
	Service TransferService
}
