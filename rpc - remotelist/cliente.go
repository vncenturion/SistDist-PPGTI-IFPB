package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
)

func exibeMenu() {
	fmt.Println("Menu de opções:")
	fmt.Println("1. Adicionar valor em uma lista")
	fmt.Println("2. Remover último valor de uma lista")
	fmt.Println("3. Listar todas as listas e seus valores")
	fmt.Println("4. Exibir o tamanho de uma lista")
	fmt.Println("5. Retornar indice de um item da lista")
	fmt.Println("0. Sair")
	fmt.Print("Opção: ")
}

func main() {
	// Conexão com o servidor RPC
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Erro ao conectar ao servidor: ", err)
	}
	defer client.Close()

	// Menu de opções
	scanner := bufio.NewScanner(os.Stdin)
	for {

		exibeMenu()

		// Lê a opção escolhida pelo usuário
		scanner.Scan()
		option := scanner.Text()

		switch option {
		case "1":
			// Lê o nome da lista e o valor a ser adicionado
			fmt.Print("Digite o nome da lista: ")
			scanner.Scan()
			listName := scanner.Text()

			fmt.Print("Digite o valor a ser adicionado: ")
			scanner.Scan()
			value := scanner.Text()

			// Chama o método Append do servidor para adicionar o valor na lista
			var reply bool
			err = client.Call("Server.Append", []string{listName, value}, &reply)
			if err != nil {
				fmt.Println("Erro ao adicionar valor na lista: ", err)
			} else {
				fmt.Println("Valor adicionado com sucesso na lista!")
			}

		case "2":
			// Lê o nome da lista da qual o último valor será removido
			fmt.Print("Digite o nome da lista: ")
			scanner.Scan()
			listName := scanner.Text()

			// Chama o método Pop do servidor para remover o último valor da lista
			var reply int
			err = client.Call("Server.Pop", []string{listName}, &reply)
			if err != nil {
				fmt.Println("Erro ao remover último valor da lista: ", err)
			} else {
				fmt.Printf("Último valor removido da lista %s: %d\n", listName, reply)
			}

		case "3":
			// Chama o método ListAllLists do servidor para listar todas as listas e seus valores
			var reply map[string][]int
			err = client.Call("Server.ListAllLists", []string{}, &reply)
			if err != nil {
				fmt.Println("Erro ao listar as listas: ", err)
			} else {
				fmt.Println("Listas e seus valores:")
				for name, values := range reply {
					fmt.Printf("%s: %v\n", name, values)
				}
			}

		case "4":
			// Lê o nome da lista da qual se deseja obter o tamanho
			fmt.Print("Digite o nome da lista: ")
			scanner.Scan()
			listName := scanner.Text()

			// Chama o método Length do servidor para obter o tamanho da lista
			var reply int
			err = client.Call("Server.Size", listName, &reply)
			if err != nil {
				fmt.Println("Erro ao obter tamanho da lista: ", err)
			} else {
				fmt.Printf("Tamanho da lista %s: %d\n", listName, reply)
			}

		case "5":
			// Lê o nome da lista e o valor a ser adicionado
			fmt.Print("Digite o nome da lista: ")
			scanner.Scan()
			listName := scanner.Text()

			fmt.Print("Digite o valor a ser adicionado: ")
			scanner.Scan()
			value := scanner.Text()

			// Chama o método Append do servidor para adicionar o valor na lista
			var index int
			client.Call("Server.GetIndex", []string{listName, value}, &index)
			if err != nil {
				fmt.Println("Erro ao buscar item na lista: ", err)
			} else {
				fmt.Printf("Index do item buscado na lista %s: %d\n", listName, index)
			}

		case "0":
			// Encerra o programa
			fmt.Println("Saindo...")
			return

		default:
			fmt.Println("Opção inválida!")
		}
		fmt.Println("\n")
		fmt.Println(strings.Repeat("-", 20))
	}
}
