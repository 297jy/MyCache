package cluster

import "gomemory/server/protocol"

func (cluster *CacheCluster) startAsSeed(listenAddr string) protocol.ErrorReply {
	return nil
}

func (cluster *CacheCluster) join(seedAddr string) protocol.ErrorReply {
	return nil
}
