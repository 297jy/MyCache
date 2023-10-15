package cluster

import "gomemory/server/protocol"

func (cluster *Cluster) startAsSeed(listenAddr string) protocol.ErrorReply {
	err := cluster.topology.StartAsSeed(listenAddr)
	if err != nil {
		return err
	}
	for i := 0; i < slotCount; i++ {
		cluster.initSlot(uint32(i), slotStateHost)
	}
	return nil
}

func (cluster *Cluster) join(seedAddr string) protocol.ErrorReply {
	return nil
}
