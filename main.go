package main

import (
	"context"
	"geeRpc"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type Foo int

type Args struct {
	Num1, Num2 int
}

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func startServer(addr chan string) {
	var foo Foo
	if err := geeRpc.Register(&foo); err != nil {
		log.Fatal("register error:", err)
	}
	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	geeRpc.HandleHTTP()
	addr <- l.Addr().String()
	_ = http.Serve(l, nil)
}

func call(addrCh chan string) {
	client, _ := geeRpc.DialHTTP("tcp", <-addrCh)
	defer func() {
		_ = client.Close()
	}()
	time.Sleep(time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := Args{Num1: i, Num2: i * i}
			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			var reply int
			if err := client.Call(ctx, "Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()

}

func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go call(addr)
	startServer(addr)
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
