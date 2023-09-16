package core

import (
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"net/http"
	"strings"
	"time"
)

const RegisterPath = "/easy-memcache-service"

const defaultBasePath = "/_easycache/"

type Service interface {
	discover() error

	register() error

	Online() error

	WatchCluster()
}

const (
	Running int = iota
	Stopped
)

type Server struct {
	Ip       string
	Name     string
	status   int
	replicas int
}

type LocalServer struct {
}

var Provides []Server

var conn *zk.Conn

func NewServer(name string) *Server {
	s := &Server{
		status: Running,
		Ip:     fmt.Sprintf("%d", time.Now().Unix()),
		Name:   name,
	}
	return s
}

func (s *Server) String() string {
	return fmt.Sprintf("[ip => %v, host => %v]", s.Ip, s.Name)
}

func (s *Server) Equal(p *Server) bool {
	return s.Name == p.Name
}

func (s *Server) Log(format string, v ...interface{}) {
	log.Printf("Server %s[%s], %s", s.Ip, s.Name, fmt.Sprintf(format, v...))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, defaultBasePath) {
		panic("Server unexpected path: " + r.URL.Path)
	}
	s.Log("%s %s", r.Method, r.URL.Path)

	parts := strings.SplitN(r.URL.Path[len(defaultBasePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	fmt.Println(groupName + key)
}

func (s *Server) WatchCluster() {
	for ; s.status == Running; {
		_, _, e, err := conn.ChildrenW(RegisterPath)
		if err != nil {
			s.Log("ChildrenW error:{%s}", err)
			continue
		}
		event := <-e
		fmt.Println("###########################")
		fmt.Println("path: ", event.Path)
		fmt.Println("type: ", event.Type.String())
		fmt.Println("state: ", event.State.String())
		fmt.Println("---------------------------")
		// 只有子节点变化才重新开始服务发现
		if event.Type == zk.EventNodeChildrenChanged {
			_ = s.discover()
		}
	}
}

func (s *Server) Online() error {
	go s.WatchCluster()
	_ = s.register()
	return nil
}

func (s *Server) discover() error {
	Provides = make([]Server, 0)
	names, _, err := conn.Children(RegisterPath)
	if err != nil {
		return err
	}
	for _, name := range names {
		data, _, gerr := conn.Get(RegisterPath + "/" + name)
		if gerr != nil {
			fmt.Println(gerr)
			continue
		}
		var server Server
		uerr := json.Unmarshal(data, &server)
		if uerr != nil {
			s.Log("service discovery json error, %s", uerr)
			continue
		}
		Provides = append(Provides, server)
	}
	return nil
}

func (s *Server) register() error {
	if e, _, err := conn.Exists(RegisterPath); !e || err != nil {
		_, cerr := conn.Create(RegisterPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if cerr != nil {
			return cerr
		}
	}
	data, merr := json.Marshal(s)
	if merr != nil {
		return merr
	}
	_, cerr := conn.Create(RegisterPath+"/"+s.Name, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if cerr != nil {
		fmt.Println(cerr)
		return cerr
	}
	return nil
}
