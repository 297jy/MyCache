package database

import (
	"gomemory/interface/database"
	"gomemory/interface/server"
	"time"
)

func (server *CacheServer) Exec(c server.Connection, cmdLine [][]byte) (result server.Reply) {
	return nil
}

// AfterClientClose does some clean after client close connection
func (server *CacheServer) AfterClientClose(c server.Connection) {
}

// Close graceful shutdown database
func (server *CacheServer) Close() {

}

// ForEach traverses all the keys in the given database
func (server *CacheServer) ForEach(dbIndex int, cb func(key string, data *database.DataEntity, expiration *time.Time) bool) {
}

// GetEntity returns the data entity to the given key
func (server *CacheServer) GetEntity(dbIndex int, key string) (*database.DataEntity, bool) {
	return nil, false
}

func (server *CacheServer) GetExpiration(dbIndex int, key string) *time.Time {
	return nil
}

// ExecMulti executes multi commands transaction Atomically and Isolated
func (server *CacheServer) ExecMulti(conn server.Connection, watching map[string]uint32, cmdLines []database.CmdLine) server.Reply {
	return nil
}

// RWLocks lock keys for writing and reading
func (server *CacheServer) RWLocks(dbIndex int, writeKeys []string, readKeys []string) {
}

// RWUnLocks unlock keys for writing and reading
func (server *CacheServer) RWUnLocks(dbIndex int, writeKeys []string, readKeys []string) {
}

// GetUndoLogs return rollback commands
func (server *CacheServer) GetUndoLogs(dbIndex int, cmdLine [][]byte) []database.CmdLine {
	return nil
}

// ExecWithLock executes normal commands, invoker should provide locks
func (server *CacheServer) ExecWithLock(conn server.Connection, cmdLine [][]byte) server.Reply {
	return nil
}

// GetDBSize returns keys count and ttl key count
func (server *CacheServer) GetDBSize(dbIndex int) (int, int) {
	return 0, 0
}

func (server *CacheServer) SetKeyInsertedCallback(cb database.KeyEventCallback) {

}

func (server *CacheServer) SetKeyDeletedCallback(cb database.KeyEventCallback) {

}
