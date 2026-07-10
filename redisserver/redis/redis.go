package redis

import (
	"fmt"
	"net"
	"sync"
)

type AppEnv struct {
	Port   int
	Addr   string
	data   sync.Map
	config map[string]string
}

func (app *AppEnv) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", app.Addr, app.Port))
	if err != nil {
		return err
	}
	defer listener.Close()

	app.config = make(map[string]string)
	app.config["save"] = "3600 1 300 100 60 10000"
	app.config["appendonly"] = "no"

	fmt.Printf("listening for connections on %s:%d\n", app.Addr, app.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		go handleClient(app, conn)
	}
}

func handleClient(app *AppEnv, conn net.Conn) {
	defer conn.Close()

	parser := NewRespParser(conn, 1025)

	for {
		elem, err := parser.GetRespPrimitive()
		if err != nil {
			return
		}

		if err := handleRespCmd(app, conn, elem); err != nil {
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

func sendErrorResp(conn net.Conn, format string, args ...any) error {
	return sendResp(conn, RespError(fmt.Sprintf(format, args...)))
}

func handleRespCmd(app *AppEnv, conn net.Conn, cmd Resp) error {
	switch v := cmd.(type) {
	case RespArray:
		return handleArrayCmd(app, conn, v)
	default:
		return fmt.Errorf("invalid type")
	}
}

func handleArrayCmd(app *AppEnv, conn net.Conn, cmd RespArray) error {
	switch v := cmd[0].(type) {
	case RespBulkString:
		switch v {
		case "INFO":
			return sendResp(conn, RespBulkString(`# Server
			redis_version:0.0.0.0
			redis_mode:standalone
			`))
		case "PING":
			return sendResp(conn, RespSimpleString("pong"))
		case "ECHO":
			return handleEcho(conn, cmd[1:])
		case "SET":
			return handleSet(app, conn, cmd[1:])
		case "GET":
			return handleGet(app, conn, cmd[1:])
		case "CONFIG":
			return handleConfig(app, conn, cmd[1:])
		default:
			return fmt.Errorf("unknown array command:%s", v)
		}
	default:
		return fmt.Errorf("unknown first elem %T", cmd[0])
	}

}

func handleEcho(conn net.Conn, args RespArray) error {
	if len(args) != 1 {
		return sendErrorResp(conn, "wrong number of arguments for 'echo' command")
	}
	return sendResp(conn, args[0])
}

func handleSet(app *AppEnv, conn net.Conn, args RespArray) error {
	if len(args) != 2 {
		return sendErrorResp(conn, "wrong number of arguments for 'set' command")
	}
	key := args[0]
	value := args[1]
	app.data.Store(key.String(), value.String())

	return sendResp(conn, RespSimpleString("OK"))

}

func handleGet(app *AppEnv, conn net.Conn, args RespArray) error {
	if len(args) != 1 {
		return sendErrorResp(conn, "wrong number of arguments for 'get' command")
	}
	key := args[0]
	value, ok := app.data.Load(key.String())
	if !ok {
		return sendResp(conn, RespNull{})
	} else {
		switch v := value.(type) {
		case string:
			return sendResp(conn, RespBulkString(v))
		default:
			return fmt.Errorf("non-string type value")
		}

	}

}

func handleConfig(app *AppEnv, conn net.Conn, args RespArray) error {
	if len(args) < 1 {
		return sendErrorResp(conn, "wrong number of arguments for 'config' command")
	}
	switch v := args[0].(type) {
	case RespBulkString:
		switch v {
		case "GET":
			return handleConfigGet(app, conn, args[1:])
		case "SET":
			return handleConfigSet(app, conn, args[1:])
		default:
			return fmt.Errorf("unknown config subcommand %s", v)
		}

	default:
		return fmt.Errorf("unknown config subcommand type %T", v)
	}
}

func handleConfigGet(app *AppEnv, conn net.Conn, keys RespArray) error {
	if len(keys) < 1 {
		return sendErrorResp(conn, "wrong number of arguments for 'config|get' command")
	}

	var responseArray RespArray

	for _, key := range keys {
		keyStr := key.String()
		value, ok := app.config[keyStr]

		if ok {
			responseArray = append(responseArray, RespBulkString(keyStr))
			responseArray = append(responseArray, RespBulkString(value))
		}
	}

	return sendResp(conn, responseArray)
}

func handleConfigSet(app *AppEnv, conn net.Conn, key_values RespArray) error {
	if len(key_values) < 2 || len(key_values)%2 != 0 {
		return sendResp(conn, RespError("wrong number of arguments for 'config|set' command"))
	}
	//TODO protect with mutex
	for idx := 0; idx < len(key_values); idx += 2 {
		app.config[key_values[idx].String()] = key_values[idx+1].String()
	}
	return sendResp(conn, RespSimpleString("OK"))
}
