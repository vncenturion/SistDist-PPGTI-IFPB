package remotelist

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type RemoteList struct {
	list []int
	size uint32
}

type RemoteMap struct {
	Map map[string]*RemoteList
	mu  sync.Mutex
}


func (m *RemoteMap) Append(args *struct {
	List_id string
	V       int
}, reply *bool) error {
	
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.Map[args.List_id]
	if !ok {
		l = &RemoteList{}
		m.Map[args.List_id] = l
	}
	fmt.Println("Exibindo lista de nome:", args.List_id)
	l.list = append(l.list, args.V)
	fmt.Println(l.list)
	l.size++
	*reply = true
	SaveMapToFile("dados.txt", m)
	return nil
}


func (m *RemoteMap) Remove(list_id string, reply *int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var valor int

	l, ok := m.Map[list_id]
	if !ok {
		fmt.Println("Lista: ", list_id, " inexistente!")
		return errors.New("Lista inexistente.\n")
	} else {
		valor = l.list[len(l.list)-1]
		l.list = l.list[:len(l.list)-1]
		fmt.Println(l.list)
		fmt.Println("Valor ", valor, "em ", list_id, "removido!")
		*reply = valor
		SaveMapToFile("dados.txt", m)
	}
	return nil
}

func (m *RemoteMap) Get(args *struct {
	List_id string
	I       int
}, reply *int) error {
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	l, ok := m.Map[args.List_id]
	if !ok {
		return fmt.Errorf("Lista de nome %s inexistente!", args.List_id)
	}
	
	if args.I >= len(l.list) {
		return fmt.Errorf("Indice %d inexistente!", args.I)
	}
	
	*reply = l.list[args.I]
	fmt.Printf("Buscando elemento na lista: %s na posicao: %d\n", args.List_id, args.I)
	fmt.Println(l.list)
	fmt.Printf("Valor na posicao %d é: %d\n", args.I, *reply)

	return nil
}

func (m *RemoteMap) Size(listID string, reply *int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	l, ok := m.Map[listID]
	if !ok {
		return fmt.Errorf("Lista de nome %s inexistente!", listID)
	}

	tam := len(l.list)
	fmt.Printf("Tamanho da lista %s: %d\n", listID, tam)

	*reply = tam

	return nil
}

func NewRemoteList() *RemoteList {
	return new(RemoteList)
}

func SaveMapToFile(filename string, m *RemoteMap) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	mapData := make(map[string][]int)
	for key, list := range m.Map {
		mapData[key] = list.list
	}

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&mapData); err != nil {
		return err
	}

	return nil
}

func LoadMapFromFile(filename string, m *RemoteMap) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	mapData := make(map[string][]int)
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&mapData); err != nil {
		return err
	}

	for key, list := range mapData {
		m.Map[key] = &RemoteList{
			list: list,
			size: uint32(len(list)),
		}
	}

	return nil
}
