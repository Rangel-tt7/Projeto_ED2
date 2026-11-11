package main //todo o arquivo gera um programa

import (
	"bufio" //Leitura rápida
	"fmt"
	"os" //Entrada e saída
)

var leitor *bufio.Scanner  //Variável global, aponta para o scanner de entrada
var escritor *bufio.Writer //Escritor da saída, envia no Flush

var largura, altura, totalTurnos int
var hospitais []Hospital
var turnoAtual int
var chamados []Chamado       //todos os chamados, na ordem em que chegaram
var ambulancias []Ambulancia //todas as ambulâncias, aplica a regra do 0 até k-1
var porID *TabelaHash        //para achar o índice de um chamado pelo ID
var expira [][]int           //uma lista de páginas, cada uma com os IDs dos chamados que expiram naquele turno
var chamadosComuns *Grade    //grade para os chamados comuns, menos que valem até 20
var chamadosUrgentes *Grade  //grade para os chamados urgentes, mais de 20

func main() {
	leitor = bufio.NewScanner(os.Stdin)
	leitor.Split(bufio.ScanWords)
	escritor = bufio.NewWriter(os.Stdout)

	lerInicio()
	fmt.Fprintln(os.Stderr, "mapa:", largura, "x", altura, "| turnos:", totalTurnos, "| hospitais:", hospitais)

	escritor.WriteString("PRONTO\n")
	escritor.Flush()

	for lerTurno() {
		jogar()
	}
}
