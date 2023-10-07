package cluster

import (
	"errors"
	"gomemory/datastruct/dict"
	"gomemory/lib/logger"
	"gomemory/lib/pool"
	"gomemory/server/client"
)

type defaultClientFactory struct {
	nodeConnections dict.Dict // map[string]*pool.Pool
}

var connectionPoolConfig = pool.Config{
	MaxIdle:   1,
	MaxActive: 16,
}

func newDefaultClientFactory() *defaultClientFactory {
	return &defaultClientFactory{
		nodeConnections: dict.MakeConcurrent(1),
	}
}

// GetPeerClient gets a client with peer form pool
func (factory *defaultClientFactory) GetPeerClient(peerAddr string) (peerClient, error) {
	var connectionPool *pool.Pool
	raw, ok := factory.nodeConnections.Get(peerAddr)
	if !ok {
		creator := func() (interface{}, error) {
			c, err := client.MakeClient(peerAddr)
			if err != nil {
				return nil, err
			}
			c.Start()
			return c, nil
		}
		finalizer := func(x interface{}) {
			logger.Debug("destroy client")
			cli, ok := x.(client.Client)
			if !ok {
				return
			}
			cli.Close()
		}
		connectionPool = pool.New(creator, finalizer, connectionPoolConfig)
		factory.nodeConnections.Put(peerAddr, connectionPool)
	} else {
		connectionPool = raw.(*pool.Pool)
	}
	raw, err := connectionPool.Get()
	if err != nil {
		return nil, err
	}
	conn, ok := raw.(*client.Client)
	if !ok {
		return nil, errors.New("connection pool make wrong type")
	}
	return conn, nil
}

// ReturnPeerClient returns client to pool
func (factory *defaultClientFactory) ReturnPeerClient(peer string, peerClient peerClient) error {
	return nil
}

func (factory *defaultClientFactory) NewStream(peerAddr string, cmdLine CmdLine) (peerStream, error) {
	return nil, nil
}

func (factory *defaultClientFactory) Close() error {
	return nil
}
