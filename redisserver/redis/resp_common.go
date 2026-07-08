package redis

type Resp interface {
	Build() []byte
	SerialLen() int
}
type RespSimpleString string
type RespError string
type RespInteger int
type RespBulkString string
type RespArray []Resp
type RespNull struct{}
type RespBool bool
