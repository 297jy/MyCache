package database

import (
	"fmt"
	"gomemory/config"
	"gomemory/datastruct/dict"
	"gomemory/interface/database"
	"gomemory/interface/server"
	"gomemory/server/protocol"
	"os"
	"strings"
	"sync/atomic"
)

const (
	dataDictSize = 1 << 16
	ttlDictSize  = 1 << 10
)

// CmdLine is alias for [][]byte, represents a command line
type CmdLine = [][]byte

// ExecFunc is interface for command executor
// args don't include cmd line
type ExecFunc func(db *DB, args [][]byte) server.Reply

// PreFunc analyses command line when queued command to `multi`
// returns related write keys and read keys
type PreFunc func(args [][]byte) ([]string, []string)

// UndoFunc returns undo logs for the given command line
// execute from head to tail when undo
type UndoFunc func(db *DB, args [][]byte) []CmdLine

// Standalone 单机
type Standalone struct {
	db *atomic.Value
}

func NewStandaloneServer() *Standalone {
	server := &Standalone{}

	config.GetTempDir()
	err := os.MkdirAll(config.GetTempDir(), os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("create tmp dir failed: %v", err))
	}

	db := makeDB()
	holder := &atomic.Value{}
	holder.Store(db)
	server.db = holder

	return nil
}

// DB stores data and execute user's commands
type DB struct {
	index int
	// key -> DataEntity
	data *dict.ConcurrentDict
	// key -> expireTime (time.Time)
	ttlMap *dict.ConcurrentDict
	// key -> version(uint32)
	versionMap *dict.ConcurrentDict

	// callbacks
	insertCallback database.KeyEventCallback
	deleteCallback database.KeyEventCallback
}

// makeDB create DB instance
func makeDB() *DB {
	db := &DB{
		data:       dict.MakeConcurrent(dataDictSize),
		ttlMap:     dict.MakeConcurrent(ttlDictSize),
		versionMap: dict.MakeConcurrent(dataDictSize),
	}
	return db
}

func (db *DB) Exec(c server.Connection, cmdLine [][]byte) server.Reply {
	return db.execNormalCommand(cmdLine)
}

func (db *DB) execNormalCommand(cmdLine [][]byte) server.Reply {
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return protocol.MakeErrReply("ERR unknown command '" + cmdName + "'")
	}
	fun := cmd.executor
	return fun(db, cmdLine[1:])
}

func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	ret := db.data.PutIfAbsentWithLock(key, entity)
	return ret
}
