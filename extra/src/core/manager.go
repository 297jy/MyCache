package core

import (
	"EasyMemcache/src/cache"
	"EasyMemcache/src/common"
	"sort"
	"strconv"
	"sync"
)

type NodeRecordManage interface {
	// 更新版本信息
	UpdateRecord()
	// 清除无效低版本信息
	ClearRecord()
	// 检查key是否因为hash漂移而失效
	CheckValid(serverName string, key string, version int) bool
}

type HttpServer interface {
	Put(key string, value string) (ok bool)

	Get(key string) (value string)
	// 路由key可能存在的服务器节点
	Route(key string) *Server
	// 重新加载服务器节点
	ReloadRouter()
}

type Hash func(data []byte) uint32

type neighborRecord struct {
	// 不同版本号对应的邻居节点的散列值列表
	lastHash int
	// 到当前版本号最大的散列值
	maxHashs []int
	// 版本号到hashs下标的映射
	versionIndex map[int]int
}

type neighborRecordMap struct {
	// 全局最大版本号
	maxVersion int
	data       map[string]*neighborRecord
	rwMu       sync.RWMutex
}

type consistentHash struct {
	hash      Hash
	keys      []int
	serverMap map[int]Server
	// 本地虚拟节点名与下标映射，用来快速找到虚拟节点名的左邻居节点
	localIndexMap map[string]int
}

type ServerContainer struct {
	mu          sync.Mutex // 全局锁
	provides    []*Server
	localServer *Server
	hashMap     *consistentHash
	recordMap   neighborRecordMap
	groupManage cache.GroupManage
}

func (m *ServerContainer) Route(key string) *Server {
	if len(m.hashMap.keys) == 0 {
		return nil
	}

	h := int(m.hashMap.hash([]byte(key)))
	idx := sort.Search(len(m.hashMap.keys), func(i int) bool {
		return m.hashMap.keys[i] >= h
	})

	s := m.hashMap.serverMap[m.hashMap.keys[idx%len(m.hashMap.keys)]]
	return &s
}

func (m *ServerContainer) ReloadRouter() {
	m.hashMap.serverMap = map[int]Server{}
	m.hashMap.keys = m.hashMap.keys[0:0]

	for i := 0; i < len(m.provides); i++ {
		for j := 0; j < m.provides[i].replicas; j++ {
			h := int(m.hashMap.hash([]byte(m.provides[i].Name + strconv.Itoa(j))))
			m.hashMap.keys = append(m.hashMap.keys, h)
			m.hashMap.serverMap[h] = Server{
				Name: m.provides[i].Name + strconv.Itoa(j),
				Ip:   m.provides[i].Ip,
			}
		}
	}
	sort.Ints(m.hashMap.keys)

	// 找到每个本地虚拟节点对应的左邻居节点
	for i := 0; i < len(m.hashMap.keys); i++ {
		s := m.hashMap.serverMap[m.hashMap.keys[i]]
		if m.localServer.Ip == s.Ip {
			m.hashMap.localIndexMap[s.Name] = (i + len(m.hashMap.keys) - 1) % len(m.hashMap.keys)
		}
	}
}

func (m *ServerContainer) UpdateRecord() {
	m.recordMap.rwMu.Lock()
	defer m.recordMap.rwMu.Unlock()

	if m.recordMap.maxVersion == 0 {
		m.recordMap.maxVersion = 1
	}

	for name, idx := range m.hashMap.localIndexMap {
		neighborHash := m.hashMap.keys[idx]
		record, ok := m.recordMap.data[name]
		if !ok {
			m.recordMap.maxVersion++
			m.recordMap.data[name] = &neighborRecord{
				lastHash:     neighborHash,
				maxHashs:     []int{neighborHash},
				versionIndex: map[int]int{m.recordMap.maxVersion: 0},
			}
		} else {
			if record.lastHash != neighborHash {
				m.recordMap.maxVersion++
				record.versionIndex[m.recordMap.maxVersion] = neighborHash
				record.maxHashs = append(record.maxHashs, neighborHash)
				for i := len(record.maxHashs) - 3; i >= 0; i-- {
					record.maxHashs[i] = common.Max(record.maxHashs[i], record.maxHashs[i+1])
				}
			}
		}
	}
}

func (m *ServerContainer) ClearRecord() {
	m.mu.Lock()
	defer m.mu.Unlock()
	// 获取可删除版本号的列表
	versions := m.groupManage.ClearableVersionNum()
	for version := range versions {
		for _, v := range m.recordMap.data {
			index, ok := v.versionIndex[version]
			if ok {
				for oldV, oldIndex := range v.versionIndex {
					if index < oldIndex {
						v.versionIndex[oldV] = oldIndex - 1
					}
				}

				delete(v.versionIndex, index)
				v.maxHashs = append(v.maxHashs[0:index], v.maxHashs[index+1:len(v.maxHashs)]...)
			}
		}
	}
}

func (m *ServerContainer) CheckValid(serverName string, key string, version int) bool {
	m.recordMap.rwMu.Lock()
	defer m.recordMap.rwMu.Unlock()

	h := m.hashMap.hash([]byte(key))
	record := m.recordMap.data[serverName]
	if record == nil {
		return false
	}
	index := record.versionIndex[version]
	// 最近的两个版本不会发生哈希漂移
	if (len(record.maxHashs) - index) <= 1 {
		return true
	}
	return int(h) > record.maxHashs[index]
}
