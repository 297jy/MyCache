package main

import (
	"EasyMemcache/src/common"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.etcd.io/etcd/clientv3"
	"net/http"
	"sync"
	"time"
)

type SessionNode interface {
	Listen(context *gin.Context)
	Publish(data string)
}

type sessionServer struct {
	common.ServerNode
	configCenterAddr string
	r                *gin.Engine
	ls               *listenerManager
}

type listenerManager struct {
	listeners []listener
	rwlock    sync.Mutex
}

type listener struct {
	timer   *time.Timer
	context *gin.Context
}

func (lm *listenerManager) trigger(data string) (err error) {
	lm.rwlock.Lock()
	defer lm.rwlock.Unlock()
	if lm.listeners == nil {
		return nil
	}
	for i := 0; i < len(lm.listeners); i += 1 {
		// 先停止阻塞
		lm.listeners[i].timer.Stop()
		lm.listeners[i].context.String(http.StatusOK, data)
	}

	lm.listeners = make([]listener, 10)
	return nil
}

func (lm *listenerManager) append(context *gin.Context, waitDuration time.Duration) (l listener, err error) {
	lm.rwlock.Lock()
	defer lm.rwlock.Unlock()
	l = listener{context: context, timer: time.NewTimer(waitDuration)}
	if lm.listeners == nil {
		lm.listeners = make([]listener, 2)
	}
	lm.listeners = append(lm.listeners, l)
	return l, err
}

func (l listener) wait() {
	for {
		<-l.timer.C
		// 阻塞时间到
	}
}

func (s *sessionServer) Initialize(os *common.Options) (err error) {
	s.Addr = os.Addr
	s.configCenterAddr = os.ConfigCenterAddr
	s.r = gin.Default()
	s.r.GET("/listen", s.Listen)
	return err
}

func (s *sessionServer) Run() {

	s.watchConfigCeneter()
	_ = s.r.Run(s.Addr)
}

func (s *sessionServer) watchConfigCeneter() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{s.configCenterAddr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()

	rch := cli.Watch(context.Background(), "easy-data-node/", clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			if ev.IsModify() || ev.IsCreate() {
				// 需要发送的数据
				s.Publish("")
			}
		}
	}
}

func (s *sessionServer) Listen(context *gin.Context) {
	if s.ls == nil {
		s.ls = new(listenerManager)
	}

	l, _ := s.ls.append(context, 30*time.Second)
	l.wait()
}

func (s *sessionServer) Publish(data string) {
	if s.ls == nil {
		return
	}

	s.ls.trigger(data)
}
