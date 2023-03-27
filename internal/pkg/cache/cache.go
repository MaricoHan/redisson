package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

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

// Project 返回项目信息切缓存项目
func (c cache) Project(key string) (entity.Project, bool, error) {
	var projectInfo entity.Project
	existWalletService := false
	err := initialize.RedisClient.GetObject(fmt.Sprintf("%s%s", constant.KeyProjectApikey, key), &projectInfo)
	if err != nil {
		return projectInfo, existWalletService, errors.Wrap(err, "get project from cache")
	}
	err = initialize.RedisClient.GetObject(fmt.Sprintf("%s%s", constant.KeyExistWalletService, key), &existWalletService)
	if err != nil {
		return projectInfo, existWalletService, errors.Wrap(err, "get project exist wallet services from cache")
	}
	if projectInfo.Id < 1 {
		// 查询 project 信息
		projectRepo := project.NewProjectRepo(initialize.MysqlDB)
		projectInfo, err = projectRepo.GetProjectByApiKey(key)
		if err != nil {
			return projectInfo, existWalletService, errors.Wrap(err, "get project from db")
		}

		if projectInfo.Id > 0 {
			// save cache
			if err := initialize.RedisClient.SetObject(fmt.Sprintf("%s%s", constant.KeyProjectApikey, key), projectInfo, time.Minute*5); err != nil {
				return projectInfo, existWalletService, errors.Wrap(err, "save project cache")
			}
		}

		// project 关联的 serviceIds
		existWalletService, err = projectRepo.ExistServices(projectInfo.Id, entity.ServiceTypeWallet)
		if err != nil {
			return projectInfo, existWalletService, errors.Wrap(err, "get project services from db")
		}
		if err := initialize.RedisClient.SetObject(fmt.Sprintf("%s%s", constant.KeyExistWalletService, key), existWalletService, time.Minute*5); err != nil {
			return projectInfo, existWalletService, errors.Wrap(err, "save project exist wallet services cache")
		}
	}
	logrus.Warnln("====",existWalletService)
	return projectInfo, existWalletService, nil
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
