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

func handleRespCmd(conn net.Conn, cmd Resp) error {
	switch v := cmd.(type) {
	case RespSimpleString:
		return handleSimpleCmd(conn, v)
	default:
		return fmt.Errorf("invalid type")
	}
}

func handleSimpleCmd(conn net.Conn, cmd RespSimpleString) error {
	switch cmd {
	case "PING":
		response := RespSimpleString("PONG").Build()
		n, err := conn.Write(response)
		if err != nil {
			return err
		}
		if n != len(response) {
			return fmt.Errorf("didn't write enough data: %d", n)
		}
		return nil

	default:
		return fmt.Errorf("unknown command:%s", cmd)

	}

}
