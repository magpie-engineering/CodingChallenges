package redis

import (
	"strconv"
	"strings"
)

type Resp interface {
	Build() []byte
	SerialLen() int
	String() string
}
type RespSimpleString string
type RespError string
type RespInteger int
type RespBulkString string
type RespArray []Resp
type RespNull struct{}
type RespBool bool
type RespMap map[Resp]Resp

func (s RespSimpleString) String() string {
	return string(s)
}

func (s RespError) String() string {
	return string(s)
}

func (s RespInteger) String() string {
	return strconv.Itoa(int(s))
}

func (s RespBulkString) String() string {
	return string(s)
}

func (s RespArray) String() string {
	var b strings.Builder
	b.WriteRune('[')
	for _, elem := range s {
		b.WriteString(elem.String())
		b.WriteString(", ")
	}
	b.WriteRune(']')
	return b.String()
}

func (s RespNull) String() string {
	return "Null"
}

func (s RespBool) String() string {
	return strconv.FormatBool(bool(s))
}

func (s RespMap) String() string {
	var b strings.Builder
	b.WriteRune('{')
	for key, value := range s {
		b.WriteString(key.String())
		b.WriteString(": ")
		b.WriteString(value.String())
		b.WriteString(", ")
	}
	b.WriteRune('}')
	return b.String()
}
