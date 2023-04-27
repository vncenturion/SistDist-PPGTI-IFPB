package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
)

func PrintMenu() {
	fmt.Println(`
Selecione uma opção:

1: Adicionar item
2: Remover ultimo item
3: Pesquisar item
4: Exibir tamanho
0: Encerrar

`)
}

func Menu() {
	reader := bufio.NewReader(os.Stdin)

	for {

		PrintMenu()

		fmt.Print("OPÇÃO: ")

		opcao, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		switch opcao {
		case "1\n":
			Append()
		case "2\n":
			Remove()
		case "3\n":
			Get()
		case "4\n":
			Size()
		case "0\n":
			fmt.Println("Encerrando...")
			os.Exit(0)
		default:
			fmt.Println("Opção inválida!")
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

	fmt.Print("Lista: ")
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

	fmt.Print("Lista: ")
	fmt.Scanln(&list_id)

	err = client.Call("RemoteMap.Size", list_id, &reply)
	if err != nil {
		fmt.Print("Error Size:", err)
	} else {
		fmt.Println("Tamanho:", list_id, ":", reply)
	}
}

func Get() {
	client, err := rpc.Dial("tcp", ":5000")
	if err != nil {
		fmt.Print("dialing:", err)
	}

	var list_id string
	fmt.Print("Lista: ")
	fmt.Scanln(&list_id)

	var i int
	fmt.Print("Posicao: ")
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
	fmt.Print("Lista: ")
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
		fmt.Println("OK!")
	}
}

func main() {

	Menu()

}
