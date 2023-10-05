package server

// Reply 序列化后的 服务器响应消息体
type Reply interface {
	ToBytes() []byte
}
