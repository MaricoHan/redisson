package service

import (
	"context"
	"fmt"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/modext"
	"strings"
	"time"
)

type NftClass struct {
}

func NewNftClass() *NftClass {
	return &NftClass{}
}

func (svc *NftClass) CreateNftClass(params dto.CreateNftClassP) ([]string, error) {

	// 写入数据库
	// sdk 创建账户
	db, err := orm.GetDB().Begin()
	if err != nil {
		return nil, types.ErrMysqlConn
	}

	tClass := &models.TClass{
		AppID:       params.AppID,
		Name:        null.StringFromPtr(&params.Name),
		Symbol:      null.StringFromPtr(&params.Symbol),
		URI:         null.StringFromPtr(&params.Uri),
		URIHash:     null.StringFromPtr(&params.UriHash),
		Description: null.StringFromPtr(&params.Description),
		Metadata:    null.StringFromPtr(&params.Data),
		Owner:       params.Owner,

		ClassID:   "",
		TXHash:    "",
		Timestamp: null.NewTime(time.Now(), true),
	}

	err = tClass.Insert(context.Background(), db, boil.Infer())
	if err != nil {
		return nil, types.ErrAccountCreate
	}

	return nil, nil
}

func (svc *NftClass) NftClasses(params dto.NftClassesP) (*dto.NftClassesRes, error) {
	result := &dto.NftClassesRes{}
	result.Offset = params.Offset
	result.Limit = params.Limit
	result.NftClasses = []*dto.NftClass{}
	queryMod := []qm.QueryMod{
		qm.From(models.TableNames.TClasses),
		qm.Select(models.TClassColumns.ID, models.TClassColumns.Name, models.TClassColumns.Symbol,
			models.TClassColumns.Offset, models.TClassColumns.URI, models.TClassColumns.Owner,
			models.TClassColumns.TXHash, models.TClassColumns.Timestamp),
		models.TClassWhere.AppID.EQ(params.AppID),
	}
	if params.Id != "" {
		queryMod = append(queryMod, models.TClassWhere.ClassID.EQ(params.Id))
	}
	if params.Name != "" { //支持模糊查询
		queryMod = append(queryMod, models.TClassWhere.Name.NEQ(null.StringFromPtr(&params.Name)))
	}
	if params.Owner != "" {
		queryMod = append(queryMod, models.TClassWhere.Owner.EQ(params.Owner))
	}
	if params.TxHash != "" {
		queryMod = append(queryMod, models.TClassWhere.TXHash.EQ(params.TxHash))
	}

	if params.StartDate != nil {
		queryMod = append(queryMod, models.TClassWhere.CreateAt.GTE(*params.StartDate))
	}
	if params.EndDate != nil {
		queryMod = append(queryMod, models.TClassWhere.CreateAt.LTE(*params.EndDate))
	}
	if params.SortBy != "" {
		orderBy := ""
		switch params.SortBy {
		case "DATE_DESC":
			orderBy = fmt.Sprintf("%s desc", models.TClassColumns.CreateAt)
		case "DATE_ASC":
			orderBy = fmt.Sprintf("%s ASC", models.TClassColumns.CreateAt)
		}
		queryMod = append(queryMod, qm.OrderBy(orderBy))
	}

	var modelResults []*models.TClass
	total, err := modext.PageQuery(
		context.Background(),
		orm.GetDB(),
		queryMod,
		&modelResults,
		params.Offset,
		params.Limit,
	)
	if err != nil {
		// records not exist
		if strings.Contains(err.Error(), "records not exist") {
			return result, nil
		}

		return nil, types.ErrMysqlConn
	}

	result.TotalCount = total

	var classIds []string
	for _, modelResult := range modelResults {
		classIds = append(classIds, modelResult.ClassID)
	}
	q1 := []qm.QueryMod{
		qm.From(models.TableNames.TNFTS),
		qm.Select(models.TNFTColumns.ClassID),
		qm.Select("count(class_id) as count AND class_id"),
		qm.GroupBy(models.TNFTColumns.ClassID),
		//SELECT sex,COUNT(sex) FROM employee GROUP BY sex;
	}
	q1 = append(q1, models.TNFTWhere.ClassID.IN(classIds))
	var countRes []*dto.NftCount
	models.NewQuery(q1...).Bind(context.Background(), orm.GetDB(), &countRes)

	var nftClasses []*dto.NftClass
	for _, modelResult := range modelResults {
		nftClass := &dto.NftClass{
			Id:        modelResult.ClassID,
			Name:      modelResult.Name.String,
			Symbol:    modelResult.Symbol.String,
			NftCount:  uint64(0),
			Uri:       modelResult.URI.String,
			Owner:     modelResult.Owner,
			TxHash:    modelResult.TXHash,
			Timestamp: modelResult.Timestamp.Time.String(),
		}
		for _, r := range countRes {
			if r.ClassId == modelResult.ClassID {
				nftClass.NftCount = uint64(r.Count)
			}
		}
		nftClasses = append(nftClasses, nftClass)
	}
	result.NftClasses = nftClasses
	return result, nil
}

func (svc *NftClass) NftClassById(params dto.NftClassesP) (*dto.NftClassRes, error) {
	if params.Id == "" {
		return nil, types.ErrNftClassDetailsGet
	}

	classOne, err := models.TClasses(
		models.TClassWhere.ClassID.EQ(params.Id),
		models.TClassWhere.AppID.EQ(params.AppID),
	).OneG(context.Background())

	if err != nil {
		return nil, types.ErrTxResult
	}

	count, err := models.TNFTS(
		models.TNFTWhere.ClassID.EQ(params.Id),
		models.TNFTWhere.AppID.EQ(params.AppID),
	).CountG(context.Background())

	result := &dto.NftClassRes{}
	result.Id = classOne.ClassID
	result.Timestamp = classOne.Timestamp.Time.String()
	result.Name = classOne.Name.String
	result.Uri = classOne.URI.String
	result.Owner = classOne.Owner
	result.Symbol = classOne.Symbol.String
	result.Data = classOne.Metadata.String
	result.Description = classOne.Description.String
	result.UriHash = classOne.URIHash.String
	result.NftCount = uint64(count)

	return result, nil

}
