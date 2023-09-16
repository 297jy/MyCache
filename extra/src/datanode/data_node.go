package main

import (
	"EasyMemcache/src/common"
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type DataNode interface {
	Register() (err error)
	Close() (err error)
}

type dataServer struct {
	common.ServerNode
	configCenterAddr string
	meta             registerMeta
}

type registerMeta struct {
	cli           *clientv3.Client //etcd client
	leaseID       clientv3.LeaseID //租约ID
	leaseKey      string
	leaseTime     int64 // 租约续期时间
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
}

func (s *dataServer) Initialize(os *common.Options) (err error) {

	return err
}

func (s *dataServer) Run() {
	s.Register()
	//监听续租相应chan，保证服务不会过期
	go s.listenLeaseRespChan()
}

func (s *dataServer) Register() (err error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{s.configCenterAddr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	s.meta.cli = cli
	if err := s.putKeyWithLease(); err != nil {
		return err
	}
	return err
}

func (s *dataServer) putKeyWithLease() error {
	//设置租约时间
	resp, err := s.meta.cli.Grant(context.Background(), s.meta.leaseTime)
	if err != nil {
		return err
	}

	registerVal, err := s.renderRegisterVal()
	if err != nil {
		return err
	}
	//注册服务并绑定租约
	_, err = s.meta.cli.Put(context.Background(), s.meta.leaseKey, registerVal, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	//设置续租 定期发送需求请求
	leaseRespChan, err := s.meta.cli.KeepAlive(context.Background(), resp.ID)

	if err != nil {
		return err
	}
	s.meta.leaseID = resp.ID
	fmt.Println(s.meta.leaseID)
	s.meta.keepAliveChan = leaseRespChan
	fmt.Printf("Put key:%s  val:%s  success!", s.meta.leaseKey, registerVal)
	return nil
}

func (s *dataServer) renderRegisterVal() (val string, err error) {
	marshal, err := json.Marshal(*s)
	return common.ByteToString(marshal), err
}

func (s *dataServer) listenLeaseRespChan() {
	for leaseKeepResp := range s.meta.keepAliveChan {
		fmt.Println("续约成功", leaseKeepResp)
	}
	fmt.Println("关闭续租")
}

// Close 注销服务
func (s *dataServer) Close() error {
	//撤销租约
	if _, err := s.meta.cli.Revoke(context.Background(), s.meta.leaseID); err != nil {
		return err
	}
	fmt.Println("撤销租约")
	return s.meta.cli.Close()
}
