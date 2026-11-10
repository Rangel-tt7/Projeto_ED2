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

type Chamado struct {
	ID, X, Y, Valor, Prazo int
	Aberto                 bool
}

type Ambulancia struct {
	ID, X, Y int
	Livre    bool
}

var leitor *bufio.Scanner  //Variável global, aponta para o scanner de entrada
var escritor *bufio.Writer //Escritor da saída, envia no Flush

var largura, altura, totalTurnos int
var hospitais []Hospital
var turnoAtual int
var chamados []Chamado       //todos os chamados, na ordem em que chegaram
var ambulancias []Ambulancia //todas as ambulâncias, aplica a regra do 0 até k-1
var porID *TabelaHash        //para achar o índice de um chamado pelo ID
var expira [][]int           //uma lista de páginas, cada uma com os IDs dos chamados que expiram naquele turno

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
	ambulancias = make([]Ambulancia, qtdAmbulancias)
	for i := 0; i < qtdAmbulancias; i++ {
		id := lerNumero()
		x := lerNumero()
		y := lerNumero()
		ambulancias[id] = Ambulancia{id, x, y, true}
	}

	lerPalavra() // BLOQUEIOS
	qtdBloqueios := lerNumero()
	for i := 0; i < qtdBloqueios*4; i++ {
		lerNumero() // x1 y1 alvoX alvoY (só na etapa 2)
	}

	lerPalavra()                                            // FIM
	porID = CriarTabela(proximoPrimo(largura * altura / 2)) //tamanho inicial da tabela hash
	chamados = make([]Chamado, 0, largura*altura/4)         //capacidade inicial do slice, make([]Tipo, tamanho, capacidade)
	expira = make([][]int, totalTurnos+2000)                //uma página para cada turno
}

// lerTurno lê a mensagem de um turno. Devolve false quando a partida acaba
func lerTurno() bool {
	if lerPalavra() == nil { // TURNO, ou nil se o servidor fechou a entrada
		return false
	}
	turnoAtual = lerNumero()

	for _, id := range expira[turnoAtual-1] { //percorre o slice e entrega a posição e o item, descarta a posição
		fechar(id) //fecha o chamado
	}
	expira[turnoAtual-1] = nil

	lerPalavra() // CHAMADOS
	qtdChamados := lerNumero()

	for i := 0; i < qtdChamados; i++ {
		id := lerNumero() // id do chamado
		x := lerNumero()
		y := lerNumero()
		valor := lerNumero()
		prazo := lerNumero()
		posicao := len(chamados)                                           //posição do chamado no slice
		chamados = append(chamados, Chamado{id, x, y, valor, prazo, true}) //adiciona o chamado ao slice
		porID.Inserir(id, posicao)                                         // guarda na hash
		expira[prazo] = append(expira[prazo], id)                          // anota o id na página do prazo
	}

	lerPalavra() // RECOLHIDOS
	qtdRecolhidos := lerNumero()
	for i := 0; i < qtdRecolhidos; i++ {
		id := lerNumero() // id
		fechar(id)        //fecha o chamado
	}

	lerPalavra() // AMBULANCIAS
	qtdAmbulancias := lerNumero()
	for i := 0; i < qtdAmbulancias; i++ {
		ambulanciaID := lerNumero() // id
		x := lerNumero()            // x
		y := lerNumero()            // y
		estado := lerPalavra()      // LIVRE ou OCUPADA
		if string(estado) == "LIVRE" {
			ambulancias[ambulanciaID].Livre = true
		} else {
			ambulancias[ambulanciaID].Livre = false
		}
		ambulancias[ambulanciaID].X = x
		ambulancias[ambulanciaID].Y = y
	}

	lerPalavra() // PLACAR
	qtdEquipes := lerNumero()
	for i := 0; i < qtdEquipes*2; i++ {
		lerNumero() // equipe e pontos
	}

	lerPalavra() // FIM
	return true
}

func direcao(x1, y1, alvoX, alvoY int) string { //Função que indica a direção que a ambulância tem que seguir
	if x1 < alvoX {
		return "L"
	} else if x1 > alvoX {
		return "O"
	} else if y1 < alvoY {
		return "S"
	} else if y1 > alvoY {
		return "N"
	}
	return "P"
}

func abs(x int) int { //Função que devolve o valor absoluto de um número
	if x < 0 {
		return -x
	}
	return x
}

func manhattan(x1, y1, x2, y2 int) int { //Função calcula a distância entre dois pontos
	return abs(x1-x2) + abs(y1-y2)
}

func hospitalMaisProximo(x, y int) (int, int) { //Encontra o hospital mais próximo de um ponto
	menorDistancia := largura + altura
	var hospitalX, hospitalY int
	for _, hospital := range hospitais {
		distancia := manhattan(x, y, hospital.X, hospital.Y)
		if distancia < menorDistancia {
			menorDistancia = distancia
			hospitalX = hospital.X
			hospitalY = hospital.Y
		}
	}
	return hospitalX, hospitalY
}

func fechar(id int) {
	busca, achou := porID.Buscar(id) //retorna o índice do chamado no slice e se achou
	if !achou {
		return // o chamado não está na lista
	}
	chamados[busca].Aberto = false //muda o chamado para fechado
	porID.Remover(id)
}

func jogar() {
	for _, ambulancia := range ambulancias { //percorre o slice, entrega posição e item
		movimento := "P"       // ambulância não se mexe
		if !ambulancia.Livre { //ambulância ocupada procura hospital mais próximo
			hospitalX, hospitalY := hospitalMaisProximo(ambulancia.X, ambulancia.Y)
			movimento = direcao(ambulancia.X, ambulancia.Y, hospitalX, hospitalY)
		}
		fmt.Fprintf(escritor, "%d %s\n", ambulancia.ID, movimento)
	}
	escritor.WriteString("FIM\n")
	escritor.Flush()
}

func main() {
	leitor = bufio.NewScanner(os.Stdin)
	leitor.Split(bufio.ScanWords)
	escritor = bufio.NewWriter(os.Stdout)
	depuracao, _ := os.Create("depuracao.txt")

	lerInicio()
	fmt.Fprintln(os.Stderr, "mapa:", largura, "x", altura, "| turnos:", totalTurnos, "| hospitais:", hospitais)

	escritor.WriteString("PRONTO\n")
	escritor.Flush()

	for lerTurno() {
		jogar()
		if turnoAtual%100 == 0 {
			fmt.Fprintln(depuracao, "turno", turnoAtual, "abertos", porID.quantidade)
		}
	}
}
