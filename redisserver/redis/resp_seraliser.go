package redis

import (
	"bytes"
	"fmt"
	"strconv"
)

func (s RespSimpleString) Build() []byte {
	return wrapResp('+', []byte(s))
}

func (s RespError) Build() []byte {
	return wrapResp('-', []byte(s))
}

func (i RespInteger) Build() []byte {
	s := strconv.Itoa(int(i))
	return wrapResp(':', []byte(s))
}

func (s RespBulkString) Build() []byte {
	var buf bytes.Buffer
	len_string := strconv.Itoa(len(s))
	buf.Grow(1 + len(len_string) + 2 + len(s) + 2)

	buf.WriteByte('$')
	buf.Write([]byte(len_string))
	buf.WriteString("\r\n")
	buf.Write([]byte(s))
	buf.WriteString("\r\n")

	return buf.Bytes()

}

func (s RespArray) Build() []byte {
	var buf bytes.Buffer
	len_string := strconv.Itoa(len(s))

	var elems_len int
	for _, elem := range s {
		elems_len += elem.SerialLen()
	}

	buf.Grow(1 + len(len_string) + 2 + elems_len)

	buf.WriteByte('*')
	buf.Write([]byte(len_string))
	buf.WriteString("\r\n")

	for _, elem := range s {
		buf.Write(elem.Build())
	}

	return buf.Bytes()

}

func (s RespNull) Build() []byte {
	return []byte("_\r\n")
}

func (s RespBool) Build() []byte {
	var bool_char byte
	if bool(s) {
		bool_char = 't'
	} else {
		bool_char = 'f'
	}
	return []byte(fmt.Sprintf("#%c\r\n", bool_char))
}

func (s RespMap) Build() []byte {
	var buf bytes.Buffer
	len_string := strconv.Itoa(len(s))

	buf.Grow(s.SerialLen())

	buf.WriteByte('%')
	buf.Write([]byte(len_string))
	buf.WriteString("\r\n")

	for key, value := range s {
		buf.Write(key.Build())
		buf.Write(value.Build())
	}

	return buf.Bytes()

}

func wrapResp(prefixChar byte, original []byte) []byte {
	var buf bytes.Buffer
	buf.Grow(1 + len(original) + 2)

	buf.WriteByte(prefixChar)
	buf.Write(original)
	buf.WriteString("\r\n")

	return buf.Bytes()
}

func (s RespSimpleString) SerialLen() int {
	return 1 + len(s) + 2
}

func (s RespError) SerialLen() int {
	return 1 + len(s) + 2
}

func (i RespInteger) SerialLen() int {
	s := strconv.Itoa(int(i))
	return 1 + len(s) + 2
}

func (s RespBulkString) SerialLen() int {
	len_string := strconv.Itoa(len(s))
	return 1 + len(len_string) + 2 + len(s) + 2

}

func (s RespArray) SerialLen() int {

	len_string := strconv.Itoa(len(s))

	var elems_len int
	for _, elem := range s {
		elems_len += elem.SerialLen()
	}
	return 1 + len(len_string) + 2 + elems_len
}

func (s RespNull) SerialLen() int {
	return 3
}

func (s RespBool) SerialLen() int {
	return 3
}

func (s RespMap) SerialLen() int {

	len_string := strconv.Itoa(len(s))

	var elems_len int
	for key, value := range s {
		elems_len += key.SerialLen()
		elems_len += value.SerialLen()
	}
	return 1 + len(len_string) + 2 + elems_len
}
