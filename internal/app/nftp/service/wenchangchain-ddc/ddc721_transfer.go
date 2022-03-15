package wenchangchain_ddc

import (
	"context"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
)

type DDC721Transfer struct {
	Base
}

func NEWDDC721Transfer(base *service.Base) *service.TransferBase {
	return &service.TransferBase{
		Module: service.DDC,
		Service: &DDC721Transfer{
			NewBase(base),
		},
	}
}
func (svc DDC721Transfer) TransferNFTClass(params dto.TransferNftClassByIDP) (*dto.TxRes, error) {
	panic("...")
}
func (svc DDC721Transfer) TransferNFT(params dto.TransferNftByNftIdP) (*dto.TxRes, error) {
	// ValidateSigner
	if err := svc.base.ValidateSigner(params.Owner, params.ProjectID); err != nil {
		return nil, err
	}

	// ValidateRecipient
	if err := svc.base.ValidateRecipient(params.Recipient, params.ProjectID); err != nil {
		return nil, err
	}

	//查出ddc
	res, err := models.TNFTS(
		models.TNFTWhere.NFTID.EQ(params.NftId),
		models.TNFTWhere.ClassID.EQ(params.ClassID),
		models.TNFTWhere.ProjectID.EQ(params.ProjectID),
		models.TNFTWhere.Owner.EQ(params.Owner),
	).OneG(context.Background())

}
