package main

import (
	"context"
	"fmt"
	"geeRpc"
	"geeRpc/registry"
	"geeRpc/xclient"
	"log"
	"net"
	"net/http"
	"reflect"
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

func (f Foo) Sleep(args Args, reply *int) error {
	time.Sleep(time.Second * time.Duration(args.Num1))
	*reply = args.Num1 + args.Num2
	return nil
}

func startRegistry(wg *sync.WaitGroup) {
	l, _ := net.Listen("tcp", ":9999")
	registry.HandleHTTP()
	wg.Done()
	_ = http.Serve(l, nil)
}

func startServer(registryAddr string, wg *sync.WaitGroup) {
	var foo Foo
	l, _ := net.Listen("tcp", ":0")
	server := geeRpc.NewServer()
	_ = server.Register(&foo)
	registry.Heartbeat(registryAddr, "tcp@"+l.Addr().String(), 0)
	wg.Done()
	server.Accept(l)
}

func foo(xc *xclient.XClient, ctx context.Context, typ, serviceMethod string, args *Args) {
	var reply int
	var err error
	switch typ {
	case "call":
		err = xc.Call(ctx, serviceMethod, args, &reply)
	case "broadcast":
		err = xc.Broadcast(ctx, serviceMethod, args, &reply)
	}
	if err != nil {
		log.Printf("%s %s error: %v", typ, serviceMethod, err)
	} else {
		log.Printf("%s %s success: %d + %d = %d", typ, serviceMethod, args.Num1, args.Num2, reply)
	}
}

func call(registryAddr string) {
	d := xclient.NewGeeRegistryDiscovery(registryAddr, 0)
	xc := xclient.NewXClient(d, xclient.RandomSelect, nil)
	defer func() {
		_ = xc.Close()
	}()
	time.Sleep(time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := Args{Num1: i, Num2: i * i}
			foo(xc, context.Background(), "call", "Foo.Sum", &args)
		}(i)
	}
	wg.Wait()
}

func broadcast(registryAddr string) {
	d := xclient.NewGeeRegistryDiscovery(registryAddr, 0)
	xc := xclient.NewXClient(d, xclient.RandomSelect, nil)
	defer func() { _ = xc.Close() }()
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			foo(xc, context.Background(), "broadcast", "Foo.Sum", &Args{Num1: i, Num2: i * i})
			// expect 2 - 5 timeout
			ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
			foo(xc, ctx, "broadcast", "Foo.Sleep", &Args{Num1: i, Num2: i * i})
		}(i)
	}
	wg.Wait()
}

//func main() {
//	log.SetFlags(0)
//	registryAddr := "http://localhost:9999/_geerpc_/registry"
//	var wg sync.WaitGroup
//	wg.Add(1)
//	go startRegistry(&wg)
//	wg.Wait()
//
//	time.Sleep(time.Second)
//	wg.Add(2)
//	go startServer(registryAddr, &wg)
//	go startServer(registryAddr, &wg)
//	wg.Wait()
//
//	time.Sleep(time.Second)
//	//call(registryAddr)
//	broadcast(registryAddr)
//}

type User struct {
	name string
	age  int
	addr string
}

func main() {
	uu := User{"tom", 27, "beijing"}
	u := &uu
	v := reflect.ValueOf(u).Interface()

	fmt.Println("ValueOf=", reflect.ValueOf(v).Elem())

	t := reflect.TypeOf(v)
	fmt.Println("TypeOf=", t)

	t1 := reflect.Indirect(reflect.ValueOf(v)).Type()
	fmt.Println("t1=", t1)

	t2 := reflect.TypeOf(v).Elem()
	fmt.Println("t2=", t2)
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
