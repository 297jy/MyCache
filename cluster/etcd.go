package cluster

import "gomemory/server/protocol"

type EtcdTopology struct {
	cluster   *Cluster
	closeChan chan struct{}
}

func newEtcdTopology(cluster *Cluster) *EtcdTopology {
	return &EtcdTopology{
		cluster:   cluster,
		closeChan: make(chan struct{}),
	}
}

func (t *EtcdTopology) GetSelfNodeID() string {
	return ""
}

func (t *EtcdTopology) GetNodes() []*Node {
	return nil
}

func (t *EtcdTopology) GetNode(nodeID string) *Node {
	return nil
}

func (t *EtcdTopology) GetSlots() []*Slot {
	return nil
}

func (t *EtcdTopology) StartAsSeed(addr string) protocol.ErrorReply {
	return nil
}

func (t *EtcdTopology) SetSlot(slotIDs []uint32, newNodeID string) protocol.ErrorReply {
	return nil
}

func (t *EtcdTopology) LoadConfigFile() protocol.ErrorReply {
	return nil
}

func (t *EtcdTopology) Join(seed string) protocol.ErrorReply {
	return nil
}

func (t *EtcdTopology) Close() error {
	return nil
}
