package main

type Grade struct {
	Lado, Colunas, Linhas, ValorMax int
	Baldes                          [][]int
}

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
