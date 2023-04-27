package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
)

func New() {
	var reply string
	var novalista string

	client, err := rpc.Dial("tcp", ":5000")
	if err != nil {
		fmt.Print("dialing:", err)
	}
	fmt.Print("Nome da nova lista: ")
	fmt.Scanln(&novalista)
	err = client.Call("RemoteMap.New", novalista, &reply)
	if err != nil {
		fmt.Print("Error New:", err)
	} else {
		fmt.Println("Nova lista criada com sucesso! Nome:", reply)
	}
}

func Menu() {
	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Println("************************************************")
		fmt.Println("Selecione uma opção:")
		fmt.Println("************************************************")
		fmt.Println("1 - ALL - \t Exibir todas as listas")
		fmt.Println("2 - NEW - \t Criar nova lista")
		fmt.Println("3 - GET - \t Pesquisar valor em lista")
		fmt.Println("4 - APPEND - \t Adicionar valor em lista existente ou nova")
		fmt.Println("5 - SIZE - \t Exibir tamanho da lista")
		fmt.Println("6 - REMOVE - \t Remover valor(ultimo) em lista")
		fmt.Println("0 - Sair")
		fmt.Println("************************************************")
		fmt.Print("Resposta: ")

		opcao, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		switch opcao {
		case "1\n":
			fmt.Println("1 - ALL - Exibir todas as listas")
			All()

		case "2\n":
			fmt.Println("2 - NEW - Criar nova lista: ")
			New()

		case "3\n":
			fmt.Println("3 - GET - Pesquisar valor em lista: ")
			Get()

		case "0\n":
			fmt.Println("Saindo do programa...")
			os.Exit(0)
		case "4\n":
			fmt.Println("4 - APPEND - Adicionar valor em lista:")
			Append()
		case "5\n":
			fmt.Println("5 - SIZE - Exibir tamanho da lista:")
			Size()
		case "6\n":
			fmt.Println("6 - REMOVE - Remover valor(ultimo) em lista:")
			Remove()
		default:
			fmt.Println("Opção inválida, tente novamente.")
		}
	}
	return
}

func Remove() {
	client, err := rpc.Dial("tcp", ":5000")
	if err != nil {
		fmt.Print("dialing:", err)
	}

	var reply int
	var list_id string

	fmt.Print("Nome da lista: ")
	fmt.Scanln(&list_id)

	err = client.Call("RemoteMap.Remove", list_id, &reply)
	if err != nil {
		fmt.Print("Error Remove:", err)
	} else {
		fmt.Println("Ultima posicao da lista", list_id, "de valor", reply, "foi removida com sucesso!")
	}
}

func Size() {
	client, err := rpc.Dial("tcp", ":5000")
	if err != nil {
		fmt.Print("dialing:", err)
	}

	var reply int
	var list_id string

	fmt.Print("Nome da lista: ")
	fmt.Scanln(&list_id)

	err = client.Call("RemoteMap.Size", list_id, &reply)
	if err != nil {
		fmt.Print("Error Size:", err)
	} else {
		fmt.Println("Tamanho da lista:", list_id, ":", reply)
	}
}

func Get() {
	client, err := rpc.Dial("tcp", ":5000")
	if err != nil {
		fmt.Print("dialing:", err)
	}

	var list_id string
	fmt.Print("Nome da lista: ")
	fmt.Scanln(&list_id)

	var i int
	fmt.Print("Posicao na lista: ")
	fmt.Scanln(&i)

	var reply int

	err = client.Call("RemoteMap.Get", &struct {
		List_id string
		I       int
	}{list_id, i}, &reply)
	if err != nil {
		fmt.Print("Error Get:", err)
	} else {
		fmt.Println("Valor na posicao:", i, "e:", reply)
	}
}

func Append() {
	client, err := rpc.Dial("tcp", ":5000")
	if err != nil {
		fmt.Print("dialing:", err)
	}

	var list_id string
	fmt.Print("Nome da lista(se nao exitir, sera criado uma nova): ")
	fmt.Scanln(&list_id)

	var v int
	fmt.Print("Valor: ")
	fmt.Scanln(&v)

	var reply bool

	err = client.Call("RemoteMap.Append", &struct {
		List_id string
		V       int
	}{list_id, v}, &reply)
	if err != nil && reply == true {
		fmt.Print("Error Append:", err)
	} else {
		fmt.Println("Valores inseridos com sucesso!")
	}
}

func All() {
	client, err := rpc.Dial("tcp", ":5000")
	if err != nil {
		fmt.Print("dialing:", err)
	}

	var reply bool
	err = client.Call("RemoteMap.All", 0, &reply)
	if err != nil {
		fmt.Print("Error All:", err)
	}
}

func main() {
	/*client, err := rpc.Dial("tcp", ":5000")
	if err != nil {
		fmt.Print("dialing:", err)
	}*/

	Menu()

	// Synchronous call
	/*var reply bool
	//var reply_i int
	//err = client.Call("RemoteMap.Menu", 0, &reply)
	//fmt.Print("Error:", err)

	err = client.Call("RemoteMap.Append", &struct {
		Key   string
		Value int
	}{"marcio", 10}, &reply)
	fmt.Print("Error:", err)*/
	/*err = client.Call("RemoteList.Append", &struct {
		Key   string
		Value int
	}{"marcio", 20}, &reply)
	err = client.Call("RemoteList.Append", &struct {
		Key   string
		Value int
	}{"marcio", 30}, &reply)
	err = client.Call("RemoteList.Append", &struct {
		Key   string
		Value int
	}{"marcio", 40}, &reply)
	err = client.Call("RemoteList.Append", &struct {
		Key   string
		Value int
	}{"marcio", 50}, &reply)*/
	//fmt.Print("Error:", err)
	/*fmt.Println("Teste Get:")
	err = client.Call("RemoteList.Get", 2, &reply_i)
	if err != nil {
		fmt.Print("Error:", err)
	} else {
		fmt.Println("Elemento consultado:", reply_i)
	}

	fmt.Println("Teste Size:")
	err = client.Call("RemoteList.Size", 0, &reply_i)
	if err != nil {
		fmt.Print("Error:", err)
	} else {
		fmt.Println("Tamanho da Lista Consultada:", reply_i)
	}

	err = client.Call("RemoteList.Remove", 0, &reply_i)
	if err != nil {
		fmt.Print("Error:", err)
	} else {
		fmt.Println("Elemento retirado:", reply_i)
	}
	err = client.Call("RemoteList.Remove", 0, &reply_i)
	if err != nil {
		fmt.Print("Error:", err)
	} else {
		fmt.Println("Elemento retirado:", reply_i)
	}*/
}
