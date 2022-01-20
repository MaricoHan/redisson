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

	if tNft == nil || tNft.Status == "burned" {
		return 0, types.ErrNftMissing
	}

	//assign
	tNft.URI = null.StringFrom(params.Uri)
	tNft.Metadata = null.StringFrom(params.Data)

	//update database
	rowsAff, err := tNft.Update(context.Background(), db, boil.Infer())

	//return the affected rows amount
	return rowsAff, err

}

func (svc *Nft) EditNftByBatch(params dto.EditNftByBatchP) (int64, error) {

	// get database object
	db, err := orm.GetDB().Begin()
	if err != nil {
		return 0, types.ErrMysqlConn
	}
	var rowsAff int64
	for _, EditNft := range params.EditNfts { //edit every NFT
		//get NFT by app_id,class_id and index
		tNft, err := models.TNFTS(qm.Where("app_id=? AND class_id=? AND index=?", params.AppID, params.ClassId, EditNft.Index)).One(context.Background(), db)

		if err != nil {
			return rowsAff, types.ErrInternal
		}

		if tNft == nil || tNft.Status == "burned" {
			return rowsAff, types.ErrNftMissing
		}

		//assign
		tNft.URI = null.StringFrom(EditNft.Uri)
		tNft.Metadata = null.StringFrom(EditNft.Data)

		//update database
		i, err := tNft.Update(context.Background(), db, boil.Infer())
		rowsAff += i

		if err != nil {
			return rowsAff, err
		}
	}
	return rowsAff, err
}

func (svc *Nft) DeleteNftByIndex(params dto.DeleteNftByIndexP) (int64, error) {

	// get database object
	db, err := orm.GetDB().Begin()
	if err != nil {
		return 0, types.ErrMysqlConn
	}

	//get NFT by app_id,class_id and index
	tNft, err := models.TNFTS(qm.Where("app_id=? AND class_id=? AND index=?", params.AppID, params.ClassId, params.Index)).One(context.Background(), db)
	if err != nil {
		return 0, err
	}
	if tNft == nil || tNft.Status == "burned" {
		return 0, types.ErrNftMissing
	}
	if tNft.Status == "pendding" {
		return 0, types.ErrNftBurnPend
	}
	//just burn 🔥
	tNft.Status = "burned"
	rowsAff, err := tNft.Update(context.Background(), db, boil.Infer())

	//return the affected rows amount
	return rowsAff, nil
}
func (svc *Nft) DeleteNftByBatch(params dto.DeleteNftByBatchP) (int64, error) {

	// get database object
	db, err := orm.GetDB().Begin()
	if err != nil {
		return 0, types.ErrMysqlConn
	}
	var rowsAff int64
	for _, index := range params.Indices { //burn every NFT
		//get NFT by app_id,class_id and index
		tNft, err := models.TNFTS(qm.Where("app_id=? AND class_id=? AND index=?", params.AppID, params.ClassId, index)).One(context.Background(), db)
		if err != nil {
			return rowsAff, err
		}
		if tNft == nil || tNft.Status == "burned" {
			return rowsAff, types.ErrNftMissing
		}
		if tNft.Status == "pendding" {
			return rowsAff, types.ErrNftBurnPend
		}

		//just burn
		tNft.Status = "burned"
		i, err := tNft.Update(context.Background(), db, boil.Infer())
		rowsAff += i

		if err != nil {
			return rowsAff, err
		}
	}

	//return the affected rows amount
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
