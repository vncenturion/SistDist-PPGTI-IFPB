package main

import (
    "fmt"
    "log"
    "net"
    "net/rpc"
    "strconv"
    "sync"
)

// Definição da estrutura de lista
type List struct {
    Name   string
    Values []int
    Lock   sync.Mutex // Mutex para garantir exclusão mútua durante operações na lista
}

// Mapa que irá armazenar as listas criadas pelos clientes
var lists = make(map[string]*List)

// Estrutura para o servidor
type Server struct{}

// Método para adicionar um valor em uma lista existente ou criar uma nova lista caso não exista
func (s *Server) Append(args []string, reply *bool) error {
    listName := args[0]
    value, err := strconv.Atoi(args[1])
    if err != nil {
        return err
    }

    // Verifica se a lista já existe
    list, ok := lists[listName]
    if !ok {
        // Se a lista não existe, cria uma nova lista com o nome fornecido
        list = &List{Name: listName, Values: make([]int, 0)}
        lists[listName] = list
    }

    // Garante que apenas um cliente por vez esteja realizando operações na lista
    list.Lock.Lock()
    defer list.Lock.Unlock()

    // Adiciona o valor à lista
    list.Values = append(list.Values, value)

    *reply = true
    return nil
}

// Método para remover o último valor de uma lista
func (s *Server) Pop(args []string, reply *int) error {
    listName := args[0]

    // Verifica se a lista já existe
    list, ok := lists[listName]
    if !ok {
        return fmt.Errorf("Lista %s não encontrada", listName)
    }

    // Garante que apenas um cliente por vez esteja realizando operações na lista
    list.Lock.Lock()
    defer list.Lock.Unlock()

    // Remove o último valor da lista, se houver
    if len(list.Values) > 0 {
        lastIdx := len(list.Values) - 1
        *reply = list.Values[lastIdx]
        list.Values = list.Values[:lastIdx]
    } else {
        // Se a lista estiver vazia, retorna erro
        return fmt.Errorf("A lista %s está vazia", listName)
    }

    return nil
}

func main() {
    // Criação do servidor RPC
    server := new(Server)
    rpc.Register(server)

    // Configuração do listener do servidor
    listener, err := net.Listen("tcp", ":1234")
    if err != nil {
        log.Fatal("Erro ao iniciar servidor: ", err)
    }

    // Aceita e atende conexões dos clientes
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        go rpc.ServeConn(conn)
    }
}
