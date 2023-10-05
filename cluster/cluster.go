package cluster

import (
	"gomemory/config"
	database2 "gomemory/database"
	"gomemory/datastruct/dict"
	"gomemory/datastruct/set"
	"gomemory/interface/cache"
	"gomemory/interface/server"
	"gomemory/lib/idgenerator"
	"gomemory/server/parser"
	"sync"
)

type CmdLine = [][]byte

// hostSlot stores status of host which hosted by current node
type hostSlot struct {
	state uint32
	mu    sync.RWMutex
	// OldNodeID is the node which is moving out this slot
	// only valid during slot is importing
	oldNodeID string
	// OldNodeID is the node which is importing this slot
	// only valid during slot is moving out
	newNodeID string

	/* importedKeys stores imported keys during migrating progress
	 * While this slot is migrating, if importedKeys does not have the given key, then current node will import key before execute commands
	 *
	 * In a migrating slot, the slot on the old node is immutable, we only delete a key in the new node.
	 * Therefore, we must distinguish between non-migrated key and deleted key.
	 * Even if a key has been deleted, it still exists in importedKeys, so we can distinguish between non-migrated and deleted.
	 */
	importedKeys *set.Set
	// keys stores all keys in this slot
	// Cluster.makeInsertCallback and Cluster.makeDeleteCallback will keep keys up to time
	keys *set.Set
}

type peerClient interface {
	Send(args [][]byte) server.Reply
}

type peerStream interface {
	Stream() <-chan *parser.Payload
	Close() error
}

type clientFactory interface {
	GetPeerClient(peerAddr string) (peerClient, error)
	ReturnPeerClient(peerAddr string, peerClient peerClient) error
	NewStream(peerAddr string, cmdLine CmdLine) (peerStream, error)
	Close() error
}

type CacheCluster struct {
	self          string
	addr          string
	db            cache.DBEngine
	transactions  *dict.SimpleDict
	transactionMu sync.RWMutex
	topology      topology
	slotMu        sync.RWMutex
	slots         map[uint32]hostSlot
	idGenerator   *idgenerator.IDGenerator
	clientFactory clientFactory
}

func MakeCacheCluster() *CacheCluster {
	cluster := &CacheCluster{
		self:          config.CacheProperties.Self,
		addr:          config.CacheProperties.AnnounceAddress(),
		db:            database2.NewStandaloneCacheServer(),
		transactions:  dict.MakeSimple(),
		idGenerator:   idgenerator.MakeGenerator(config.CacheProperties.Self),
		clientFactory: newDefaultClientFactory(),
	}
	cluster.topology = newEtcdTopology(cluster)

	var err error
	// 如果是主节点
	if config.CacheProperties.ClusterAsSeed {
		err = cluster.startAsSeed(config.CacheProperties.AnnounceAddress())
	} else {
		err = cluster.join(config.CacheProperties.ClusterSeed)
	}
	if err != nil {
		panic(err)
	}
	return cluster
}

func (cluster *CacheCluster) Exec(c server.Connection, cmdLine CmdLine) server.Reply {
	return nil
}

func (cluster *CacheCluster) AfterClientClose(c server.Connection) {
}

func (cluster *CacheCluster) Close() {
}
