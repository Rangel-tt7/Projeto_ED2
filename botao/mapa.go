package main

// Hospital guarda a posição de um hospital
type Hospital struct {
	X, Y int
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
