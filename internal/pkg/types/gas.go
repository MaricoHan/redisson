package types

const (
	//mint nft
	MintMinNFTDataSize         = 443
	MintMinNFTIncreaseDataSize = 266
	MintMinNFTGas              = 71388
	MintMinNFTIncreaseGas      = 17951
	MintNFTCoefficient         = 40

	//create denom
	CreateMinDENOMDataSize = 527
	CreateMinDENOMGas      = 75821
	CreateDENOMCoefficient = 73

	//transfer denom
	TransferMinDENOMDataSize = 274
	TransferMinDENOMGas      = 63404
	TransferDENOMCoefficient = 33.3

	//transfer nft
	TransferMinNFTDataSize         = 446
	TransferMinNFTIncreaseDataSize = 269
	TransferMinNFTGas              = 69636
	TransferMinNFTIncreaseGas      = 16077
	TransferNFTCoefficient         = 43

	//account
	CreateAccountGas = 80000

	// edit nft
	EditNFTBaseGas            = 48130
	EditNFTLenCoefficient     = 7
	EditNFTSignLenCoefficient = 42

	// edit batch nft
	EditBatchNFTBaseGas            = 65420
	EditBatchNFTLenCoefficient     = 8
	EditBatchNFTSignLenCoefficient = 42

	//delete nft
	DeleteNFTBaseLen     = 371
	DeleteNFTBaseGas     = 62019
	DeleteNFTCoefficient = 4

	//delete batch nft
	DeleteBatchNFTBaseLen            = 373
	DeleteBatchNFTBaseLenCoefficient = 373
	DeleteBatchNFTBaseGas            = 63105
	DeleteBatchNFTBaseGasCoefficient = 9612
	DeleteBatchCoefficient           = 3
)
