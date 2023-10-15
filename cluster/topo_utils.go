package cluster

func (cluster *Cluster) getHostSlot(slotId uint32) *hostSlot {
	cluster.slotMu.RLock()
	defer cluster.slotMu.RUnlock()
	return cluster.slots[slotId]
}

func (cluster *Cluster) initSlot(slotId uint32, state uint32) {
	cluster.slotMu.Lock()
	defer cluster.slotMu.Unlock()
	cluster.slots[slotId] = &hostSlot{}
}
