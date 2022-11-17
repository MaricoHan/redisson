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

	mapset "github.com/deckarep/golang-set"
	httptransport "github.com/go-kit/kit/transport/http"
	log "github.com/sirupsen/logrus"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
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

type IAuth interface {
	Verify(ctx context.Context, verify *vo.AuthVerify) (*dto.AuthVerify, error)
	GetUser(ctx context.Context, user *vo.AuthGetUser) (*dto.AuthGetUser, error)
}

type auth struct {
	logger *log.Logger
}

func NewAuth(logger *log.Logger) *auth {
	return &auth{logger: logger}
}

func (a *auth) Verify(ctx context.Context, params *vo.AuthVerify) (*dto.AuthVerify, error) {
	logger := a.logger.WithField("params", params).WithField("func", "verify")
	path := ctx.Value(httptransport.ContextKeyRequestPath).(string)[len(configs.Cfg.App.RouterPrefix)+1:]
	res := &dto.AuthVerify{}
	project, err := a.getProject(params.ProjectID)
	if err != nil {
		logger.WithError(err).Error("query project")
		return res, err
	}
	user, err := a.getUser(uint64(project.UserId))
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
	hash, err := a.hash(request, &project)
	body, err := a.request(ctx, fmt.Sprintf("%s%s?hash=%s&project_id=%s", url.Url, path, params.Hash, params.ProjectID), project.ApiKey, hash, user.Code, nil)
	if err != nil {
		return res, err
	}
	var service dto.AuthUpstreamVerify
	if err := json.Unmarshal(body, &service); err != nil {
		logger.WithError(err).Error("body un marshal")
		return res, constant.ErrUpstreamInternalEntity
	}
	if service.Data.Exists != dto.AuthVerifyExists {
		return res, constant.ErrAuthVerifyExists
	}
	res.Exists = service.Data.Exists
	return res, nil
}

func (a *auth) GetUser(ctx context.Context, params *vo.AuthGetUser) (*dto.AuthGetUser, error) {
	logger := a.logger.WithField("params", params).WithField("func", "get user")
	path := ctx.Value(httptransport.ContextKeyRequestPath).(string)[len(configs.Cfg.App.RouterPrefix)+1:]
	res := &dto.AuthGetUser{}
	project, err := a.getProject(params.ProjectID)
	if err != nil {
		logger.WithError(err).Error("query project")
		return res, err
	}
	user, err := a.getUser(uint64(project.UserId))
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
	hash, err := a.hash(request, &project)
	body, err := a.request(ctx, fmt.Sprintf("%s%s?hash=%s&project_id=%s", url.Url, path, params.Hash, params.ProjectID), project.ApiKey, hash, user.Code, nil)
	if err != nil {
		return res, err
	}
	var service dto.AuthUpstreamGetUser
	if err := json.Unmarshal(body, &service); err != nil {
		logger.WithError(err).Error("body un marshal")
		return res, constant.ErrUpstreamInternalEntity
	}
	if service.Data.Address == "" {
		return res, constant.ErrAuthUserAddress
	}
	chainName := mapset.NewSetFromSlice(dto.AuthChainName)
	if !chainName.Contains(service.Data.ChainName) {
		return res, constant.ErrAuthUserChainName
	}
	res.Address = service.Data.Address
	res.ChainName = service.Data.ChainName
	return res, nil
}

// getProject 获取项目信息
func (a *auth) getProject(projectCode string) (entity.Project, error) {
	projectRepo := project.NewProjectRepo(initialize.MysqlDB)
	project, err := projectRepo.GetProjectByCode(projectCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return project, constant.ErrProjectOrUserNotFound
		}
		return project, err
	}
	return project, nil
}

// getUser 获取用户信息, id&&code
func (a *auth) getUser(userID uint64) (entity.User, error) {
	userRepo := user.NewUserRepo(initialize.MysqlDB)
	user, err := userRepo.GetUser(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, constant.ErrProjectOrUserNotFound
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
			return sru, constant.ErrServiceRedirectUrlNotFound
		}
		return sru, err
	}
	return sru, nil
}

// request 服务请求
// 1.向上游服务方发起请求, err不等于nil说明存在异常返回服务错误
// 2.当请求成功，返回的状态码为404时, 则返回NOT_FOUNT
func (a *auth) request(ctx context.Context, url, apikey, hash, code string, request map[string]interface{}) ([]byte, error) {
	logger := a.logger.WithFields(map[string]interface{}{
		"url":  url,
		"code": code,
	}).WithField("func", "request")
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(configs.Cfg.App.HttpTimeout))
	defer cancel()

	results, err := utils.Get(ctx, url, apikey, hash, code, request)
	if err != nil {
		logger.WithError(err).Error("get")
		return nil, constant.ErrInternal
	}
	defer results.Body.Close()
	body, err := ioutil.ReadAll(results.Body)
	if err != nil {
		logger.WithError(err).Error("read body")
		return nil, constant.ErrInternal
	}
	if results.StatusCode == http.StatusNotFound {
		logger.WithError(fmt.Errorf(string(body))).Error("not found")
		var resp constant.ErrorResp
		if err := json.Unmarshal(body, &resp); err != nil {
			logger.WithError(fmt.Errorf(string(body))).Error("not found json un marshal")
			return nil, constant.ErrInternal
		}
		return nil, constant.Register(constant.AuthCodeSpace, constant.NotFound, resp.Message)
	}
	return body, nil
}

// hash 生成签名
func (a *auth) hash(str map[string]interface{}, project *entity.Project) (string, error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(str)

	apiSecret, err := aes.Decode(project.ApiSecret, configs.Cfg.Project.SecretPwd)
	if err != nil {
		return "", err
	}
	oriTextHashBytes := sha256.Sum256([]byte(strings.TrimRight(bf.String(), "\n") + apiSecret))
	return hex.EncodeToString(oriTextHashBytes[:]), nil
}
