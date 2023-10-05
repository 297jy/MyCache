package database

import "gomemory/interface/server"

type CmdLine = [][]byte

type DB interface {
	Exec(client server.Connection, cmdLine CmdLine) server.Reply
	AfterClientClose(client server.Connection)
	Close()
}

// KeyEventCallback 是一个回调接口，将在 key事件发生后被调用, 例如key的插入或者删除
type KeyEventCallback func(dbIndex int, key string, entity *DataEntity)

// DataEntity 存储key对应的数据：string,list,hash,set,zset 等等
type DataEntity struct {
	Data interface{}
}
