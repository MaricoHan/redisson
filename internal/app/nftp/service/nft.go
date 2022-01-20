package service

import (
	"context"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
	"strconv"
)

type Nft struct {
}

func NewNft() *Nft {
	return &Nft{}
}
func (svc *Nft) EditNftByIndex(params dto.EditNftByIndexP) (int64, error) {
	// get database object
	db, err := orm.GetDB().Begin()
	if err != nil {
		return 0, types.ErrMysqlConn
	}

	//get NFT by app_id,class_id and index
	tNft, err := models.TNFTS(qm.Where("app_id=? AND class_id=? AND index=?", params.AppID, params.ClassId, params.Index)).One(context.Background(), db)

	if err != nil {
		return 0, types.ErrInternal
	}

	//!!!not sure the error msg
	if tNft == nil || tNft.Status == "burned" {
		return 0, types.ErrNftDetailsGet
	}

	//assignment
	tNft.URI = null.StringFrom(params.Uri)
	tNft.Metadata = null.StringFrom(params.Data)

	//update database
	rowsAff, err := tNft.Update(context.Background(), db, boil.Infer())

	//return the affected rows amount
	return rowsAff, err

}

func (svc *Nft) EditNftByBatch(params dto.EditNftByBatchP) (int64, error) {

	rowsAff := int64((0))
	return rowsAff, nil
}

func (svc *Nft) DeleteNftByIndex(params dto.DeleteNftByIndexP) (int64, error) {

	// get database object
	db, err := orm.GetDB().Begin()
	if err != nil {
		return 0, types.ErrMysqlConn
	}

	//get NFT by app_id,class_id and index
	tNft, err := models.TNFTS(qm.Where("app_id=? AND class_id=? AND index=?", params.AppID, params.ClassId, params.Index)).One(context.Background(), db)
	//!!!not sure the error msg
	if tNft == nil || tNft.Status == "burned" {
		return 0, types.ErrNftDetailsGet
	}
	if tNft.Status == "pendding" {
		return 0, types.ErrNftDetailsGet
	}
	//just burn
	tNft.Status = "burned"
	rowsAff, err := tNft.Update(context.Background(), db, boil.Infer())

	//return the affected rows amount
	return rowsAff, nil
}
func (svc *Nft) DeleteNftByBatch(params dto.DeleteNftByBatchP) (int64, error) {

	rowsAff := int64((0))
	return rowsAff, nil
}

func (svc *Nft) NftByIndex(params dto.NftByIndexP) (*dto.NftByIndexP, error) {

	// get database object
	db, err := orm.GetDB().Begin()
	if err != nil {
		return nil, types.ErrMysqlConn
	}
	//get NFT by app_id,class_id and index
	tNft, err := models.TNFTS(qm.Where("app_id=? AND class_id=? AND index=?", params.AppID, params.ClassId, params.Index)).One(context.Background(), db)
	//get class by class_id
	class, err := models.TClasses(qm.Where("class_id=?", params.ClassId)).One(context.Background(), db)
	result := &dto.NftByIndexP{
		Id:          strconv.FormatInt(int64(tNft.ID), 10),
		Index:       tNft.Index,
		Name:        tNft.Name.String,
		ClassId:     tNft.ClassID,
		ClassName:   class.Name.String,
		ClassSymbol: class.Symbol.String,
		Uri:         tNft.URI.String,
		UriHash:     tNft.URIHash.String,
		Data:        tNft.Metadata.String,
		Owner:       tNft.Owner,
		Status:      tNft.Status,
		TxHash:      tNft.TXHash,
		TimeStamp:   tNft.Timestamp.Time.String(),
	}

	return result, nil
}
