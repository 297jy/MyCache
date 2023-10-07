package server

// Connection 表示来自客户端的连接
type Connection interface {
	Write([]byte) (int, error)

	Close() error
}
