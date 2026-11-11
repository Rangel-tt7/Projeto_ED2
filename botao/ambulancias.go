package main

import "fmt"

type Ambulancia struct {
	ID, X, Y, Alvo int
	Livre          bool
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
