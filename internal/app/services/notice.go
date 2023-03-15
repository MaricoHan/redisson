package services

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	httptransport "github.com/go-kit/kit/transport/http"
	log "github.com/sirupsen/logrus"
	pb_notice "gitlab.bianjie.ai/avata/chains/api/pb/v2/notice"
	authErr "gitlab.bianjie.ai/avata/utils/errors/auth"

	noticeResp "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/notice"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	notice2 "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/notice"
	"gitlab.bianjie.ai/avata/open-api/internal/app/repository/db/project"
	"gitlab.bianjie.ai/avata/open-api/internal/app/repository/db/service_redirect_url"
	"gitlab.bianjie.ai/avata/open-api/internal/app/repository/db/user"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/configs"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	"gitlab.bianjie.ai/avata/open-api/utils"
	"gitlab.bianjie.ai/avata/utils/commons/aes"
	"gitlab.bianjie.ai/avata/utils/errors"
)

type INotice interface {
	TransferNFTS(ctx context.Context, params *notice2.TransferNFTS) (*noticeResp.TransferNFTS, error)
	TransferClasses(ctx context.Context, params *notice2.TransferClasses) (*noticeResp.TransferClasses, error)
}

type notice struct {
	logger *log.Logger
}

func NewNotice(logger *log.Logger) *notice {
	return &notice{logger: logger}
}

// TransferNFTS 转让NFT通知
func (a *notice) TransferNFTS(ctx context.Context, params *notice2.TransferNFTS) (*noticeResp.TransferNFTS, error) {
	logger := a.logger.WithField("params", params).WithField("func", "transfer nft")
	path := ctx.Value(httptransport.ContextKeyRequestPath).(string)[len(configs.Cfg.App.RouterPrefix)+1:]
	res := &noticeResp.TransferNFTS{}
	project, err := a.getProject(params.ProjectID)
	if err != nil {
		logger.WithError(err).Error("query project")
		return res, err
	}
	authData, err := utils.AuthData(ctx)
	if err != nil {
		logger.WithError(err).Error("auth data")
		return res, errors.ErrInternal
	}
	user, err := a.getUser(authData.UserId)
	if err != nil {
		logger.WithError(err).Error("query project")
		return res, err
	}
	url, err := a.getServiceRedirectUrl(uint64(project.Id))
	if err != nil {
		logger.WithError(err).Error("query service redirect url")
		return res, err
	}
	mapKey := fmt.Sprintf("%s-%s", authData.Code, authData.Module)
	grpcClient, ok := initialize.NoticeClientMap[mapKey]
	if !ok {
		return nil, errors.New(errors.InternalError, errors.ErrService)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	//  此处 err 由微服务返回, 已是定义好的异常，可直接返回给到用户
	resp, err := grpcClient.Nft(ctx, &pb_notice.NoticeRequest{TxHash: params.TxHash})
	if err != nil {
		logger.WithError(err).Error("grpc request")
		return nil, err
	}
	if resp == nil {
		return nil, errors.New(errors.InternalError, errors.ErrGrpc)
	}

	request := map[string]interface{}{
		"tx_hash":      params.TxHash,
		"project_id":   params.ProjectID,
		"sender":       resp.Sender,
		"recipient":    resp.Recipient,
		"class_id":     resp.ClassId,
		"nft_id":       resp.NftId,
		"block_height": resp.BlockHeight,
		"timestamp":    resp.Timestamp,
	}
	// 组合签名
	timestamp := utils.TimeToUnix(time.Now())
	hash, err := a.hash(request, path, timestamp, &project)
	if err != nil {
		logger.WithError(err).Error("hash")
		return res, errors.ErrInternal
	}
	_, err = a.request(context.Background(), fmt.Sprintf("%s%s", url.Url, path), project.ApiKey, hash, user.Code, timestamp, request)
	if err != nil {
		return res, err
	}
	return res, nil
}

// TransferClasses 转让Class通知
func (a *notice) TransferClasses(ctx context.Context, params *notice2.TransferClasses) (*noticeResp.TransferClasses, error) {
	logger := a.logger.WithField("params", params).WithField("func", "transfer class")
	path := ctx.Value(httptransport.ContextKeyRequestPath).(string)[len(configs.Cfg.App.RouterPrefix)+1:]
	res := &noticeResp.TransferClasses{}
	project, err := a.getProject(params.ProjectID)
	if err != nil {
		logger.WithError(err).Error("query project")
		return res, err
	}
	authData, err := utils.AuthData(ctx)
	if err != nil {
		logger.WithError(err).Error("auth data")
		return res, errors.ErrInternal
	}
	user, err := a.getUser(authData.UserId)
	if err != nil {
		logger.WithError(err).Error("query project")
		return res, err
	}
	url, err := a.getServiceRedirectUrl(uint64(project.Id))
	if err != nil {
		logger.WithError(err).Error("query service redirect url")
		return res, err
	}
	mapKey := fmt.Sprintf("%s-%s", authData.Code, authData.Module)
	grpcClient, ok := initialize.NoticeClientMap[mapKey]
	if !ok {
		return nil, errors.New(errors.InternalError, errors.ErrService)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	//  此处 err 由微服务返回, 已是定义好的异常，可直接返回给到用户
	resp, err := grpcClient.Class(ctx, &pb_notice.NoticeRequest{TxHash: params.TxHash})
	if err != nil {
		logger.WithError(err).Error("grpc request")
		return nil, err
	}
	if resp == nil {
		return nil, errors.New(errors.InternalError, errors.ErrGrpc)
	}
	request := map[string]interface{}{
		"tx_hash":      params.TxHash,
		"project_id":   params.ProjectID,
		"sender":       resp.Sender,
		"recipient":    resp.Recipient,
		"class_id":     resp.ClassId,
		"block_height": resp.BlockHeight,
		"timestamp":    resp.Timestamp,
	}
	timestamp := utils.TimeToUnix(time.Now())
	hash, err := a.hash(request, path, timestamp, &project)
	if err != nil {
		logger.WithError(err).Error("hash")
		return res, errors.ErrInternal
	}
	_, err = a.request(context.Background(), fmt.Sprintf("%s%s", url.Url, path), project.ApiKey, hash, user.Code, timestamp, request)
	if err != nil {
		return res, err
	}
	return res, nil
}

// getProject 获取项目信息
func (a *notice) getProject(projectCode string) (entity.Project, error) {
	projectRepo := project.NewProjectRepo(initialize.MysqlDB)
	project, err := projectRepo.GetProjectByCode(projectCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return project, errors.New(errors.NotFound, authErr.ErrProjectOrUserNotFound)
		}
		return project, err
	}
	return project, nil
}

// getUser 获取用户信息, id&&code
func (a *notice) getUser(userID uint64) (entity.User, error) {
	userRepo := user.NewUserRepo(initialize.MysqlDB)
	user, err := userRepo.GetUser(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New(errors.NotFound, authErr.ErrProjectOrUserNotFound)
		}
		return user, err
	}
	return user, nil
}

