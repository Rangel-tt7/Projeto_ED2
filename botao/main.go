package main //todo o arquivo gera um programa

import (
	"bufio" //Leitura rápida
	"fmt"   //Printf
	"os"    //Entrada e saída
)

// Hospital guarda a posição de um hospital
type Hospital struct {
	X, Y int
}

var leitor *bufio.Scanner  //Variável global, aponta para o scanner de entrada
var escritor *bufio.Writer //Escritor da saída, envia no Flush

var largura, altura, totalTurnos int
var hospitais []Hospital
var turnoAtual int

// lerPalavra lê a próxima palavra da entrada. Devolve nil quando a entrada acaba
func lerPalavra() []byte {
	if !leitor.Scan() { //Avança para o próximo pedaço de texto
		return nil
	}
	return leitor.Bytes() //Retorna o pedaço de texto lido
}

// lerNumero lê a próxima palavra e converte para número
func lerNumero() int {
	numero := 0
	for _, letra := range lerPalavra() { //percorre o slice e entrega a posição e o item, descarta a posição
		numero = numero*10 + int(letra-'0') // letra-'0' converte o caractere no dígito (ex.: '7'-'0' = 7); numero*10 desloca os dígitos já lidos uma casa para a esquerda
	}
	return numero
}

// lerInicio lê a mensagem de início e guarda o mapa, os turnos e os hospitais
func lerInicio() {
	lerPalavra() // INICIO
	largura = lerNumero()
	altura = lerNumero()
	totalTurnos = lerNumero()
	for i := 0; i < 4; i++ {
		lerNumero() // limites, equipe e número de equipes (não usados)
	}

	lerPalavra() // HOSPITAIS
	qtdHospitais := lerNumero()
	for i := 0; i < qtdHospitais; i++ {
		x := lerNumero()
		y := lerNumero()
		hospitais = append(hospitais, Hospital{x, y})
	}

	lerPalavra() // AMBULANCIAS
	qtdAmbulancias := lerNumero()
	for i := 0; i < qtdAmbulancias*3; i++ {
		lerNumero() // id x y (não usados ainda)
	}

	lerPalavra() // BLOQUEIOS
	qtdBloqueios := lerNumero()
	for i := 0; i < qtdBloqueios*4; i++ {
		lerNumero() // x1 y1 x2 y2 (só na etapa 2)
	}

	lerPalavra() // FIM
}

// lerTurno lê a mensagem de um turno. Devolve false quando a partida acaba
func lerTurno() bool {
	if lerPalavra() == nil { // TURNO, ou nil se o servidor fechou a entrada
		return false
	}
	turnoAtual = lerNumero()

	lerPalavra() // CHAMADOS
	qtdChamados := lerNumero()
	for i := 0; i < qtdChamados*5; i++ {
		lerNumero() // id x y valor prazo (guardados no passo 2)
	}

	lerPalavra() // RECOLHIDOS
	qtdRecolhidos := lerNumero()
	for i := 0; i < qtdRecolhidos; i++ {
		lerNumero() // id
	}

	lerPalavra() // AMBULANCIAS
	qtdAmbulancias := lerNumero()
	for i := 0; i < qtdAmbulancias; i++ {
		lerNumero()  // id
		lerNumero()  // x
		lerNumero()  // y
		lerPalavra() // LIVRE ou OCUPADA
	}

	lerPalavra() // PLACAR
	qtdEquipes := lerNumero()
	for i := 0; i < qtdEquipes*2; i++ {
		lerNumero() // equipe e pontos
	}

	lerPalavra() // FIM
	return true
}

func main() {
	leitor = bufio.NewScanner(os.Stdin)
	leitor.Split(bufio.ScanWords)
	escritor = bufio.NewWriter(os.Stdout)

	lerInicio()
	fmt.Fprintln(os.Stderr, "mapa:", largura, "x", altura, "| turnos:", totalTurnos, "| hospitais:", hospitais)

	escritor.WriteString("PRONTO\n")
	escritor.Flush()

	for lerTurno() {
		escritor.WriteString("FIM\n") // nenhuma ambulância se mexe ainda
		escritor.Flush()
	}
}
