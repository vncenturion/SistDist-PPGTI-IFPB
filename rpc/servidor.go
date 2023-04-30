package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

// Método para salvar a lista em um arquivo
func (l *List) SaveList(filename string) error {
	// Cria um novo arquivo e escreve a lista
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, v := range l.Values {
		_, err = file.WriteString(strconv.Itoa(v) + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

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

	// Salva a lista em um arquivo
	err = list.SaveList(listName + ".txt")
	if err != nil {
		log.Printf("Erro ao saver a lista list %s no arquivo: %s\n", listName, err)
	}

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

		// Chama o método SaveList para salvar a lista em um arquivo
		err := list.SaveList(listName + ".txt")
		if err != nil {
			log.Printf("Erro ao salvar a lista %s: %v", listName, err)
		}

	} else {
		// Se a lista estiver vazia, retorna erro
		return fmt.Errorf("A lista %s está vazia", listName)
	}

	return nil
}

// Método para listar todas as listas e seus valores
func (s *Server) ListAllLists(args []string, reply *map[string][]int) error {
	*reply = make(map[string][]int)
	for name, list := range lists {
		(*reply)[name] = list.Values
	}
	return nil
}

// Método para obter o tamanho de uma lista
func (s *Server) Size(listName string, reply *int) error {
	// Verifica se a lista existe
	list, ok := lists[listName]
	if !ok {
		// Se a lista não existe, retorna um erro indicando que a lista não foi encontrada
		return fmt.Errorf("Lista não encontrada")
	}

	// Garante que apenas um cliente por vez esteja realizando operações na lista
	list.Lock.Lock()
	defer list.Lock.Unlock()

	// Retorna o tamanho da lista
	*reply = len(list.Values)
	return nil
}

// Método para retornar o valor de um determinado índice de uma lista
func (s *Server) Get(args []string, reply *int) error {
	listName := args[0]
	index, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}

	// Verifica se a lista já existe
	list, ok := lists[listName]
	if !ok {
		return fmt.Errorf("Lista %s não encontrada", listName)
	}

	// Garante que apenas um cliente por vez esteja realizando operações na lista
	list.Lock.Lock()
	defer list.Lock.Unlock()

	// Verifica se o índice está dentro do range da lista
	if index < 0 || index >= len(list.Values) {
		return fmt.Errorf("Índice inválido para a lista %s: %d", listName, index)
	}

	// Retorna o valor do índice
	*reply = list.Values[index]
	return nil
}

// Método para retornar o valor de index de um item da lista
func (s *Server) GetIndex(args []string, reply *int) error {
	listName := args[0]
	value, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}

	// Verifica se a lista já existe
	list, ok := lists[listName]
	if !ok {
		return fmt.Errorf("Lista %s não encontrada", listName)
	}

	// Garante que apenas um cliente por vez esteja realizando operações na lista
	list.Lock.Lock()
	defer list.Lock.Unlock()

	// Encontra o index do valor na lista
	for i, v := range list.Values {
		if v == value {
			*reply = i
			return nil
		}
	}

	// Se o valor não for encontrado, retorna erro
	return fmt.Errorf("Valor %d não encontrado na lista %s", value, listName)
}

// Método para carregar todas as listas salvas na pasta atual
func LoadLists() error {
	// Abre a pasta atual
	dir, err := os.Open(".")
	if err != nil {
		return err
	}
	defer dir.Close()

	// Lista todos os arquivos na pasta atual
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	// Itera sobre todos os arquivos na pasta atual
	for _, fileInfo := range fileInfos {
		// Verifica se o arquivo é do tipo .txt
		if !fileInfo.IsDir() && filepath.Ext(fileInfo.Name()) == ".txt" {
			// Lê o arquivo e carrega a lista correspondente no mapa lists
			file, err := os.Open(fileInfo.Name())
			if err != nil {
				return err
			}
			defer file.Close()

			// Cria uma nova lista com o nome do arquivo (sem a extensão .txt)
			listName := strings.TrimSuffix(fileInfo.Name(), ".txt")
			list := &List{Name: listName, Values: make([]int, 0)}
			lists[listName] = list

			// Lê todos os valores do arquivo e adiciona à lista
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				value, err := strconv.Atoi(scanner.Text())
				if err != nil {
					return err
				}
				list.Values = append(list.Values, value)
			}
			if err := scanner.Err(); err != nil {
				return err
			}

			fmt.Printf("Lista carregada: %s \n", listName)
		}
	}

	return nil
}

func main() {
	// Criação do servidor RPC
	server := new(Server)
	err := rpc.Register(server)
	if err != nil {
		return
	}
	fmt.Println("Servidor inicializado")

	// Carrega todas as listas salvas na pasta atual
	err = LoadLists()
	if err != nil {
		log.Fatal("Erro ao carregar as listas: ", err)
	}

	// Configuração do listener do servidor
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Erro ao iniciar servidor: ", err)
	}
	fmt.Println("Aguardando conexões")

	// Aceita e atende conexões dos clientes
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(conn)
	}
}
