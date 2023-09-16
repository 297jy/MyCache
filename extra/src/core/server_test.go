package core

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestZookeeper(t *testing.T) {

	s := NewServer("test")

	// 创建zk连接地址
	hosts := []string{"192.168.136.119:2181"}
	// 连接zk
	localConn, _, err := zk.Connect(hosts, time.Second*5)
	conn = localConn
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	s.Online()

	addr := "localhost:9999"
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, s))
}

func TestStr(t *testing.T) {
	x := []int{1,2,3}
	fmt.Printf("%d_%d_%d", 1, 2, 3)
	fmt.Println(x[0:1])
}

func getPoint() *Server{
	return &Server{
		Name: "test",
	}
}
