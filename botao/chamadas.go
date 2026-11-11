package main

type Chamado struct {
	ID, X, Y, Valor, Prazo, DistanciaHospital, Reserva int
	Aberto                                             bool
}

func fechar(id int) {
	busca, achou := porID.Buscar(id) //retorna o índice do chamado no slice e se achou
	if !achou {
		return // o chamado não está na lista
	}
	chamados[busca].Aberto = false //muda o chamado para fechado
	porID.Remover(id)
}
