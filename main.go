package main

import "fmt"

//type Foo int
//
//type Args struct {
//	Num1, Num2 int
//}
//
//func (f Foo) Sum(args Args, reply *int) error {
//	*reply = args.Num1 + args.Num2
//	return nil
//}
//
//func startServer(addr chan string) {
//	var foo Foo
//	if err := geeRpc.Register(&foo); err != nil {
//		log.Fatal("register error:", err)
//	}
//	l, err := net.Listen("tcp", ":0")
//	if err != nil {
//		log.Fatal("network error:", err)
//	}
//	log.Println("start rpc server on", l.Addr())
//	addr <- l.Addr().String()
//	geeRpc.Accept(l)
//}
//func main() {
//	log.SetFlags(0)
//	addr := make(chan string)
//	go startServer(addr)
//	client, _ := geeRpc.Dial("tcp", <-addr, geeRpc.DefaultOption)
//	time.Sleep(time.Second)
//	var wg sync.WaitGroup
//	for i := 0; i < 5; i++ {
//		wg.Add(1)
//		go func(i int) {
//			defer wg.Done()
//			args := Args{Num1: i, Num2: i * i}
//			var reply int
//			if err := client.Call("Foo.Sum", args, &reply); err != nil {
//				log.Fatal("call Foo.Sum error:", err)
//			}
//			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
//		}(i)
//	}
//	wg.Wait()
//}

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

func main() {
	sli := make([]interface{}, 0)
	x := []interface{}{1, 2, 3}
	sli = append(sli, x)
	fmt.Println(sli)
}
