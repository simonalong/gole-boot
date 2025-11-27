// Package grpc
/**
gole:///xxx-service 解析器
*/
package grpc

import (
	"sync"

	cmap "github.com/orcaman/concurrent-map"
	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/util"
	"google.golang.org/grpc/resolver"
)

const DefaultBaseBootGrpcScheme = "gole"

var GlobalResolver *BaseResolver

type BaseResolverBuilder struct{}
type BaseResolver struct {
	target     resolver.Target
	clientConn resolver.ClientConn
	// map[string] map[string]string
	// 示例：key: serviceName, key: serviceUniqName, value: 192.168.xx.xx:8080
	addressStore cmap.ConcurrentMap
}

func InitResolver() {
	GlobalResolver = &BaseResolver{
		addressStore: cmap.New(),
	}
	resolver.Register(&BaseResolverBuilder{})
}

func (*BaseResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	GlobalResolver.setTarget(target)
	GlobalResolver.setClientConn(cc)
	GlobalResolver.refreshState()
	return GlobalResolver, nil
}
func (*BaseResolverBuilder) Scheme() string { return DefaultBaseBootGrpcScheme }

func (baseResolver *BaseResolver) setTarget(target resolver.Target) {
	baseResolver.target = target
}

func (baseResolver *BaseResolver) setClientConn(clientConn resolver.ClientConn) {
	baseResolver.clientConn = clientConn
}

func (baseResolver *BaseResolver) setAddressStore(addressStore map[string]map[string]string) {
	for k, v := range addressStore {
		baseResolver.addressStore.Set(k, v)
	}
}

func (baseResolver *BaseResolver) addAddressStore(serviceName, serviceUniqName, address string) bool {
	val, exist := baseResolver.addressStore.Get(serviceName)
	if !exist {
		serviceUniqAddressMap := map[string]string{}
		serviceUniqAddressMap[serviceUniqName] = address
		baseResolver.addressStore.Set(serviceName, serviceUniqAddressMap)
	} else {
		serviceUniqAddressMap := val.(map[string]string)
		if valOld, exit := serviceUniqAddressMap[serviceUniqName]; exit {
			// 相同不更新
			if valOld == address {
				return false
			}
		}
		serviceUniqAddressMap[serviceUniqName] = address
		baseResolver.addressStore.Set(serviceName, serviceUniqAddressMap)
	}
	return true
}

func (baseResolver *BaseResolver) deleteAddressStore(serviceName, serviceUniqName string) {
	val, exist := baseResolver.addressStore.Get(serviceName)
	if !exist {
		return
	} else {
		serviceUniqAddressMap := val.(map[string]string)
		delete(serviceUniqAddressMap, serviceUniqName)
		baseResolver.addressStore.Set(serviceName, serviceUniqAddressMap)
	}
}

var lock sync.RWMutex

func (baseResolver *BaseResolver) refreshState() {
	lock.Lock()
	defer lock.Unlock()

	addrStrsVal, exit := baseResolver.addressStore.Get(baseResolver.target.Endpoint())
	if !exit {
		return
	}
	addrStrsMap := addrStrsVal.(map[string]string)
	var addrStrs []string
	for _, v := range addrStrsMap {
		addrStrs = append(addrStrs, v)
	}
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	err := baseResolver.clientConn.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		logger.Errorf("更新grpc解析名地址失败：: %v", err)
		return
	}
	logger.Infof("服务【%s】列表更新完毕，服务地址列表：%s", baseResolver.target.Endpoint(), util.ToJsonString(addrStrsVal))
}
func (*BaseResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*BaseResolver) Close()                                  {}
