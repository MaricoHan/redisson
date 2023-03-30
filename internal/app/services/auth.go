package services

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"

	mapset "github.com/deckarep/golang-set"
	httptransport "github.com/go-kit/kit/transport/http"
	log "github.com/sirupsen/logrus"
	"gitlab.bianjie.ai/avata/utils/commons/aes"
	"gitlab.bianjie.ai/avata/utils/errors"
	authErr "gitlab.bianjie.ai/avata/utils/errors/auth"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/repository/db/project"
	"gitlab.bianjie.ai/avata/open-api/internal/app/repository/db/service_redirect_url"
	userRepo "gitlab.bianjie.ai/avata/open-api/internal/app/repository/db/user"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/configs"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	"gitlab.bianjie.ai/avata/open-api/utils"
)

type IAuth interface {
	Verify(ctx context.Context, verify *vo.AuthVerify) (*dto.AuthVerify, error)
	GetUser(ctx context.Context, user *vo.AuthGetUser) ([]*dto.AuthGetUser, error)
}

type auth struct {
	logger *log.Logger
}

func NewAuth(logger *log.Logger) *auth {
	return &auth{logger: logger}
}

func (a *auth) Verify(ctx context.Context, params *vo.AuthVerify) (*dto.AuthVerify, error) {
	logger := a.logger.WithContext(ctx).WithField("params", params).WithField("func", "verify")
	path := ctx.Value(httptransport.ContextKeyRequestPath).(string)[len(configs.Cfg.App.RouterPrefix)+1:]
	res := &dto.AuthVerify{}
	project, err := a.getProject(params.ProjectID)
	if err != nil {
		logger.WithError(err).Error("query project")
		return res, err
	}
	authData, err := utils.AuthData(ctx)
	if err != nil {
		logger.WithError(err).Error("query auth data")
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
	request, err := params.Map()
	if err != nil {
		logger.WithError(err).Error("struct to map")
		return res, errors.ErrInternal
	}
	timestamp := utils.TimeToUnix(time.Now())
	hash, err := a.hash(request, path, timestamp, &project)
	if err != nil {
		logger.WithError(err).Error("hash")
		return res, errors.ErrInternal
	}
	body, err := a.request(context.Background(), fmt.Sprintf("%s%s?hash=%s&project_id=%s&type=%s", url.Url, path, params.Hash, params.ProjectID, params.Type), project.ApiKey, hash, user.Code, timestamp, nil)
	if err != nil {
		return res, err
	}
	var service dto.AuthUpstreamVerify
	if err := json.Unmarshal(body, &service); err != nil {
		logger.WithError(err).Error("body un marshal")
		return res, errors.New(errors.UpstreamInternalFailed, authErr.ErrUpstreamInternal)
	}
	if service.Data.Exists != dto.AuthVerifyExists {
		return res, errors.New(errors.UpstreamInternalFailed, authErr.ErrAuthVerifyExists)
	}
	res.Exists = service.Data.Exists
	return res, nil
}

func (a *auth) GetUser(ctx context.Context, params *vo.AuthGetUser) ([]*dto.AuthGetUser, error) {
	logger := a.logger.WithContext(ctx).WithField("params", params).WithField("func", "get user")
	path := ctx.Value(httptransport.ContextKeyRequestPath).(string)[len(configs.Cfg.App.RouterPrefix)+1:]
	var res []*dto.AuthGetUser
	project, err := a.getProject(params.ProjectID)
	if err != nil {
		logger.WithError(err).Error("query project")
		return res, err
	}
	authData, err := utils.AuthData(ctx)
	if err != nil {
		logger.WithError(err).Error("query auth data")
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
	request, err := params.Map()
	if err != nil {
		logger.WithError(err).Error("struct to map")
		return res, errors.ErrInternal
	}
	timestamp := utils.TimeToUnix(time.Now())
	hash, err := a.hash(request, path, timestamp, &project)
	if err != nil {
		logger.WithError(err).Error("hash")
		return res, errors.ErrInternal
	}
	body, err := a.request(context.Background(), fmt.Sprintf("%s%s?hash=%s&project_id=%s&type=%s", url.Url, path, params.Hash, params.ProjectID, params.Type), project.ApiKey, hash, user.Code, timestamp, nil)
	if err != nil {
		return res, err
	}
	var service dto.AuthUpstreamGetUser
	if err := json.Unmarshal(body, &service); err != nil {
		logger.WithError(err).Error("body un marshal")
		return res, errors.New(errors.UpstreamInternalFailed, authErr.ErrUpstreamInternal)
	}
	if len(service.Data) == 0 {
		return []*dto.AuthGetUser{}, nil
	}
	for _, v := range service.Data {
		if v.Address == "" {
			return res, errors.New(errors.UpstreamInternalFailed, authErr.ErrAuthUserAddress)
		}
		chainName := mapset.NewSetFromSlice(dto.AuthChainName)
		if !chainName.Contains(v.ChainName) {
			return res, errors.New(errors.UpstreamInternalFailed, authErr.ErrAuthUserChainName)
		}
		res = append(res, &dto.AuthGetUser{
			Address:   v.Address,
			ChainName: v.ChainName,
		})
	}
	return res, nil
}

// getProject 获取项目信息
func (a *auth) getProject(projectCode string) (entity.Project, error) {
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
func (a *auth) getUser(userID uint64) (entity.User, error) {
	userRepo := userRepo.NewUserRepo(initialize.MysqlDB)
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
func (a *auth) getServiceRedirectUrl(projectID uint64) (entity.ServiceRedirectUrl, error) {
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
func (a *auth) request(ctx context.Context, url, apikey, hash, code, timestamp string, request map[string]interface{}) ([]byte, error) {
	logger := a.logger.WithContext(ctx).WithFields(map[string]interface{}{
		"url":  url,
		"code": code,
		"hash": hash,
	}).WithField("func", "request").WithContext(ctx)
	logger.Info("start request")
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(configs.Cfg.App.HttpTimeout))
	defer cancel()
	results, err := utils.Get(ctx, url, apikey, hash, code, timestamp, request)
	if err != nil {
		logger.WithError(err).Error("get")
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
func (a *auth) hash(str map[string]interface{}, path, timestamp string, project *entity.Project) (string, error) {
	hashStr := make(map[string]interface{}, len(str))
	for k, v := range str {
		hashStr[fmt.Sprintf("query_%s", k)] = v
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
