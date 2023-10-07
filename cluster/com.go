package cluster

import (
	"gomemory/interface/server"
	"gomemory/server/protocol"
)

/**
func (cluster *Cluster) ensureKeyWithoutLock(key string) protocol.ErrorReply {
	cluster.db.RWLocks(0, []string{key}, nil)
	defer cluster.db.RWUnLocks(0, []string{key}, nil)
	return cluster.ensureKey(key)
}

func (cluster *Cluster) ensureKey(key string) protocol.ErrorReply {
	slotId := getSlot(key)
	cluster.slotMu.RLock()
	slot := cluster.slots[slotId]
	cluster.slotMu.RUnlock()
	if slot == nil {
		return nil
	}
	if slot.state != slotStateImporting || slot.importedKeys.Has(key) {
		return nil
	}

	resp := cluster.relay(slot.oldNodeID, connection.NewFakeConn(), utils.ToCmdLine("DumpKey_", key))
}
**/

func (cluster *Cluster) relay(peerId string, conn server.Connection, cmdLine [][]byte) server.Reply {
	if peerId == cluster.self {
		return cluster.Exec(conn, cmdLine)
	}

	cli, err := cluster.clientFactory.GetPeerClient(peerId)
	if err != nil {
		return protocol.MakeErrReply(err.Error())
	}
	defer func() {
		_ = cluster.clientFactory.ReturnPeerClient(peerId, cli)
	}()
	return cli.Send(cmdLine)
}
