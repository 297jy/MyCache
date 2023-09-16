package main

import (
	"EasyMemcache/src/core"
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"net/http"
	"time"
)

var (
	hosts       = []string{"192.168.136.119:2181"}
	path        = "/easy-memcache-service"
	flags int32 = zk.FlagEphemeral
	data        = []byte("zk data 001")
	acls        = zk.WorldACL(zk.PermAll)
)

func main() {
	/**
	core.LocalServer = core.Server{

	}
	// 创建监听的option，用于初始化zk
	eventCallbackOption := zk.WithEventCallback(callback)
	// 连接zk
	conn, _, err := zk.Connect(hosts, time.Second*5, eventCallbackOption)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 开始监听path
	_, _, _, err = conn.ExistsW(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 触发创建数据操作
	create(conn, path, data)

	//再次监听path
	_, _, _, err = conn.ExistsW(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 触发删除数据操作
	del(conn, path)
	**/
	s := core.NewServer("test")
	localConn, _, err := zk.Connect(hosts, time.Second*5)
	defer localConn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}


	go func() {
		for {
			data, merr := json.Marshal(s)
			if merr != nil {
				continue
			}

			s.Name = "test" + fmt.Sprintf("%d", time.Now().Unix())
			_, cerr := localConn.Create("/easy-memcache-service/"+s.Name, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
			if cerr != nil {
				fmt.Println(cerr)
				//return
			}
			fmt.Println("上线:" + s.Name)
			time.Sleep(time.Duration(3) * time.Second)
		}
	}()

	go watchZkEvent(localConn)

	addr := "localhost:9998"
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, s))
}

// zk watch 回调函数
func callback(event zk.Event) {
	// zk.EventNodeCreated
	// zk.EventNodeDeleted
	fmt.Println("#########监听到服务器列表变化#########")
	fmt.Println("path: ", event.Path)
	fmt.Println("type: ", event.Type.String())
	fmt.Println("state: ", event.State.String())
	fmt.Println("---------------------------")
}

func watchZkEvent(localConn *zk.Conn) {
	for {
		fmt.Println("fdsfsf")
		_, _, event1, err := localConn.ChildrenW(path)
		if err != nil {
			fmt.Println(err)
			continue
		}
		event := <-event1
		fmt.Println("###########################")
		fmt.Println("path: ", event.Path)
		fmt.Println("type: ", event.Type.String())
		fmt.Println("state: ", event.State.String())
		fmt.Println("---------------------------")
	}
}

// 创建数据
func create(conn *zk.Conn, path string, data []byte) {
	_, err := conn.Create(path, data, flags, acls)
	if err != nil {
		fmt.Printf("创建数据失败: %v\n", err)
		return
	}
	fmt.Println("创建数据成功")
}

// 删除数据
func del(conn *zk.Conn, path string) {
	_, stat, _ := conn.Get(path)
	err := conn.Delete(path, stat.Version)
	if err != nil {
		fmt.Printf("删除数据失败: %v\n", err)
		return
	}
	fmt.Println("删除数据成功")
}
