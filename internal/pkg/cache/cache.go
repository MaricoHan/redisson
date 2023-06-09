package cache

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/repository/db/auth"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/app/repository/db/chain"
	"gitlab.bianjie.ai/avata/open-api/internal/app/repository/db/project"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
)

type cache struct {
}

var (
	cacheOnce sync.Once
	caches    *cache
)

func NewCache() *cache {
	cacheOnce.Do(func() {
		caches = &cache{}
	})
	return caches
}

// Project 返回项目信息且缓存项目
func (c cache) Project(key string) (authData vo.AuthData, err error) {
	var projectInfo entity.Project
	//err = initialize.RedisClient.GetObject(fmt.Sprintf("%s%s", constant.KeyProjectApikey, key), &authData)
	if err != nil {
		return authData, errors.Wrap(err, "get project from cache")
	}
	if projectInfo.Id < 1 {
		// 查询项目信息
		projectRepo := project.NewProjectRepo(initialize.MysqlDB)
		projectInfo, err = projectRepo.GetProjectByApiKey(key)
		if err != nil {
			return authData, errors.Wrap(err, "get project from db")
		}

		// 查询链信息
		chainInfo, err := c.Chain(projectInfo.ChainId)
		if err != nil {
			log.WithError(err).Error("get chain from cache")
			return authData, errors.Wrap(err, "get project from db")
		}

		authData = vo.AuthData{
			ProjectId:          uint64(projectInfo.Id),
			ChainId:            uint64(chainInfo.Id),
			PlatformId:         uint64(projectInfo.UserId),
			Module:             chainInfo.Module,
			Code:               chainInfo.Code,
			AccessMode:         projectInfo.AccessMode,
			UserId:             uint64(projectInfo.UserId),
			ApiSecret:          projectInfo.ApiSecret,
			ExistWalletService: false,
		}

		// 查询是否开通钱包服务
		existWalletService, err := projectRepo.ExistServices(projectInfo.Id, constant.ServiceTypeWallet)
		if existWalletService {
			authData.ExistWalletService = true
		}

		// 增加缓存
		if err := initialize.RedisClient.SetObject(fmt.Sprintf("%s%s", constant.KeyProjectApikey, key), authData, time.Minute*5); err != nil {
			return authData, errors.Wrap(err, "save auth cache")
		}

	}
	return authData, nil
}

// Chain 返回链信息且缓存链信息
func (c cache) Chain(chainID uint) (entity.Chain, error) {
	var chainInfo entity.Chain
	err := initialize.RedisClient.GetObject(fmt.Sprintf("%s%d", constant.KeyChain, chainID), &chainInfo)
	if err != nil {
		return chainInfo, errors.Wrap(err, "get chain from cache")
	}
	if chainInfo.Id < 1 {
		// 获取链信息
		chainRepo := chain.NewChainRepo(initialize.MysqlDB)
		chainInfo, err = chainRepo.QueryChainById(uint64(chainID))
		if err != nil {
			return chainInfo, errors.Wrap(err, "query chain from db by id")
		}
		if chainInfo.Id > 0 {
			// save cache
			if err := initialize.RedisClient.SetObject(fmt.Sprintf("%s%d", constant.KeyChain, chainID), chainInfo, time.Minute*5); err != nil {
				return chainInfo, errors.Wrap(err, "save project to cache")
			}
		}
	}
	return chainInfo, nil
}

//
//  ProjectAuth
//  @Description: 获取项目权限
//  @param key
//
func (c cache) ProjectAuth(pid int) (list []entity.Permission, err error) {
	// 查询缓存
	key := fmt.Sprintf("%s:%d", constant.KeyAuth, pid)
	//err = initialize.RedisClient.GetObject(key, &list)
	//if err == nil && len(list) > 0 {
	//	// 有缓存
	//	return list, err
	//}

	// 无缓存 查询数据库
	projectAuthRepo := auth.NewProjectAuthRepo(initialize.MysqlDB)
	list, err = projectAuthRepo.GetProjectPermission(pid)
	if err != nil {
		if err != nil {
			return list, errors.Wrap(err, "get project auth from db")
		}
	}

	// 增加缓存
	if err = initialize.RedisClient.SetObject(key, list, time.Minute*5); err != nil {
		return list, errors.Wrap(err, "save project auth cache")
	}

	return list, nil
}
