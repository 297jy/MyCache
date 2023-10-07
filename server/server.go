package server

import (
	"context"
	"gomemory/cluster"
	"gomemory/config"
	"gomemory/interface/database"
	"gomemory/lib/logger"
	"gomemory/lib/sync/atomic"
	"gomemory/server/connection"
	"gomemory/server/parser"
	"gomemory/server/protocol"
	"net"
	"sync"
)

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

// Handler implements tcp.Handler and serves as a redis server
type Handler struct {
	activeConn sync.Map // *client -> placeholder
	db         database.DB
	closing    atomic.Boolean // 当服务器正在停止服务时，值为true，拒绝所有客户端的请求
}

func MakeCacheHandler() *Handler {
	var db database.DB
	if config.CacheProperties.ClusterEnable {
		db = cluster.MakeCluster()
	} else {
		//db = database2.NewStandaloneServer()
	}
	return &Handler{
		db: db,
	}
}

// Handle receives and executes redis commands
func (h *Handler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		_ = conn.Close()
		return
	}

	client := connection.NewConn(conn)
	h.activeConn.Store(client, struct{}{})

	ch := parser.ParseStream(conn)
	for payload := range ch {
		if payload.Data == nil {
			logger.Error("empty payload")
			continue
		}
		r, ok := payload.Data.(*protocol.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk protocol")
			continue
		}
		result := h.db.Exec(client, r.Args)
		if result != nil {
			_, _ = client.Write(result.ToBytes())
		} else {
			_, _ = client.Write(unknownErrReplyBytes)
		}
	}
}

// Close stops handler
func (h *Handler) Close() error {
	logger.Info("handler shutting down...")
	h.closing.Set(true)
	h.activeConn.Range(func(key, value any) bool {
		client := key.(*connection.Connection)
		_ = client.Close()
		return true
	})
	h.db.Close()
	return nil
}
