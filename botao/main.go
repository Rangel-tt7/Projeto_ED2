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

var entrada *bufio.Scanner //Variável global, aponta para o scanner de entrada
var saida *bufio.Writer    //Escritor da saída, envia no Flush

var largura, altura, turnos int
var hospitais []Hospital

// palavra lê o próximo campo da entrada. Devolve nil quando a entrada acaba
func palavra() []byte {
	if !entrada.Scan() { //Avança para o próximo pedaço de texto
		return nil
	}
	return entrada.Bytes() //Retorna o pedaço de texto lido
}

// inteiro lê o próximo field e converte para número
func inteiro() int {
	v := 0
	for _, c := range palavra() { //percorre o slice e entrega a posição e o item, descarta a posição
		v = v*10 + int(c-'0') // c-'0' converte o caractere no dígito (ex.: '7'-'0' = 7); v*10 desloca os dígitos já lidos uma casa para a esquerda
	}
	return v
}

// lerInicio lê a mensagem de início e guarda o mapa, os turnos e os hospitais
func lerInicio() {
	palavra() // INICIO
	largura = inteiro()
	altura = inteiro()
	turnos = inteiro()
	for i := 0; i < 4; i++ {
		inteiro() // limites, equipe e número de equipes (não usados)
	}

	palavra() // HOSPITAIS
	h := inteiro()
	for i := 0; i < h; i++ {
		x := inteiro()
		y := inteiro()
		hospitais = append(hospitais, Hospital{x, y})
	}

	palavra() // AMBULANCIAS
	k := inteiro()
	for i := 0; i < k*3; i++ {
		inteiro() // id x y (não usados ainda)
	}

	palavra() // BLOQUEIOS
	b := inteiro()
	for i := 0; i < b*4; i++ {
		inteiro() // x1 y1 x2 y2 (só na etapa 2)
	}

	palavra() // FIM
}

func main() {
	entrada = bufio.NewScanner(os.Stdin)
	entrada.Split(bufio.ScanWords)
	saida = bufio.NewWriter(os.Stdout)

	lerInicio()
	fmt.Fprintln(os.Stderr, "mapa:", largura, "x", altura, "| turnos:", turnos, "| hospitais:", hospitais)

	saida.WriteString("PRONTO\n")
	saida.Flush()
}
