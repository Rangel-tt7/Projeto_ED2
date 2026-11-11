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
	ID, X, Y, Valor, Prazo, DistanciaHospital, Reserva int
	Aberto                                             bool
}

type Ambulancia struct {
	ID, X, Y, Alvo int
	Livre          bool
}

type Grade struct {
	Lado, Colunas, Linhas, ValorMax int
	Baldes                          [][]int
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
var chamadosComuns *Grade    //grade para os chamados comuns, menos que valem até 20
var chamadosUrgentes *Grade  //grade para os chamados urgentes, mais de 20

func criarBalde(lado, valorMax int) *Grade {
	colunas := (largura + lado - 1) / lado
	linhas := (altura + lado - 1) / lado
	baldes := make([][]int, colunas*linhas)
	return &Grade{
		Lado:     lado,
		Colunas:  colunas,
		Linhas:   linhas,
		ValorMax: valorMax,
		Baldes:   baldes,
	}
}

func (grade *Grade) balde(x, y int) int {
	coluna := x / grade.Lado
	linha := y / grade.Lado
	return linha*grade.Colunas + coluna
}

func pontosPorTurno(ambulancia *Ambulancia, chamado *Chamado) float64 { // calcula a nota do chamado para a ambulância, ou 0 se não puder atender
	if !chamado.Aberto {
		return 0 // já foi recolhido ou venceu
	}
	if chamado.Reserva != -1 && chamado.Reserva != ambulancia.ID {
		return 0 // outra ambulância já vai buscar
	}

	distAteChamado := manhattan(ambulancia.X, ambulancia.Y, chamado.X, chamado.Y) // a sua conta, guardada numa variável

	if turnoAtual+distAteChamado-1 > chamado.Prazo {
		return 0 // não chega antes do prazo
	}
	if turnoAtual+distAteChamado+chamado.DistanciaHospital-1 > totalTurnos {
		return 0 // não entrega antes do fim da partida
	}

	return float64(chamado.Valor) / float64(distAteChamado+chamado.DistanciaHospital)
}

func buscarChamado(ambulancia *Ambulancia) int { //devolve a posição do melhor chamado para a ambulância -1 se não achar nenhum
	melhor := -1        // posição do melhor chamado achado
	melhorPontos := 0.0 // pontos por turno dele

	for _, grade := range []*Grade{chamadosComuns, chamadosUrgentes} { // olha as duas grades
		colunaAmb := ambulancia.X / grade.Lado // em que balde a ambulância está
		linhaAmb := ambulancia.Y / grade.Lado

		for anel := 0; anel < grade.Colunas || anel < grade.Linhas; anel++ {
			if anel > 0 {
				distMin := (anel-1)*grade.Lado + 1 // o mais perto que um chamado deste anel pode estar
				if float64(grade.ValorMax)/float64(distMin) <= melhorPontos {
					break // nada daqui para frente ganha do melhor achado
				}
				if distMin > totalTurnos-turnoAtual {
					break // longe demais para dar tempo antes do fim
				}
			}

			for linha := linhaAmb - anel; linha <= linhaAmb+anel; linha++ {
				for coluna := colunaAmb - anel; coluna <= colunaAmb+anel; coluna++ {
					if linha < 0 || linha >= grade.Linhas || coluna < 0 || coluna >= grade.Colunas {
						continue // fora do mapa
					}
					if linha != linhaAmb-anel && linha != linhaAmb+anel &&
						coluna != colunaAmb-anel && coluna != colunaAmb+anel {
						continue // está dentro do anel: já foi visto antes
					}

					for _, posicao := range grade.Baldes[linha*grade.Colunas+coluna] {
						pontos := pontosPorTurno(ambulancia, &chamados[posicao]) // a nota deste chamado para esta ambulância
						if pontos > melhorPontos {
							melhorPontos = pontos
							melhor = posicao
						}
					}
				}
			}
		}
	}
	return melhor
}

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
		ambulancias[id] = Ambulancia{id, x, y, -1, true}
	}

	lerPalavra() // BLOQUEIOS
	qtdBloqueios := lerNumero()
	for i := 0; i < qtdBloqueios*4; i++ {
		lerNumero()
	}

	lerPalavra()                                            // FIM
	porID = CriarTabela(proximoPrimo(largura * altura / 2)) //tamanho inicial da tabela hash
	chamados = make([]Chamado, 0, largura*altura/4)         //capacidade inicial do slice, make([]Tipo, tamanho, capacidade)
	expira = make([][]int, totalTurnos+2000)                //uma página para cada turno
	chamadosComuns = criarBalde(16, 20)                     //balde para os chamados comuns, menos que valem até 20
	chamadosUrgentes = criarBalde(32, 100)                  //balde para os chamados urgentes, mais de 20
}

// lerTurno lê a mensagem de um turno. Devolve false quando a partida acaba
func lerTurno() bool {
	if lerPalavra() == nil { // TURNO, ou nil se o servidor fechou a entrada
		return false
	}
	turnoAtual = lerNumero()

	for _, id := range expira[turnoAtual-1] { //fecha os chamados dos turnos que expiraram no turno anterior
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
		posicao := len(chamados)
		hospitalX, hospitalY := hospitalMaisProximo(x, y)
		distHospital := manhattan(x, y, hospitalX, hospitalY)
		chamados = append(chamados, Chamado{id, x, y, valor, prazo, distHospital, -1, true}) //adiciona o chamado ao slice
		porID.Inserir(id, posicao)                                                           // guarda na hash
		expira[prazo] = append(expira[prazo], id)                                            // anota o id na página do prazo
		grade := chamadosComuns
		if valor > 20 {
			grade = chamadosUrgentes
		}
		indiceBalde := grade.balde(x, y)
		grade.Baldes[indiceBalde] = append(grade.Baldes[indiceBalde], posicao)
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
	for i := range ambulancias {
		ambulancia := &ambulancias[i]
		movimento := "P" // ambulância não se mexe
		if ambulancia.Alvo != -1 {
			chamados[ambulancia.Alvo].Reserva = -1 // larga a reserva do turno anterior
			ambulancia.Alvo = -1
		}
		if ambulancia.Livre { //ambulância livre procura chamado
			posicao := buscarChamado(ambulancia)
			if posicao != -1 {
				chamado := &chamados[posicao]
				chamado.Reserva = ambulancia.ID //reserva o chamado para a ambulância
				ambulancia.Alvo = posicao       // guarda o alvo para poder largar no próximo turno
				movimento = direcao(ambulancia.X, ambulancia.Y, chamado.X, chamado.Y)
			}
		} else {
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
