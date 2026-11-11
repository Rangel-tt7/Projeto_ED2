package main

type TabelaHash struct {
	chaves     []int
	valores    []int
	ocupada    []bool
	quantidade int // quantos pares estão guardados
}

// CriarTabela cria uma tabela vazia com "tamanho" posições
func CriarTabela(tamanho int) *TabelaHash {
	return &TabelaHash{
		chaves:  make([]int, tamanho),
		valores: make([]int, tamanho),
		ocupada: make([]bool, tamanho),
	}
}

// tamanho devolve quantas posições o vetor tem
func (tabela *TabelaHash) tamanho() int {
	return len(tabela.chaves)
}

// hash diz em qual posição a chave deveria ficar
func (tabela *TabelaHash) hash(chave int) int {
	return chave % tabela.tamanho()
}

// Buscar devolve o valor da chave e true, ou 0 e false se ela não está na tabela
func (tabela *TabelaHash) Buscar(chave int) (int, bool) {
	pos := tabela.hash(chave)
	for tabela.ocupada[pos] {
		if tabela.chaves[pos] == chave {
			return tabela.valores[pos], true
		}
		pos = (pos + 1) % tabela.tamanho() // sondagem linear: tenta a próxima, dando a volta no fim
	}
	return 0, false
}

// Inserir guarda o par. Se a chave já existe, só troca o valor
func (tabela *TabelaHash) Inserir(chave, valor int) {
	pos := tabela.hash(chave)
	for tabela.ocupada[pos] {
		if tabela.chaves[pos] == chave {
			tabela.valores[pos] = valor
			return
		}
		pos = (pos + 1) % tabela.tamanho()
	}
	tabela.chaves[pos] = chave
	tabela.valores[pos] = valor
	tabela.ocupada[pos] = true
	tabela.quantidade++

	if 10*tabela.quantidade > 7*tabela.tamanho() { // fator de carga passou de 0,7
		tabela.aumentarTabela(proximoPrimo(2 * tabela.tamanho()))
	}
}

// aumentarTabela troca o vetor por um com "novoTamanho" posições e reinsere todos os pares
func (tabela *TabelaHash) aumentarTabela(novoTamanho int) {
	tabelaNova := CriarTabela(novoTamanho)
	for pos := 0; pos < tabela.tamanho(); pos++ {
		if tabela.ocupada[pos] {
			tabelaNova.Inserir(tabela.chaves[pos], tabela.valores[pos])
		}
	}
	*tabela = *tabelaNova
}

// proximoPrimo devolve o menor primo maior ou igual a numero
func proximoPrimo(numero int) int {
	for {
		ehPrimo := numero >= 2
		for divisor := 2; divisor*divisor <= numero; divisor++ {
			if numero%divisor == 0 {
				ehPrimo = false
				break
			}
		}
		if ehPrimo {
			return numero
		}
		numero++
	}
}

// distancia devolve quantos passos a sondagem dá da posição inicio até a posição
// fim, dando a volta no fim do vetor. Com tamanho 11: distancia(9, 1) = 3 (9, 10, 0, 1)
func (tabela *TabelaHash) distancia(inicio, fim int) int {
	return (fim - inicio + tabela.tamanho()) % tabela.tamanho()
}

// Remover tira a chave da tabela. Devolve true se ela estava lá
func (tabela *TabelaHash) Remover(chave int) bool {
	buraco := tabela.hash(chave)
	for tabela.ocupada[buraco] && tabela.chaves[buraco] != chave {
		buraco = (buraco + 1) % tabela.tamanho()
	}
	if !tabela.ocupada[buraco] {
		return false // a chave não estava na tabela
	}
	tabela.ocupada[buraco] = false // abre o buraco
	tabela.quantidade--

	// Puxa para o buraco as chaves seguintes que passaram por ele ao serem
	// inseridas. Só uma posição vazia encerra o percurso.
	atual := (buraco + 1) % tabela.tamanho()
	for tabela.ocupada[atual] {
		posCerta := tabela.hash(tabela.chaves[atual])                             // onde a chave de "atual" deveria estar
		if tabela.distancia(buraco, atual) <= tabela.distancia(posCerta, atual) { // o buraco está no caminho de posCerta até atual
			tabela.chaves[buraco] = tabela.chaves[atual]
			tabela.valores[buraco] = tabela.valores[atual]
			tabela.ocupada[buraco] = true
			tabela.ocupada[atual] = false
			buraco = atual // o buraco passa a ser a posição atual
		}
		atual = (atual + 1) % tabela.tamanho()
	}
	return true
}
