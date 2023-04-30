# Prática - RPC
 
Implemente um sistema distribuído que consiste em um componente denominado de RemoteList, que gerencia um conjunto de listas de valores inteiros, e um conjunto de clientes, que utilizam o serviço oferecido pelo RemoteList (ou seja, inserção, consulta e remoção de elementos em listas). Dessa forma, o RemoteList funciona como um servidor e armazena dados submetidos por clientes em listas (mais de uma lista pode existir e cada lista deve possuir um identicador único). Também permite que clientes remotos consultem dados da lista, de qualquer posição, bem como removam dados da lista, mas nesse caso apenas o valor do nal da lista. Vários clientes podem utilizar os serviços oferecidos por RemoteList em simultâneo e usar listas em comum.

Deve-se implementar um esquema de comunicação síncrono (o cliente obtém conrmação da operação realizada) e persistente (os dados continuam armazenados mesmo se os clientes ou o servidor pararem de executar), utilizando Remote Procedure Call (RPC). Dessa forma, as seguintes operações precisam ser disponibilizadas para os clientes via RPC:

1. Append(list_id, v) → coloca o valor v no nal da lista com identicador list_id
1. Get(list_id, i) → retorna o valor da posição i, na lista com identicador list_id
1. Remove(list_id) → remove e retorna o último elemento da lista com identicador list_id
1. Size(list_id) → obtém a quantidade de elementos armazenados na lista com identicador list_id

É permitido utilizar o código a seguir como base, sem restrições: https://github.com/ruandg/SD_PPGTI/tree/main/remotelist