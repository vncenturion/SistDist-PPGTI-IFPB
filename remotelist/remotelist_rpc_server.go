package main

import (
	"fmt"
	"net"
	"net/rpc"
	remotelist "ppgti/remotelist/pkg"
)

func main() {
	
	Map := &remotelist.RemoteMap{Map: make(map[string]*remotelist.RemoteList)}
	err := remotelist.LoadMapFromFile("dados.txt", Map)

	if err != nil {
		fmt.Println("Erro ao carregar listas:", err)
	}

	rpcs := rpc.NewServer()
	rpcs.Register(Map)
	l, e := net.Listen("tcp", "[localhost]:5000")
	defer l.Close()
	if e != nil {
		fmt.Println("listen error:", e)
	}

	for {
		conn, err := l.Accept()
		if err == nil {
			go rpcs.ServeConn(conn)
		} else {
			break
		}
	}
}
