//package main
//
//import (
//	"encoding/json"
//	"fmt"
//	"geeRpc"
//	"geeRpc/codec"
//	"geeorm"
//	"log"
//	"net"
//	"time"
//)
//
//func startServer(addr chan string) {
//	l, err := net.Listen("tcp", ":0")
//	if err != nil {
//		log.Fatal("network error:", err)
//	}
//	log.Println("start rpc server on", l.Addr())
//	addr <- l.Addr().String()
//	geeRpc.Accept(l)
//}
//func main() {
//	addr := make(chan string)
//	go startServer(addr)
//	conn, _ := net.Dial("tcp", <-addr)
//	defer func() { _ = conn.Close() }()
//	time.Sleep(time.Second)
//	_ = json.NewEncoder(conn).Encode(geeRpc.DefaultOption)
//	cc := codec.NewGobCodec(conn)
//	for i := 0; i < 5; i++ {
//		h := &codec.Header{
//			ServiceMethod: "Foo.Sum",
//			Seq:           uint64(i),
//		}
//		_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
//		_ = cc.ReadHeader(h)
//		var reply string
//		_ = cc.ReadBody(&reply)
//		log.Println("reply:", reply)
//	}
//}

package main

import (
	"fmt"
	"geeorm"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	engine, _ := geeorm.NewEngine("sqlite3", "gee.db")
	defer engine.Close()
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}
