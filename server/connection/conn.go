package connection

import (
	"gomemory/lib/logger"
	"gomemory/lib/sync/wait"
	"net"
	"sync"
	"time"
)

type Connection struct {
	conn net.Conn

	sendingData wait.Wait

	mu sync.Mutex

	flag uint64
}

var connPool = sync.Pool{
	New: func() interface{} {
		return &Connection{}
	},
}

// NewConn 通过池化技术，获取一个连接实例，避免重复创建对象
func NewConn(conn net.Conn) *Connection {
	c, ok := connPool.Get().(*Connection)
	if !ok {
		logger.Error("connection pool make wrong type")
		return &Connection{
			conn: conn,
		}
	}
	c.conn = conn
	return c
}

func (c *Connection) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	c.sendingData.Add(1)
	defer func() {
		c.sendingData.Done()
	}()
	return c.conn.Write(b)
}

func (c *Connection) Close() error {
	c.sendingData.WaitWithTimeout(10 * time.Second)
	_ = c.conn.Close()
	// 归还连接池
	connPool.Put(c)
	return nil
}
