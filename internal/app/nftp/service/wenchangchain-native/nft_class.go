package wenchangchain_native

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gitlab.bianjie.ai/irita-paas/open-api/config"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/friendsofgo/errors"
	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/modext"
)

type NFTClass struct {
	Base
}

func NewNFTClass(base *service.Base) *service.NFTClassBase {
	return &service.NFTClassBase{
		Module: service.NATIVE,
		Service: &NFTClass{
			NewBase(base),
		},
	}
}

func (svc *NFTClass) List(params dto.NftClassesP) (*dto.NftClassesRes, error) {
	result := &dto.NftClassesRes{
		PageRes: dto.PageRes{
			Offset:     params.Offset,
			Limit:      params.Limit,
			TotalCount: 0,
		},
		NftClasses: []*dto.NftClass{},
	}
	var modelResults []*models.TClass
	var countRes []*dto.NftCount
	err := modext.Transaction(func(exec boil.ContextExecutor) error {
		queryMod := []qm.QueryMod{
			qm.From(models.TableNames.TClasses),
			models.TClassWhere.ProjectID.EQ(params.ProjectID),
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
				orderBy = fmt.Sprintf("%s DESC", models.TClassColumns.CreateAt)
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
			// records not exist
			if strings.Contains(err.Error(), service.SqlNotFound) {
				return nil
			}
			log.Error("nft classes", "query nft class error:", err.Error())
			return types.ErrInternal
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
			return types.ErrInternal
		}
		return err
	})
	if err != nil {
		if strings.Contains(err.Error(), service.SqlNotFound) {
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
	if nftClasses != nil {
		result.NftClasses = nftClasses
	}
	return result, nil
}

func (svc *NFTClass) Show(params dto.NftClassesP) (*dto.NftClassRes, error) {
	var err error
	var classOne *models.TClass
	var count int64
	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		classOne, err = models.TClasses(
			models.TClassWhere.ClassID.EQ(params.Id),
			models.TClassWhere.ProjectID.EQ(params.ProjectID),
		).One(context.Background(), exec)
		if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
			(err != nil && strings.Contains(err.Error(), service.SqlNotFound)) {
			//404
			return types.ErrNotFound
		} else if err != nil {
			//500
			log.Error("nft class by id", "query class error:", err.Error())
			return types.ErrInternal
		}

		if !strings.Contains(models.TClassesStatusActive, classOne.Status) {
			return types.ErrNftClassStatus
		}

		count, err = models.TNFTS(
			models.TNFTWhere.ClassID.EQ(params.Id),
			models.TNFTWhere.ProjectID.EQ(params.ProjectID),
			models.TNFTWhere.Status.EQ(models.TNFTSStatusActive),
		).Count(context.Background(), exec)
		if err != nil {
			return types.ErrInternal
		}
		return err
	})

	if err != nil {
		return nil, err
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
	result.TxHash = classOne.TXHash

	return result, nil

}

func (svc *NFTClass) Create(params dto.CreateNftClassP) (*dto.TxRes, error) {
	//owner不能为project外的账户
	_, err := models.TAccounts(
		models.TAccountWhere.ProjectID.EQ(params.ProjectID),
		models.TAccountWhere.Address.EQ(params.Owner)).OneG(context.Background())
	if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
		(err != nil && strings.Contains(err.Error(), service.SqlNotFound)) {
		//400
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrOwnerFound)
	} else if err != nil {
		//500
		log.Error("create nft class", "validate owner error:", err.Error())
		return nil, types.ErrInternal
	}

	_, err = models.TAccounts(
		models.TAccountWhere.ProjectID.EQ(params.ProjectID),
		models.TAccountWhere.Address.EQ(params.Owner)).OneG(context.Background())

	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows || strings.Contains(err.Error(), service.SqlNotFound) {
			//400
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrOwnerFound)
		}
		//500
		log.Error("create nft class", "query owner error:", err.Error())
		return nil, types.ErrInternal
	}

	//platform address
	classOne, err := models.TAccounts(
		models.TAccountWhere.ProjectID.EQ(uint64(0)),
	).OneG(context.Background())
	if err != nil {
		log.Error("create nft class", "query platform error:", err.Error())
		return nil, types.ErrInternal
	}
	pAddress := classOne.Address

	//new classId
	var data = []byte(params.Owner)
	data = append(data, []byte(params.Name)...)
	data = append(data, []byte(strconv.FormatInt(time.Now().Unix(), 10))...)
	data = append(data, []byte(fmt.Sprintf("%d", rand.Int()))...)
	classId := nftp + strings.ToLower(hex.EncodeToString(tmhash.Sum(data)))

	//txMsg, Platform side created
	baseTx := svc.base.CreateBaseTx(pAddress, config.Get().Server.DefaultKeyPassword)
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
	originData, txHash, _ := svc.base.BuildAndSign(sdktype.Msgs{&createDenomMsg, &transferDenomMsg}, baseTx)
	baseTx.Gas = svc.base.CreateDenomGas(originData)
	err = svc.base.GasThan(params.ChainID, baseTx.Gas, params.PlatFormID)
	if err != nil {
		return nil, types.NewAppError(types.RootCodeSpace, types.ErrGasNotEnough, err.Error())
	}
	originData, txHash, err = svc.base.BuildAndSign(sdktype.Msgs{&createDenomMsg, &transferDenomMsg}, baseTx)
	if err != nil {
		log.Debug("create nft class", "BuildAndSign error:", err.Error())
		return nil, types.ErrBuildAndSign
	}

	//validate tx
	txone, err := svc.base.ValidateTx(txHash)
	if err != nil {
		return nil, err
	}
	if txone != nil && txone.Status == models.TTXSStatusFailed {
		baseTx.Memo = fmt.Sprintf("%d", txone.ID)
		originData, txHash, err = svc.base.BuildAndSign(sdktype.Msgs{&createDenomMsg, &transferDenomMsg}, baseTx)
		if err != nil {
			log.Debug("create nft class", "BuildAndSign error:", err.Error())
			return nil, types.ErrBuildAndSign
		}
	}

	message := []interface{}{createDenomMsg, transferDenomMsg}
	messageBytes, _ := json.Marshal(message)
	code := fmt.Sprintf("%s%s%s", params.Owner, models.TTXSOperationTypeIssueClass, time.Now().String())
	taskId := svc.base.EncodeData(code)
	ttx := models.TTX{
		ProjectID:     params.ProjectID,
		Hash:          txHash,
		Sender:        null.StringFrom(params.Owner),
		Timestamp:     null.Time{Time: time.Now()},
		OriginData:    null.BytesFromPtr(&originData),
		Message:       null.JSONFrom(messageBytes),
		TaskID:        null.StringFrom(taskId),
		GasUsed:       null.Int64From(int64(baseTx.Gas)),
		OperationType: models.TTXSOperationTypeIssueClass,
		Status:        models.TTXSStatusUndo,
		Tag:           null.JSONFrom(params.Tag),
		Retry:         null.Int8From(0),
	}
	err = ttx.InsertG(context.Background(), boil.Infer())
	if err != nil {
		log.Error("create nft class", "tx insert error:", err.Error())
		return nil, types.ErrInternal
	}
	return &dto.TxRes{TaskId: taskId}, nil
}
