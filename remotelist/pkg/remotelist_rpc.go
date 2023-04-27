package remotelist

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type RemoteList struct {
	//mu   sync.Mutex
	list []int
	size uint32
}

type RemoteMap struct {
	Map map[string]*RemoteList
	mu  sync.Mutex
}

//type RemoteMapRPC int

//Map := new(remotelist.RemoteMap)
//Map := make(map[string]remotelist.RemoteList)

//rpcs := rpc.NewServer()
//rpcs.Register(Map)

/*func (l *RemoteList) Get(i int, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.list) > 0 {
		*reply = l.list[i]
		//fmt.Println(l.list)
	} else {
		return errors.New("empty list")
	}
	return nil
}

func (l *RemoteList) Size(arg int, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.list) > 0 {
		*reply = len(l.list)
	} else {
		return errors.New("empty list")
	}
	return nil
}
*/

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

	// Codifica o novo mapa em JSON e salva no arquivo
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

	// Decodifica o conteúdo do arquivo em JSON para um mapa
	mapData := make(map[string][]int)
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&mapData); err != nil {
		return err
	}

	// Converte o mapa decodificado em um mapa de RemoteList
	for key, list := range mapData {
		m.Map[key] = &RemoteList{
			list: list,
			size: uint32(len(list)),
		}
	}

	return nil
}

func (m *RemoteMap) All(arg int, reply *bool) error {
	/*m.mu.Lock()
	defer m.mu.Unlock()
	m.Map[key].list = append(m.Map[key].list, value)
	fmt.Println(m.Map[key].list)
	m.Map[key].size++
	*reply = true
	return nil*/
	m.mu.Lock()
	defer m.mu.Unlock()

	fmt.Printf("*************************\n")
	for key, value := range m.Map {
		fmt.Printf("Nome: %s\n", key)
		fmt.Printf("Valores: %v\n", value.list)
		fmt.Printf("Tamanho: %d\n\n", len(value.list))
	}
	fmt.Printf("*************************\n")
	/*keys := make([]string, len(m.Map))
	i := 0
	for k := range m.Map {
		keys[i] = k
		i++
	}
	fmt.Println(keys)*/
	//fmt.Println(m.Map[key].list)

	*reply = true

	//SaveMapToFile("dados.txt", m)

	return nil
}

func (m *RemoteMap) Remove(list_id string, reply *int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var valor int
	//var l = &RemoteList{} //new(RemoteList) //&RemoteList{}

	l, ok := m.Map[list_id]
	if !ok {
		fmt.Println("Lista de nome: ", list_id, "nao existe!")
		return errors.New("Lista nao existe\n")
	} else {
		valor = l.list[len(l.list)-1]
		l.list = l.list[:len(l.list)-1]
		fmt.Println(l.list)
		fmt.Println("Valor ", valor, "da lista ", list_id, "foi removido com sucesso!")
		*reply = valor
		SaveMapToFile("dados.txt", m)
	}

	return nil
}

func (m *RemoteMap) Size(list_id string, reply *int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var tam int
	//var l = &RemoteList{} //new(RemoteList) //&RemoteList{}

	l, ok := m.Map[list_id]
	if !ok {
		fmt.Println("Lista de nome: ", list_id, "nao existe!")
		return errors.New("Lista nao existe\n")
	} else {
		tam = len(l.list)
		fmt.Println("Tamanho da lista: ", list_id, ":", tam)
	}

	*reply = tam

	return nil
}

func (m *RemoteMap) New(novalista string, reply *string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var l = &RemoteList{} //new(RemoteList) //&RemoteList{}
	m.Map[novalista] = l
	fmt.Println("Nova lista criada com sucesso: ", novalista)

	*reply = novalista
	SaveMapToFile("dados.txt", m)
	
	return nil
}

func (m *RemoteMap) Get(args *struct {
	List_id string
	I       int
}, reply *int) error {
	/*m.mu.Lock()
	defer m.mu.Unlock()
	m.Map[key].list = append(m.Map[key].list, value)
	fmt.Println(m.Map[key].list)
	m.Map[key].size++
	*reply = true
	return nil*/
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.Map[args.List_id]
	if !ok {
		fmt.Println("Lista de nome: ", args.List_id, "nao existe!")
		return errors.New("Lista nao existe\n")
	} else if len(l.list) > 0 && len(l.list) > args.I {
		*reply = l.list[args.I] //l.list = append(l.list, args.Value)
		fmt.Println("Buscando elemento na lista: ", args.List_id, " na posicao: ", args.I, ":")
		fmt.Println(l.list)
		fmt.Println("Valor na posicao:", args.I, "e:", *reply)

	} else {
		fmt.Println("Indece: ", args.I, "nao existe!")
		return errors.New("Indice nao existe\n")
	}

	//*reply = true
	return nil
}

func (m *RemoteMap) Append(args *struct {
	List_id string
	V       int
}, reply *bool) error {
	/*m.mu.Lock()
	defer m.mu.Unlock()
	m.Map[key].list = append(m.Map[key].list, value)
	fmt.Println(m.Map[key].list)
	m.Map[key].size++
	*reply = true
	return nil*/
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.Map[args.List_id]
	if !ok {
		l = &RemoteList{} //new(RemoteList) //&RemoteList{}
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

/*func (l *RemoteList) Remove(arg int, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.list) > 0 {
		*reply = l.list[len(l.list)-1]
		l.list = l.list[:len(l.list)-1]
		fmt.Println(l.list)
	} else {
		return errors.New("empty list")
	}
	return nil
}*/

func NewRemoteList() *RemoteList {
	return new(RemoteList)
}
