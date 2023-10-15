package cluster

import (
	"gomemory/server/protocol"
	"sync"
)

type RaftTopology struct {
	cluster    *Cluster
	mu         sync.RWMutex
	selfNodeID string
	slots      []*Slot
	leaderId   string
	closeChan  chan struct{}
	nodes      map[string]*Node
}

func newRaftTopology(cluster *Cluster) *RaftTopology {
	return &RaftTopology{
		cluster:   cluster,
		closeChan: make(chan struct{}),
	}
}

func (raft *RaftTopology) GetSelfNodeID() string {
	return ""
}

func (raft *RaftTopology) GetNodes() []*Node {
	return nil
}

func (raft *RaftTopology) GetNode(nodeID string) *Node {
	raft.mu.RLock()
	defer raft.mu.RUnlock()
	return raft.nodes[nodeID]
}

func (raft *RaftTopology) GetSlots() []*Slot {
	return raft.slots
}

func (raft *RaftTopology) StartAsSeed(listenAddr string) protocol.ErrorReply {
	selfNodeID := listenAddr
	raft.mu.Lock()
	defer raft.mu.Unlock()
	raft.slots = make([]*Slot, slotCount)
	for i := range raft.slots {
		raft.slots[i] = &Slot{
			ID:     uint32(i),
			NodeID: selfNodeID,
		}
	}
	raft.selfNodeID = selfNodeID
	raft.leaderId = selfNodeID

	raft.nodes = make(map[string]*Node)
	raft.nodes[selfNodeID] = &Node{
		ID:    selfNodeID,
		Addr:  listenAddr,
		Slots: raft.slots,
	}

	raft.cluster.self = selfNodeID
	return nil
}

func (raft *RaftTopology) SetSlot(slotIDs []uint32, newNodeID string) protocol.ErrorReply {
	return nil
}

func (raft *RaftTopology) LoadConfigFile() protocol.ErrorReply {
	return nil
}

func (raft *RaftTopology) Join(seed string) protocol.ErrorReply {
	return nil
}

func (raft *RaftTopology) Close() error {
	return nil
}
