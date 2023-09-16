package common

type Server interface {
	Initialize(options *Options) (err error)

	Run()
}

type ServerNode struct {
	Addr             string
}

type Options struct {
	// 服务器启动地址
	Addr string
	// 配置中心地址
	ConfigCenterAddr string
	// 租约有效时间
	LeaseTime int64
}
