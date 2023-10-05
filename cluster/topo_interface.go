package cluster

import (
	"gomemory/server/protocol"
	"time"
)

// Slot represents a hash slot,  used in cluster internal messages
type Slot struct {
	// ID is uint between 0 and 16383
	ID uint32
	// NodeID is id of the hosting node
	// If the slot is migrating, NodeID is the id of the node importing this slot (target node)
	NodeID string
	// Flags stores more information of slot
	Flags uint32
}

// Node represents a node and its slots, used in cluster internal messages
type Node struct {
	ID        string
	Addr      string
	Slots     []*Slot // ascending order by slot id
	Flags     uint32
	lastHeard time.Time
}

type topology interface {
	GetSelfNodeID() string
	GetNodes() []*Node // return a copy
	GetNode(nodeID string) *Node
	GetSlots() []*Slot
	StartAsSeed(addr string) protocol.ErrorReply
	SetSlot(slotIDs []uint32, newNodeID string) protocol.ErrorReply
	LoadConfigFile() protocol.ErrorReply
	Join(seed string) protocol.ErrorReply
	Close() error
}