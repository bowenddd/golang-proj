package main

import (
	"flag"
	"fmt"
	"geeCache"
	"net/http"
)

var db = map[string]string{
	"Tom":  "36",
	"Sam":  "28",
	"Jack": "22",
}

func createGroup() *geeCache.Group {
	return geeCache.NewGroup("scores", 2<<10, geeCache.GetterFunc(
		func(key string) ([]byte, error) {
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exists\n", key)
		}))
}

func startCacheServer(addr string, addrs []string, gee *geeCache.Group) {
	peers := geeCache.NewHTTPPool(addr, nil)
	peers.Set(addrs...)
	gee.RegisterPeer(peers)
}

func startAPIServer(apiAddr string, gee *geeCache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			key := request.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			writer.Header().Set("Content-Type", "application/octet-stream")
			writer.Write(view.ByteSlice())
		}))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], []string(addrs), gee)
}
