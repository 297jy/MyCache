package cluster

func (cluster *Cluster) getHostSlot(slotId uint32) *hostSlot {
	cluster.slotMu.RLock()
	defer cluster.slotMu.RUnlock()
	return cluster.slots[slotId]
}
