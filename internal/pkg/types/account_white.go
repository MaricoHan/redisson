package types

import (
	"gitlab.bianjie.ai/irita-paas/open-api/config"
	"strconv"
	"strings"
)

var AccountWhiteList map[uint64]int64

// NewAccountWhiteListCache 初始化白名单缓存
func NewAccountWhiteListCache(conf config.Server) {
	whiteList := strings.Split(conf.AccountWhiteList, ",")
	accountCounts := strings.Split(conf.AccountCount, ",")
	if len(whiteList) == 0 {
		return
	}

	AccountWhiteList = make(map[uint64]int64, len(whiteList))

	for k, v := range whiteList {
		val, _ := strconv.ParseInt(v, 10, 64)
		accountCount, _ := strconv.ParseInt(accountCounts[k], 10, 64)
		AccountWhiteList[uint64(val)] = accountCount
	}
}
