package codec

import "io"

type Header struct {
	ServiceMethod string // format "Service.Method"
	Seq           uint64 // this is to distinguish different client
	Err           string
}

type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

type NewCodecFunc func(closer io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json" //not implement
)

var NewCodecFuncMap map[Type]NewCodecFunc

// init()函数会自动执行，进行一些初始化操作
// init()函数的执行顺序  import -> const -> var -> init()
func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	// 工厂模式，不过这里返回的是构造方法而不是实例
	NewCodecFuncMap[GobType] = NewGobCodec
}
