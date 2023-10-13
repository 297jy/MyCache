package database

import (
	"gomemory/interface/database"
	"gomemory/interface/server"
	"gomemory/server/protocol"
)

const (
	upsertPolicy = iota // default
	insertPolicy        // set nx
	updatePolicy        // set ex
)

const unlimitedTTL int64 = 0

func execSet(db *DB, args [][]byte) server.Reply {
	key := string(args[0])
	value := args[1]
	//policy := upsertPolicy
	// ttl := unlimitedTTL

	/**
	arg := strings.ToUpper(string(args[2]))
	if arg == "NX" { // insert
		if policy == updatePolicy {
			return &protocol.SyntaxErrReply{}
		}
		policy = insertPolicy
	} else if arg == "XX" { // update policy
		if policy == insertPolicy {
			return &protocol.SyntaxErrReply{}
		}
		policy = updatePolicy
	} **/

	entity := &database.DataEntity{
		Data: value,
	}

	result := db.PutIfAbsent(key, entity)
	if result > 0 {
		return &protocol.OkReply{}
	}
	return &protocol.NullBulkReply{}
}

func execGet(db *DB, args [][]byte) server.Reply {
	key := string(args[0])
	bytes, err := db.getAsString(key)
	if err != nil {
		return err
	}
	if bytes == nil {
		return &protocol.NullBulkReply{}
	}
	return protocol.MakeBulkReply(bytes)
}

func (db *DB) getAsString(key string) ([]byte, protocol.ErrorReply) {
	entity, ok := db.GetEntity(key)
	if !ok {
		return nil, nil
	}
	bytes, ok := entity.Data.([]byte)
	if !ok {
		return nil, &protocol.WrongTypeErrReply{}
	}
	return bytes, nil
}

// GetEntity returns DataEntity bind to given key
func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	raw, ok := db.data.GetWithLock(key)
	if !ok {
		return nil, false
	}
	/**
	先不考虑过期的场景
	if db.IsExpired(key) {
		return nil, false
	}**/
	entity, _ := raw.(*database.DataEntity)
	return entity, true
}

func init() {
	registerCommand("Set", execSet, nil, nil, -3, flagWrite).
		attachCommandExtra([]string{redisFlagWrite, redisFlagDenyOOM}, 1, 1, 1)
	registerCommand("Get", execGet, nil, nil, 2, flagReadOnly).
		attachCommandExtra([]string{redisFlagReadonly, redisFlagFast}, 1, 1, 1)
}
