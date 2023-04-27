package main

import (
	"fmt"
	"net"
	"net/rpc"
	remotelist "ppgti/remotelist/pkg"
)

func main() {
	//list := new(remotelist.RemoteList)
	//Map := new(remotelist.RemoteMap)
	//Map := make(map[string]remotelist.RemoteList)
	//Rrpc := new(remotelist.RemoteMapRPC)

	/*Map, err := remotelist.loadMapFromFile("dados.txt")
	if err != nil || len(Map.Map) == 0 {
		Map := &remotelist.RemoteMap{Map: make(map[string]*remotelist.RemoteList)}
	}*/
	//fmt.Println("Carregando listas...\n")
	Map := &remotelist.RemoteMap{Map: make(map[string]*remotelist.RemoteList)}
	err := remotelist.LoadMapFromFile("dados.txt", Map)

	if err != nil {
		fmt.Println("Erro ao carregar listas:", err)
		// tratar erro
	}
	rpcs := rpc.NewServer()
	rpcs.Register(Map)
	l, e := net.Listen("tcp", "[localhost]:5000")
	defer l.Close()
	if e != nil {
		fmt.Println("listen error:", e)
	}
	//fmt.Println("Pronto esperando requisicoes...\n")
	for {
		conn, err := l.Accept()
		if err == nil {
			go rpcs.ServeConn(conn)
		} else {
			break
		}
	}
}