// getServiceRedirectUrl 获取服务回调地址
func (a *notice) getServiceRedirectUrl(projectID uint64) (entity.ServiceRedirectUrl, error) {
	serviceRedirectUrlRepo := service_redirect_url.NewServiceRedirectUrlRepo(initialize.MysqlDB)
	sru, err := serviceRedirectUrlRepo.GetServiceRedirectUrlByProjectID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sru, errors.New(errors.NotFound, authErr.ErrServiceRedirectUrlNotFound)
		}
		return sru, err
	}
	if sru.Url == "" {
		return sru, errors.New(errors.NotFound, authErr.ErrServiceRedirectUrlNotFound)
	}
	return sru, nil
}

// request 服务请求
// 1.向上游服务方发起请求, err不等于nil说明存在异常返回服务错误
// 2.当请求成功，返回的状态码为404时, 则返回NOT_FOUNT
func (a *notice) request(ctx context.Context, url, apikey, hash, code, timestamp string, request map[string]interface{}) ([]byte, error) {
	logger := a.logger.WithFields(map[string]interface{}{
		"url":  url,
		"code": code,
		"hash": hash,
	}).WithField("func", "request")
	logger.Info("start request")
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(configs.Cfg.App.HttpTimeout))
	defer cancel()
	results, err := utils.Post(ctx, url, apikey, hash, code, timestamp, request)
	if err != nil {
		logger.WithError(err).Error("post")
		return nil, errors.New(errors.UpstreamInternalFailed, authErr.ErrUpstreamInternal)
	}
	defer results.Body.Close()
	body, err := ioutil.ReadAll(results.Body)
	if err != nil {
		logger.WithError(err).Error("read body")
		return nil, constant.ErrInternal
	}
	// 404
	if results.StatusCode == http.StatusNotFound {
		logger.WithError(fmt.Errorf(string(body))).Error("not found")
		var resp constant.ErrorResp
		if err := json.Unmarshal(body, &resp); err != nil {
			logger.WithError(fmt.Errorf(string(body))).Error("not found json un marshal")
			return nil, errors.New(errors.UpstreamInternalFailed, authErr.ErrUpstreamInternal)
		}
		return nil, errors.New(errors.NotFound, resp.Message)
	}
	// 403
	if results.StatusCode == http.StatusForbidden {
		logger.WithError(fmt.Errorf(string(body))).Error("forbidden")
		return nil, errors.New(errors.Authentication, authErr.ErrUpstreamForbidden)
	}
	return body, nil
}

// hash 生成签名
func (a *notice) hash(str map[string]interface{}, path, timestamp string, project *entity.Project) (string, error) {
	hashStr := make(map[string]interface{}, len(str))
	for k, v := range str {
		hashStr[fmt.Sprintf("body_%s", k)] = v
	}
	hashStr["path_url"] = path
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(hashStr)

	apiSecret, err := aes.Decode(project.ApiSecret, configs.Cfg.Project.SecretPwd)
	if err != nil {
		return "", err
	}
	oriTextHashBytes := sha256.Sum256([]byte(strings.TrimRight(bf.String(), "\n") + timestamp + apiSecret))
	return hex.EncodeToString(oriTextHashBytes[:]), nil
}
