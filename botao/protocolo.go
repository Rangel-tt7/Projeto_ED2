package main

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
