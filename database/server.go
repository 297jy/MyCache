package database

import (
	"fmt"
	"gomemory/interface/database"
	"gomemory/interface/server"
	"gomemory/lib/logger"
	"gomemory/server/protocol"
	"runtime/debug"
	"time"
)

// Exec executes command
func (server *Standalone) Exec(c server.Connection, cmdLine [][]byte) (result server.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
			result = &protocol.UnknownErrReply{}
		}
	}()

	// cmdName := strings.ToLower(string(cmdLine[0]))
	// dbIndex := c.GetDBIndex()
	db := server.db.Load().(*DB)
	return db.Exec(c, cmdLine)
}

// AfterClientClose does some clean after client close connection
func (server *Standalone) AfterClientClose(c server.Connection) {
}

// Close graceful shutdown database
func (server *Standalone) Close() {

}

// ForEach traverses all the keys in the given database
func (server *Standalone) ForEach(dbIndex int, cb func(key string, data *database.DataEntity, expiration *time.Time) bool) {
}

// GetEntity returns the data entity to the given key
func (server *Standalone) GetEntity(dbIndex int, key string) (*database.DataEntity, bool) {
	return nil, false
}

func (server *Standalone) GetExpiration(dbIndex int, key string) *time.Time {
	return nil
}

// ExecMulti executes multi commands transaction Atomically and Isolated
func (server *Standalone) ExecMulti(conn server.Connection, watching map[string]uint32, cmdLines []database.CmdLine) server.Reply {
	return nil
}

// RWLocks lock keys for writing and reading
func (server *Standalone) RWLocks(dbIndex int, writeKeys []string, readKeys []string) {
}

// RWUnLocks unlock keys for writing and reading
func (server *Standalone) RWUnLocks(dbIndex int, writeKeys []string, readKeys []string) {
}

// GetUndoLogs return rollback commands
func (server *Standalone) GetUndoLogs(dbIndex int, cmdLine [][]byte) []database.CmdLine {
	return nil
}

// ExecWithLock executes normal commands, invoker should provide locks
func (server *Standalone) ExecWithLock(conn server.Connection, cmdLine [][]byte) server.Reply {
	return nil
}

// GetDBSize returns keys count and ttl key count
func (server *Standalone) GetDBSize(dbIndex int) (int, int) {
	return 0, 0
}

func (server *Standalone) SetKeyInsertedCallback(cb database.KeyEventCallback) {

}

func (server *Standalone) SetKeyDeletedCallback(cb database.KeyEventCallback) {

}
