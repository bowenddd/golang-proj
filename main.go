package main

import (
	"fmt"
	"geeRpc"
	"log"
	"net"
	"sync"
	"time"
)

func startServer(addr chan string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	geeRpc.Accept(l)
}
func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)
	client, _ := geeRpc.Dial("tcp", <-addr, geeRpc.DefaultOption)
	time.Sleep(time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("geerpc req %d", i)
			var reply string
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Println("reply:", reply)
		}(i)
	}
	wg.Wait()
}

//package main
//
//import (
//	"fmt"
//	"geeorm"
//	_ "github.com/mattn/go-sqlite3"
//)
//
//func main() {
//	engine, _ := geeorm.NewEngine("sqlite3", "gee.db")
//	defer engine.Close()
//	s := engine.NewSession()
//	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
//	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
//	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
//	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
//	count, _ := result.RowsAffected()
//	fmt.Printf("Exec success, %d affected\n", count)
//}
