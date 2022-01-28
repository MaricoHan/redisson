package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/modext"
)

type NftClass struct {
	base *Base
}

func NewNftClass(base *Base) *NftClass {
	return &NftClass{base: base}
}

func (svc *NftClass) CreateNftClass(params dto.CreateNftClassP) ([]string, error) {
	//platform address
	classOne, err := models.TAccounts(
		models.TAccountWhere.AppID.EQ(uint64(0)),
	).OneG(context.Background())
	if err != nil {
		return nil, types.ErrGetAccountDetails
	}
	pAddress := classOne.Address
	//new classId
	var data = []byte(params.Owner)
	data = append(data, []byte(params.Name)...)
	data = append(data, []byte(strconv.FormatInt(time.Now().Unix(), 10))...)
	classId := nftp + strings.ToLower(hex.EncodeToString(tmhash.Sum(data)))
	//txMsg, Platform side created
	baseTx := svc.base.CreateBaseTx(pAddress, defultKeyPassword)
	createDenomMsg := nft.MsgIssueDenom{
		Id:               classId,
		Name:             params.Name,
		Sender:           pAddress,
		Symbol:           params.Symbol,
		MintRestricted:   true,
		UpdateRestricted: false,
		Description:      params.Description,
		Uri:              params.Uri,
		UriHash:          params.UriHash,
		Data:             params.Data,
	}
	transferDenomMsg := nft.MsgTransferDenom{
		Id:        classId,
		Sender:    pAddress,
		Recipient: params.Owner,
	}

	originData, txHash, err := svc.base.BuildAndSign(sdktype.Msgs{&createDenomMsg, &transferDenomMsg}, baseTx)
	if err != nil {
		log.Debug("create nft class", "BuildAndSign error:", err.Error())
		return nil, err
	}

	//validate tx
	txone, err := svc.base.ValidateTx(txHash)
	if err != nil {
		return nil, err
	}
	if txone != nil && txone.Status == models.TTXSStatusFailed {
		baseTx.Memo = string(txone.ID)
		originData, txHash, err = svc.base.BuildAndSign(sdktype.Msgs{&createDenomMsg, &transferDenomMsg}, baseTx)
		if err != nil {
			log.Debug("create nft class", "BuildAndSign error:", err.Error())
			return nil, types.ErrBuildAndSign
		}
	}

	//transferInfo
	ttx := models.TTX{
		AppID:         params.AppID,
		Hash:          txHash,
		Timestamp:     null.Time{Time: time.Now()},
		OriginData:    null.BytesFromPtr(&originData),
		OperationType: models.TTXSOperationTypeIssueClass,
		Status:        models.TTXSStatusUndo,
	}
	err = ttx.InsertG(context.Background(), boil.Infer())
	if err != nil {
		return nil, err
	}
	var hashs []string
	hashs = append(hashs, txHash)
	return hashs, nil
}

func (svc *NftClass) NftClasses(params dto.NftClassesP) (*dto.NftClassesRes, error) {
	result := &dto.NftClassesRes{}
	result.Offset = params.Offset
	result.Limit = params.Limit
	result.NftClasses = []*dto.NftClass{}
	var modelResults []*models.TClass
	var countRes []*dto.NftCount
	err := modext.Transaction(func(exec boil.ContextExecutor) error {
		queryMod := []qm.QueryMod{
			qm.From(models.TableNames.TClasses),
			models.TClassWhere.AppID.EQ(params.AppID),
			models.TClassWhere.Status.EQ(models.TNFTSStatusActive),
		}
		if params.Id != "" {
			queryMod = append(queryMod, models.TClassWhere.ClassID.EQ(params.Id))
		}
		if params.Name != "" {
			//Fuzzy query support
			queryMod = append(queryMod, qm.Where("name like ? ", "%"+params.Name+"%"))
		}
		if params.Owner != "" {
			queryMod = append(queryMod, models.TClassWhere.Owner.EQ(params.Owner))
		}
		if params.TxHash != "" {
			queryMod = append(queryMod, models.TClassWhere.TXHash.EQ(params.TxHash))
		}

		if params.StartDate != nil {
			queryMod = append(queryMod, models.TClassWhere.Timestamp.GTE(null.TimeFromPtr(params.StartDate)))
		}
		if params.EndDate != nil {
			queryMod = append(queryMod, models.TClassWhere.Timestamp.LTE(null.TimeFromPtr(params.EndDate)))
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

		total, err := modext.PageQueryByOffset(
			context.Background(),
			exec,
			queryMod,
			&modelResults,
			int(params.Offset),
			int(params.Limit),
		)
		if err != nil {
			return err
		}
		result.TotalCount = total

		var classIds []string
		for _, modelResult := range modelResults {
			classIds = append(classIds, modelResult.ClassID)
		}
		q1 := []qm.QueryMod{
			qm.From(models.TableNames.TNFTS),
			qm.Select(models.TNFTColumns.ClassID),
			qm.Select("count(class_id) as count, class_id"),
			models.TNFTWhere.Status.EQ(models.TNFTSStatusActive),
			qm.GroupBy(models.TNFTColumns.ClassID),
			models.TNFTWhere.ClassID.IN(classIds),
		}

		err = models.NewQuery(q1...).Bind(context.Background(), exec, &countRes)
		if err != nil {
			return types.ErrNftCountByClass
		}
		return err
	})
	if err != nil {
		if strings.Contains(err.Error(), "records not exist") {
			return result, nil
		}
		return result, err
	}

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
	var err error
	var classOne *models.TClass
	var count int64
	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		classOne, err = models.TClasses(
			models.TClassWhere.ClassID.EQ(params.Id),
			models.TClassWhere.AppID.EQ(params.AppID),
		).One(context.Background(), exec)
		if err != nil {
			return types.ErrNftClassesGet
		}

		count, err = models.TNFTS(
			models.TNFTWhere.ClassID.EQ(params.Id),
			models.TNFTWhere.AppID.EQ(params.AppID),
		).Count(context.Background(), exec)
		if err != nil {
			return types.ErrNftClassDetailsGet
		}
		return err
	})
	if err != nil {
		return nil, types.ErrNotFound
	}

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
