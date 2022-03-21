package wenchangchain_ddc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"database/sql"
	"encoding/hex"

	"github.com/friendsofgo/errors"
	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/modext"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
)

const ddcNftp = "ddcNftp"

type DDCClass struct {
	base *service.Base
}

func NewDDCClass(base *service.Base) *service.NFTClassBase {
	return &service.NFTClassBase{
		Module: service.DDC,
		Service: &DDCClass{
			base: base,
		},
	}
}

func (svc *DDCClass) List(params dto.NftClassesP) (*dto.NftClassesRes, error) {
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
			qm.From(models.TableNames.TDDCClasses),
			models.TDDCClassWhere.ProjectID.EQ(params.ProjectID),
			models.TDDCClassWhere.Status.EQ(models.TNFTSStatusActive),
		}
		if params.Id != "" {
			queryMod = append(queryMod, models.TDDCClassWhere.ClassID.EQ(params.Id))
		}
		if params.Name != "" {
			//Fuzzy query support
			queryMod = append(queryMod, qm.Where("name like ? ", "%"+params.Name+"%"))
		}
		if params.Owner != "" {
			queryMod = append(queryMod, models.TDDCClassWhere.Owner.EQ(params.Owner))
		}
		if params.TxHash != "" {
			queryMod = append(queryMod, models.TDDCClassWhere.TXHash.EQ(params.TxHash))
		}

		if params.StartDate != nil {
			queryMod = append(queryMod, models.TDDCClassWhere.Timestamp.GTE(null.TimeFromPtr(params.StartDate)))
		}
		if params.EndDate != nil {
			queryMod = append(queryMod, models.TDDCClassWhere.Timestamp.LTE(null.TimeFromPtr(params.EndDate)))
		}
		if params.SortBy != "" {
			orderBy := ""
			switch params.SortBy {
			case "DATE_DESC":
				orderBy = fmt.Sprintf("%s DESC", models.TDDCClassColumns.CreateAt)
			case "DATE_ASC":
				orderBy = fmt.Sprintf("%s ASC", models.TDDCClassColumns.CreateAt)
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
			log.Error("nft ddc classes", "query ddc class error:", err.Error())
			return types.ErrInternal
		}
		result.TotalCount = total

		var classIds []string
		for _, modelResult := range modelResults {
			classIds = append(classIds, modelResult.ClassID)
		}
		q1 := []qm.QueryMod{
			qm.From(models.TableNames.TDDCNFTS),
			qm.Select(models.TDDCNFTColumns.ClassID),
			qm.Select("count(class_id) as count, class_id"),
			models.TDDCNFTWhere.Status.EQ(models.TNFTSStatusActive),
			qm.GroupBy(models.TDDCNFTColumns.ClassID),
			models.TDDCNFTWhere.ClassID.IN(classIds),
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

func (svc *DDCClass) Show(params dto.NftClassesP) (*dto.NftClassRes, error) {
	var err error
	var classOne *models.TDDCClass
	var count int64
	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		classOne, err = models.TDDCClasses(
			models.TDDCClassWhere.ClassID.EQ(params.Id),
			models.TDDCClassWhere.ProjectID.EQ(params.ProjectID),
		).One(context.Background(), exec)

		if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
			(err != nil && strings.Contains(err.Error(), service.SqlNotFound)) {
			//404
			return types.ErrNotFound
		} else if err != nil {
			//500
			log.Error("ddc class by id", "query class error:", err.Error())
			return types.ErrInternal
		}

		if !strings.Contains(models.TDDCClassesStatusActive, classOne.Status) {
			return types.ErrNftClassStatus
		}

		count, err = models.TDDCNFTS(
			models.TDDCNFTWhere.ClassID.EQ(params.Id),
			models.TDDCNFTWhere.ProjectID.EQ(params.ProjectID),
			models.TDDCNFTWhere.Status.EQ(models.TNFTSStatusActive),
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

func (svc *DDCClass) Create(params dto.CreateNftClassP) (*dto.TxRes, error) {
	//owner不能为project外的账户
	_, err := models.TDDCAccounts(
		models.TDDCAccountWhere.ProjectID.EQ(params.ProjectID),
		models.TDDCAccountWhere.Address.EQ(params.Owner)).OneG(context.Background())
	if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
		(err != nil && strings.Contains(err.Error(), service.SqlNotFound)) {
		//400
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrOwnerFound)
	} else if err != nil {
		//500
		log.Error("create ddc class", "validate owner error:", err.Error())
		return nil, types.ErrInternal
	}

	//new classId
	var data = []byte(params.Owner)
	data = append(data, []byte(params.Name)...)
	data = append(data, []byte(strconv.FormatInt(time.Now().Unix(), 10))...)
	data = append(data, []byte(fmt.Sprintf("%d", rand.Int()))...)
	classId := ddcNftp + strings.ToLower(hex.EncodeToString(tmhash.Sum(data)))

	createDenomMsg := nft.MsgIssueDenom{
		Id:               classId,
		Name:             params.Name,
		Sender:           params.Owner,
		Symbol:           params.Symbol,
		MintRestricted:   true,
		UpdateRestricted: false,
		Description:      params.Description,
		Uri:              params.Uri,
		UriHash:          params.UriHash,
		Data:             params.Data,
	}
	createDenomMsgByte, err := createDenomMsg.Marshal()
	if err != nil {
		log.Error("create ddc class", "createDenomMsgByte marshal error:", err.Error())
		return nil, types.ErrInternal
	}
	message := []interface{}{createDenomMsg}
	messageBytes, _ := json.Marshal(message)
	code := fmt.Sprintf("%s%s%s", params.Owner, models.TTXSOperationTypeIssueClass, time.Now().String())
	taskId := svc.base.EncodeData(code)
	hash := svc.base.EncodeData(string(createDenomMsgByte))
	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		ttx := models.TDDCTX{
			ProjectID:     params.ProjectID,
			Hash:          hash,
			Sender:        null.StringFrom(params.Owner),
			Timestamp:     null.TimeFrom(time.Now()),
			OriginData:    null.BytesFromPtr(&createDenomMsgByte),
			Message:       null.JSONFrom(messageBytes),
			TaskID:        null.StringFrom(taskId),
			GasUsed:       null.Int64From(0),
			OperationType: models.TTXSOperationTypeIssueClass,
			Status:        models.TTXSStatusSuccess,
			Tag:           null.JSONFrom(params.Tag),
			Retry:         null.Int8From(0),
			BizFee:        null.Int64From(0),
		}
		err = ttx.Insert(context.Background(), exec, boil.Infer())
		if err != nil {
			log.Error("create ddc class", "tx insert error:", err.Error())
			return types.ErrInternal
		}
		msgs := models.TDDCMSG{
			ProjectID: ttx.ProjectID,
			TXHash:    hash,
			Module:    models.TDDCMSGSModuleNFT,
			Operation: models.TDDCMSGSOperationIssueClass,
			Signer:    params.Owner,
			Timestamp: null.TimeFrom(time.Now()),
			Message:   messageBytes,
		}
		err = msgs.Insert(context.Background(), exec, boil.Infer())
		if err != nil {
			log.Error("create ddc class", "msg insert error:", err.Error())
			return types.ErrInternal
		}
		ddcClass := models.TDDCClass{
			ProjectID:   params.ProjectID,
			TXHash:      hash,
			ClassID:     createDenomMsg.Id,
			Name:        null.StringFrom(createDenomMsg.Name),
			Symbol:      null.StringFrom(createDenomMsg.Symbol),
			URI:         null.StringFrom(createDenomMsg.Uri),
			URIHash:     null.StringFrom(createDenomMsg.UriHash),
			Description: null.StringFrom(createDenomMsg.Description),
			Owner:       createDenomMsg.Sender,
			Status:      models.TClassesStatusActive,
			Timestamp:   null.TimeFrom(time.Now()),
			Metadata:    null.StringFrom(createDenomMsg.Data),
		}
		err = ddcClass.Insert(context.Background(), exec, boil.Infer())
		if err != nil {
			log.Error("create ddc class", "class insert error: ", err.Error())
			return types.ErrInternal
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &dto.TxRes{TaskId: taskId}, nil
}
