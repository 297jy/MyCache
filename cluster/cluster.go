package cluster

import (
	"fmt"
	"gomemory/config"
	database2 "gomemory/database"
	"gomemory/datastruct/dict"
	"gomemory/datastruct/set"
	"gomemory/interface/cache"
	"gomemory/interface/server"
	"gomemory/lib/idgenerator"
	"gomemory/lib/logger"
	"gomemory/server/parser"
	"gomemory/server/protocol"
	"runtime/debug"
	"strings"
	"sync"
)

const slotCount int = 16384

const (
	slotStateHost = iota
	slotStateImporting
	slotStateMovingOut
)

type CmdLine = [][]byte

type CmdFunc func(cluster *Cluster, c server.Connection, cmdLine CmdLine) server.Reply

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

type Cluster struct {
	self          string
	addr          string
	db            cache.DBEngine
	transactions  *dict.SimpleDict
	transactionMu sync.RWMutex
	topology      topology
	slotMu        sync.RWMutex
	slots         map[uint32]*hostSlot
	idGenerator   *idgenerator.IDGenerator
	clientFactory clientFactory
}

func MakeCluster() *Cluster {
	cluster := &Cluster{
		self:          config.Properties.Self,
		addr:          config.Properties.AnnounceAddress(),
		db:            database2.NewStandaloneServer(),
		transactions:  dict.MakeSimple(),
		idGenerator:   idgenerator.MakeGenerator(config.Properties.Self),
		clientFactory: newDefaultClientFactory(),
	}
	cluster.topology = newRaftTopology(cluster)
	cluster.slots = make(map[uint32]*hostSlot)

	var err error
	// 如果是主节点
	if config.Properties.ClusterAsSeed {
		err = cluster.startAsSeed(config.Properties.AnnounceAddress())
	} else {
		err = cluster.join(config.Properties.ClusterSeed)
	}
	if err != nil {
		panic(err)
	}
	return cluster
}

func (cluster *Cluster) Exec(c server.Connection, cmdLine CmdLine) (result server.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
			result = &protocol.UnknownErrReply{}
		}
	}()

	cmdName := strings.ToLower(string(cmdLine[0]))
	cmdFunc, ok := router[cmdName]
	if !ok {
		return protocol.MakeErrReply("ERR unknown command '" + cmdName + "', or not supported in cluster mode")
	}
	result = cmdFunc(cluster, c, cmdLine)
	return
}

func (cluster *Cluster) AfterClientClose(c server.Connection) {
}

func (cluster *Cluster) Close() {
}

// 获取 key slot对应的node节点
func (cluster *Cluster) pickNode(slotID uint32) *Node {
	/**
	hSlot := cluster.getHostSlot(slotID)
	if hSlot != nil {
		switch hSlot.state {

		}
	}**/
	//todo 暂时忽略再平衡时的问题，后面补上
	slot := cluster.topology.GetSlots()[int(slotID)]
	node := cluster.topology.GetNode(slot.NodeID)
	return node
}
