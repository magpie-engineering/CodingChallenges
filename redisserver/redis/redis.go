package redis

import (
	"fmt"
	"net"
)

type AppEnv struct {
	Port int
	Addr string
}

func (app *AppEnv) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", app.Addr, app.Port))
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Printf("listening for connections on %s:%d\n", app.Addr, app.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	parser := NewRespParser(conn, 1025)

	for {
		elem, err := parser.GetRespPrimitive()
		if err != nil {
			return
		}
		fmt.Println("Received:", elem)
		if err := handleRespCmd(conn, elem); err != nil {
			fmt.Println("error:", err)
		}
	}
}

func sendResp(conn net.Conn, resp Resp) error {
	response := resp.Build()
	n, err := conn.Write(response)
	if err != nil {
		return err
	}
	if n != len(response) {
		return fmt.Errorf("didn't write enough data: %d", n)
	}
	return nil
}

func handleRespCmd(conn net.Conn, cmd Resp) error {
	switch v := cmd.(type) {
	case RespArray:
		return handleArrayCmd(conn, v)
	default:
		return fmt.Errorf("invalid type")
	}
}

func handleArrayCmd(conn net.Conn, cmd RespArray) error {
	switch v := cmd[0].(type) {
	case RespBulkString:
		switch v {
		case "info":
			return sendResp(conn, RespBulkString(`# Server
			redis_version:0.0.0.0
			redis_mode:standalone
			`))
		case "ping":
			return sendResp(conn, RespSimpleString("pong"))
		case "echo":
			return handleEcho(conn, cmd[1:])
		default:
			return fmt.Errorf("unknown array command:%s", v)
		}
	default:
		return fmt.Errorf("unknown first elem %T", cmd[0])
	}

}

func handleEcho(conn net.Conn, args RespArray) error {
	if len(args) != 1 {
		return fmt.Errorf("incorrect number of args to echo:%d", len(args))
	}
	return sendResp(conn, args[0])
}
