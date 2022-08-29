package geeCache

import geecachepb "geeCache/geecachepb/proto"

// this interface is used to choose a peer that store a special key

type PeerPicker interface {
	// PickPeer 根据一个建通过一致性哈希算法选择存储该键的值的节点，返回一个实现了从该节点获取值的接口
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// a peer must implement this interface, and PeerGetter is used for a client

type PeerGetter interface {

	// Get this method is used for searching a key from a special group
	Get(in *geecachepb.Request, out *geecachepb.Response) error
}
