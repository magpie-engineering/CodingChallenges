package redis

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type RespParser struct {
	conn net.Conn
	buf  []byte // Single unified buffer
	r    int    // Read index (where the parser is currently reading)
	w    int    // Write index (where the network is currently writing)
}

func NewRespParser(conn net.Conn, buf_size int) *RespParser {
	return &RespParser{
		conn: conn,
		buf:  make([]byte, buf_size),
	}
}

func (parser *RespParser) GetN(n int) ([]byte, error) {
	// Does the buffer data ready to read
	// keep reading from the network until it does
	for parser.w-parser.r < n {
		//is there data unread in the middle of buffer
		if parser.r > 0 {
			// slide unread data back to the start of the buffer
			copy(parser.buf, parser.buf[parser.r:parser.w])
			parser.w -= parser.r
			parser.r = 0
		}

		// Is there enough space left in the buffer or do we need to grow it
		if n > len(parser.buf) {
			newBuf := make([]byte, n*2) // Double the size for headroom
			copy(newBuf, parser.buf[:parser.w])
			parser.buf = newBuf
		}

		// Read new data into the available space at the tail
		nRead, err := parser.conn.Read(parser.buf[parser.w:])
		if err != nil {

			return nil, err
		}
		parser.w += nRead
	}

	// 2. Extract the requested bytes
	out := make([]byte, n)
	copy(out, parser.buf[parser.r:parser.r+n])
	parser.r += n
	return out, nil
}

func (parser *RespParser) GetRespPrimitive() (Resp, error) {
	respType, err := parser.GetN(1)
	if err != nil {
		return nil, err
	}
	switch respType[0] {
	case '+':
		return parser.parseSimpleString()
	case '-':
		return parser.parseError()
	case ':':
		return parser.parseInteger()
	case '$':
		return parser.parseBulkString()
	case '*':
		return parser.parseArray()
	case '_':
		return parser.parseNull()
	case '#':
		return parser.parseBool()
	default:
		return nil, fmt.Errorf("Unknown command type:%c", respType[0])
	}

}

func (parser *RespParser) parseRespString() (string, error) {
	var s strings.Builder
	var c []byte
	var err error

	for c, err = parser.GetN(1); !(c[0] == '\r' || err != nil); c, err = parser.GetN(1) {
		s.WriteByte(c[0])
	}
	if err != nil {
		return "", err
	}
	c, err = parser.GetN(1)
	if err != nil {
		return "", err
	}
	if c[0] != '\n' {
		return "", fmt.Errorf("Didn't get correct line ending for simple string:%c", c[0])
	}
	return s.String(), nil
}

func (parser *RespParser) parseLineEnding() error {
	line_ending, err := parser.GetN(2)
	if err != nil {
		return err
	} else if string(line_ending) != "\r\n" {
		return fmt.Errorf("missing line ending in bulk string:%s", line_ending)
	}
	return nil
}

func (parser *RespParser) parseSimpleString() (RespSimpleString, error) {
	s, err := parser.parseRespString()
	if err != nil {
		return "", err
	}
	return RespSimpleString(s), err
}

func (parser *RespParser) parseError() (RespError, error) {
	s, err := parser.parseRespString()
	if err != nil {
		return "", err
	}
	return RespError(s), err
}

func (parser *RespParser) parseNull() (RespNull, error) {
	err := parser.parseLineEnding()
	if err != nil {
		return struct{}{}, err
	}
	return struct{}{}, nil
}

func (parser *RespParser) parseInteger() (RespInteger, error) {
	s, err := parser.parseRespString()
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return RespInteger(n), err

}

func (parser *RespParser) parseBulkString() (RespBulkString, error) {
	n, err := parser.parseInteger()
	if err != nil {
		return "", err
	}
	if n < 0 {
		return "", nil
	}
	bytes, err := parser.GetN(int(n))
	if err != nil {
		return "", err
	}
	err = parser.parseLineEnding()
	if err != nil {
		return "", err
	}

	return RespBulkString(bytes), nil

}

func (parser *RespParser) parseArray() (RespArray, error) {
	n, err := parser.parseInteger()
	if err != nil {
		return nil, err
	}
	if n < 0 {
		return nil, nil
	}
	array := make([]Resp, int(n))
	for idx := range array {
		elem, err := parser.GetRespPrimitive()
		if err != nil {
			return nil, err
		}
		array[idx] = elem
	}
	return array, nil

}

func (parser *RespParser) parseBool() (RespBool, error) {

	b, err := parser.GetN(1)
	if err != nil {
		return false, err
	}
	var val bool
	switch b[0] {
	case 't':
		val = true
	case 'f':
		val = false
	default:
		return false, fmt.Errorf("unknown boolean value:%s", b)
	}

	err = parser.parseLineEnding()
	if err != nil {
		return false, err
	}

	return RespBool(val), nil

}
