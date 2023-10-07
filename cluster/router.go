package cluster

import (
	"gomemory/interface/server"
	"strings"
)

var router = make(map[string]CmdFunc)

func registerCmd(name string, cmd CmdFunc) {
	name = strings.ToLower(name)
	router[name] = cmd
}

func defaultFunc(cluster *Cluster, c server.Connection, args [][]byte) server.Reply {
	key := string(args[1])
	slotId := getSlot(key)
	peer := cluster.pickNode(slotId)
	if peer.ID == cluster.self {
		return cluster.db.Exec(c, args)
	}
	return cluster.db.Exec(c, args)
}
func registerDefaultCmd(name string) {
	registerCmd(name, defaultFunc)
}
