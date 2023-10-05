package cluster

import "gomemory/datastruct/dict"

type defaultClientFactory struct {
	nodeConnections dict.Dict // map[string]*pool.Pool
}

func newDefaultClientFactory() *defaultClientFactory {
	return &defaultClientFactory{
		nodeConnections: dict.MakeConcurrent(1),
	}
}

// GetPeerClient gets a client with peer form pool
func (factory *defaultClientFactory) GetPeerClient(peerAddr string) (peerClient, error) {
	return nil, nil
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
