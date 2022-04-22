package etcd

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/configs"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"strings"
	"time"
)

var cli *clientv3.Client

//etcd解析器
type etcdResolver struct {
	etcdAddr   string
	clientConn resolver.ClientConn
	logger     *log.Logger
}

//初始化一个etcd解析器
func NewResolver(etcdAddr string, logger *log.Logger) resolver.Builder {
	return &etcdResolver{etcdAddr: etcdAddr, logger: logger}
}

func (r *etcdResolver) Scheme() string {
	return constant.Schema
}

//构建解析器 grpc.Dial()同步调用
func (r *etcdResolver) Build(target resolver.Target, clientConn resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var err error
	logFields := log.Fields{}
	logFields["model"] = "etcd"
	logFields["func"] = "Build"
	r.logger.WithFields(logFields).Debugf("targer endpoint:%s,Authority:%s,Scheme:%s", target.Endpoint, target.Authority, target.Scheme)
	//构建etcd client
	if cli == nil {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   strings.Split(r.etcdAddr, ";"),
			DialTimeout: 15 * time.Second,
			Username:    configs.Cfg.Etcd.Username,
			Password:    configs.Cfg.Etcd.Password,
		})
		if err != nil {
			r.logger.WithFields(logFields).Infof("连接etcd失败：%s\n", err.Error())
			return nil, err
		}
	}
	resolver.Register(r)

	r.clientConn = clientConn
	key := fmt.Sprintf("/%s/%s/", target.Scheme, target.Endpoint)
	go r.watch(key)

	return r, nil
}

func (r *etcdResolver) watch(keyPrefix string) {
	//初始化服务地址列表
	var addrList []resolver.Address
	logFields := log.Fields{}
	logFields["model"] = "etcd"
	logFields["func"] = "watch"
	r.logger.WithFields(logFields).Debugf("etcd watch key :%s", keyPrefix)
	resp, err := cli.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		r.logger.WithFields(logFields).Infof("get service addr list err:%s", err.Error())
	} else {
		for i := range resp.Kvs {
			addrList = append(addrList, resolver.Address{Addr: strings.TrimPrefix(string(resp.Kvs[i].Key), keyPrefix)})
		}
	}

	r.clientConn.NewAddress(addrList)

	//监听服务地址列表的变化
	rch := cli.Watch(context.Background(), keyPrefix, clientv3.WithPrefix())
	for n := range rch {
		for _, ev := range n.Events {
			addr := strings.TrimPrefix(string(ev.Kv.Key), keyPrefix)
			switch ev.Type {
			case mvccpb.PUT:
				if !exists(addrList, addr) {
					addrList = append(addrList, resolver.Address{Addr: addr})
					r.clientConn.NewAddress(addrList)
				}
			case mvccpb.DELETE:
				if s, ok := remove(addrList, addr); ok {
					addrList = s
					r.clientConn.NewAddress(addrList)
				}
			}
		}
	}
	r.logger.WithFields(logFields).Infof("addrList :%s", addrList)
}

func exists(l []resolver.Address, addr string) bool {
	for i := range l {
		if l[i].Addr == addr {
			return true
		}
	}
	return false
}

func remove(s []resolver.Address, addr string) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}

//watch有变化后调用
func (r *etcdResolver) ResolveNow(rn resolver.ResolveNowOptions) {
	logFields := log.Fields{}
	logFields["model"] = "etcd"
	logFields["func"] = "ResolveNow"
	//r.logger.WithFields(logFields).Debug("etcd resolver now")
}

//解析器关闭时调用
func (r *etcdResolver) Close() {
	logFields := log.Fields{}
	logFields["model"] = "etcd"
	logFields["func"] = "Close"
	r.logger.WithFields(logFields).Debug("etcd resolver stop")
}
