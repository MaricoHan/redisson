package service

import (
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
)

type NftTransfer struct {
}

func NewNftTransfer() *NftTransfer {
	return &NftTransfer{}
}

func (svc *NftTransfer) TransferNftClassByID(params dto.TransferNftClassByIDP) (string, error) {
	//db, err := orm.GetDB().Begin()
	//if err != nil {
	//	return "", types.ErrMysqlConn
	//}

	return "", nil
}

func (svc *NftTransfer) TransferNftByIndex(params dto.TransferNftByIndexP) (*dto.AccountsRes, error) {
	return nil, nil
}

func (svc *NftTransfer) TransferNftByBatch(params dto.TransferNftByBatchP) (*dto.AccountsRes, error) {
	panic("not yet implemented")
}
