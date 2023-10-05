package cache

import (
	"gomemory/interface/database"
	"gomemory/interface/server"
	"time"
)

// DBEngine 存储引擎，暴露操作数据的上层接口
type DBEngine interface {
	database.DB
	ExecWithLock(conn server.Connection, cmdLine database.CmdLine) server.Reply
	ExecMulti(conn server.Connection, watching map[string]uint32, cmdLines []database.CmdLine) server.Reply
	GetUndoLogs(dbIndex int, cmdLine [][]byte) []database.CmdLine
	ForEach(dbIndex int, cb func(key string, data *database.DataEntity, expiration *time.Time) bool)
	RWLocks(dbIndex int, writeKeys []string, readKeys []string)
	RWUnLocks(dbIndex int, writeKeys []string, readKeys []string)
	GetDBSize(dbIndex int) (int, int)
	GetEntity(dbIndex int, key string) (*database.DataEntity, bool)
	GetExpiration(dbIndex int, key string) *time.Time
	SetKeyInsertedCallback(cb database.KeyEventCallback)
	SetKeyDeletedCallback(cb database.KeyEventCallback)
}
