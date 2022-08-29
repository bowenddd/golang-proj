package geeCache

import (
	"fmt"
	"geeCache/consistenthash"
	geecachepb "geeCache/geecachepb/proto"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 30
)

type HTTPPool struct {
	self        string // message of the node. like addr and port
	basePath    string
	mu          sync.Mutex
	peers       *consistenthash.Map
	httpGetters map[string]*httpGetter // each remote peer has a httpGetters
}

func NewHTTPPool(self string, basePath interface{}) *HTTPPool {
	path := defaultBasePath
	if v, ok := basePath.(string); ok {
		path = v
	}
	return &HTTPPool{
		self:     self,
		basePath: path,
	}
}

func (p *HTTPPool) Logger(format string, v ...interface{}) {
	log.Printf("[Server %s] %s \n", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !strings.HasPrefix(req.URL.Path, p.basePath) {
		panic("HTTPPooL serving unexpected path " + req.URL.Path)
	}
	p.Logger("%s %s", req.Method, req.URL.Path)
	//   /<basePath>/<group>/<key>
	parts := strings.SplitN(req.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
	}
	groupName := parts[0]
	key := parts[1]
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}
	byteView, err := group.Get(key)
	log.Printf("%v\n", byteView.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := proto.Marshal(&geecachepb.Response{Value: byteView.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}

func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}

}

func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		return p.httpGetters[peer], true
	}
	return nil, false

}

// a http client struct, it implements the PeerGetter interface
type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(in *geecachepb.Request, out *geecachepb.Response) error {
	u := fmt.Sprintf("%v%v/%v",
		h.baseURL,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
	)
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %v", resp.Status)
	}
	value, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}
	if err = proto.Unmarshal(value, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}
	return nil
}
